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

// PermissionService はPermission関連のビジネスロジックを提供します
type PermissionService struct {
	repo   repository.PermissionRepository
	logger *zap.Logger
}

// NewPermissionService は新しいPermissionServiceを作成します
func NewPermissionService(repo repository.PermissionRepository, logger *zap.Logger) *PermissionService {
	return &PermissionService{
		repo:   repo,
		logger: logger,
	}
}

// List は全ての権限を取得します
func (s *PermissionService) List(ctx context.Context) ([]*model.PermissionResponse, error) {
	permissions, err := s.repo.FindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch permissions", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.PermissionResponse, len(permissions))
	for i, perm := range permissions {
		responses[i] = perm.ToResponse()
	}

	return responses, nil
}

// GetByID はIDで権限を取得します
func (s *PermissionService) GetByID(ctx context.Context, id uint) (*model.PermissionResponse, error) {
	permission, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("PERM_001", err)
		}
		s.logger.Error("Failed to fetch permission", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return permission.ToResponse(), nil
}

// Create は新しい権限を作成します
func (s *PermissionService) Create(ctx context.Context, req *model.CreatePermissionRequest) (*model.PermissionResponse, error) {
	// 同じ名前の権限が既に存在するかチェック
	_, err := s.repo.FindByName(ctx, req.Name)
	if err == nil {
		return nil, util.NewConflictError("PERM_002", errors.New("permission with this name already exists"))
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to check permission existence", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	permission := &model.Permission{
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, permission); err != nil {
		s.logger.Error("Failed to create permission", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Permission created", zap.Uint("id", permission.ID), zap.String("name", permission.Name))
	return permission.ToResponse(), nil
}

// Update は権限を更新します
func (s *PermissionService) Update(ctx context.Context, id uint, req *model.UpdatePermissionRequest) (*model.PermissionResponse, error) {
	permission, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("PERM_001", err)
		}
		s.logger.Error("Failed to fetch permission", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 更新内容を適用
	if req.Name != nil {
		// 名前の重複チェック（自分以外）
		existing, err := s.repo.FindByName(ctx, *req.Name)
		if err == nil && existing.ID != id {
			return nil, util.NewConflictError("PERM_002", errors.New("permission with this name already exists"))
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("Failed to check permission name", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		permission.Name = *req.Name
	}
	if req.Resource != nil {
		permission.Resource = *req.Resource
	}
	if req.Action != nil {
		permission.Action = *req.Action
	}
	if req.Description != nil {
		permission.Description = *req.Description
	}

	if err := s.repo.Update(ctx, permission); err != nil {
		s.logger.Error("Failed to update permission", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Permission updated", zap.Uint("id", id))
	return permission.ToResponse(), nil
}

// Delete は権限を削除します
func (s *PermissionService) Delete(ctx context.Context, id uint) error {
	// 存在確認
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.NewNotFoundError("PERM_001", err)
		}
		s.logger.Error("Failed to fetch permission", zap.Error(err), zap.Uint("id", id))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete permission", zap.Error(err), zap.Uint("id", id))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Permission deleted", zap.Uint("id", id))
	return nil
}
