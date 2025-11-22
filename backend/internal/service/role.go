package service

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// RoleService はRole関連のビジネスロジックを提供します
type RoleService struct {
	repo           repository.RoleRepository
	permissionRepo repository.PermissionRepository
	logger         *zap.Logger
}

// NewRoleService は新しいRoleServiceを作成します
func NewRoleService(
	repo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	logger *zap.Logger,
) *RoleService {
	return &RoleService{
		repo:           repo,
		permissionRepo: permissionRepo,
		logger:         logger,
	}
}

// List は全てのロールを取得します
func (s *RoleService) List(ctx context.Context, withPermissions bool) ([]*model.RoleResponse, error) {
	var roles []*model.Role
	var err error

	if withPermissions {
		// 権限付きで取得する場合は個別に取得してPreloadする
		allRoles, err := s.repo.FindAll(ctx)
		if err != nil {
			s.logger.Error("Failed to fetch roles", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}

		roles = make([]*model.Role, len(allRoles))
		for i, role := range allRoles {
			roleWithPerms, err := s.repo.FindByIDWithPermissions(ctx, role.ID)
			if err != nil {
				s.logger.Error("Failed to fetch role with permissions", zap.Error(err), zap.Uint("role_id", role.ID))
				return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
			}
			roles[i] = roleWithPerms
		}
	} else {
		roles, err = s.repo.FindAll(ctx)
		if err != nil {
			s.logger.Error("Failed to fetch roles", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
	}

	responses := make([]*model.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = role.ToResponse()
	}

	return responses, nil
}

// GetByID はIDでロールを取得します
func (s *RoleService) GetByID(ctx context.Context, id uint, withPermissions bool) (*model.RoleResponse, error) {
	var role *model.Role
	var err error

	if withPermissions {
		role, err = s.repo.FindByIDWithPermissions(ctx, id)
	} else {
		role, err = s.repo.FindByID(ctx, id)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("ROLE_001", err)
		}
		s.logger.Error("Failed to fetch role", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return role.ToResponse(), nil
}

// GetPermissionsForRole はロール名から権限名のリストを取得します
func (s *RoleService) GetPermissionsForRole(ctx context.Context, roleName string) ([]string, error) {
	role, err := s.repo.FindByNameWithPermissions(ctx, roleName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Role not found", zap.String("role", roleName))
			return []string{}, nil
		}
		s.logger.Error("Failed to fetch role permissions", zap.Error(err), zap.String("role", roleName))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	permissionNames := make([]string, len(role.Permissions))
	for i, perm := range role.Permissions {
		permissionNames[i] = perm.Name
	}

	return permissionNames, nil
}

// Create は新しいロールを作成します
func (s *RoleService) Create(ctx context.Context, req *model.CreateRoleRequest) (*model.RoleResponse, error) {
	// 同じ名前のロールが既に存在するかチェック
	_, err := s.repo.FindByName(ctx, req.Name)
	if err == nil {
		return nil, util.NewConflictError("ROLE_002", errors.New("role with this name already exists"))
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to check role existence", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 権限IDが指定されている場合は、その権限が存在するかチェック
	if len(req.PermissionIDs) > 0 {
		permissions, err := s.permissionRepo.FindByIDs(ctx, req.PermissionIDs)
		if err != nil {
			s.logger.Error("Failed to fetch permissions", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		if len(permissions) != len(req.PermissionIDs) {
			return nil, util.NewValidationError("ROLE_003", errors.New("some permission IDs are invalid"))
		}
	}

	role := &model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, role); err != nil {
		s.logger.Error("Failed to create role", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 権限を割り当て
	if len(req.PermissionIDs) > 0 {
		if err := s.repo.AssignPermissions(ctx, role.ID, req.PermissionIDs); err != nil {
			s.logger.Error("Failed to assign permissions to role", zap.Error(err), zap.Uint("role_id", role.ID))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
	}

	s.logger.Info("Role created", zap.Uint("id", role.ID), zap.String("name", role.Name))

	// 権限付きで取得して返す
	return s.GetByID(ctx, role.ID, true)
}

// Update はロールを更新します
func (s *RoleService) Update(ctx context.Context, id uint, req *model.UpdateRoleRequest) (*model.RoleResponse, error) {
	role, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("ROLE_001", err)
		}
		s.logger.Error("Failed to fetch role", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 更新内容を適用
	if req.Name != nil {
		// 名前の重複チェック（自分以外）
		existing, err := s.repo.FindByName(ctx, *req.Name)
		if err == nil && existing.ID != id {
			return nil, util.NewConflictError("ROLE_002", errors.New("role with this name already exists"))
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("Failed to check role name", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := s.repo.Update(ctx, role); err != nil {
		s.logger.Error("Failed to update role", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 権限の更新
	if req.PermissionIDs != nil {
		if len(*req.PermissionIDs) > 0 {
			// 権限が存在するかチェック
			permissions, err := s.permissionRepo.FindByIDs(ctx, *req.PermissionIDs)
			if err != nil {
				s.logger.Error("Failed to fetch permissions", zap.Error(err))
				return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
			}
			if len(permissions) != len(*req.PermissionIDs) {
				return nil, util.NewValidationError("ROLE_003", errors.New("some permission IDs are invalid"))
			}

			if err := s.repo.AssignPermissions(ctx, id, *req.PermissionIDs); err != nil {
				s.logger.Error("Failed to assign permissions to role", zap.Error(err), zap.Uint("role_id", id))
				return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
			}
		} else {
			// 空配列が指定された場合は全ての権限を削除
			if err := s.repo.RemoveAllPermissions(ctx, id); err != nil {
				s.logger.Error("Failed to remove permissions from role", zap.Error(err), zap.Uint("role_id", id))
				return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
			}
		}
	}

	s.logger.Info("Role updated", zap.Uint("id", id))

	// 権限付きで取得して返す
	return s.GetByID(ctx, id, true)
}

// Delete はロールを削除します
func (s *RoleService) Delete(ctx context.Context, id uint) error {
	// 存在確認
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.NewNotFoundError("ROLE_001", err)
		}
		s.logger.Error("Failed to fetch role", zap.Error(err), zap.Uint("id", id))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete role", zap.Error(err), zap.Uint("id", id))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Role deleted", zap.Uint("id", id))
	return nil
}
