package service

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// UserService はユーザー関連のビジネスロジックを提供します
type UserService struct {
	repo              *repository.UserRepository
	logger            *zap.Logger
	auditLogService   *AuditLogService
}

// NewUserService は新しいUserServiceを作成します
func NewUserService(repo *repository.UserRepository, logger *zap.Logger, auditLogService *AuditLogService) *UserService {
	return &UserService{
		repo:              repo,
		logger:            logger,
		auditLogService:   auditLogService,
	}
}

// List はユーザー一覧を取得します
func (s *UserService) List(ctx context.Context, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	users, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		s.logger.Error("Failed to fetch users", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return util.NewPaginatedResponse(responses, total, params), nil
}

// GetByID はIDでユーザーを取得します
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError(util.ErrCodeUserNotFound, err)
		}
		s.logger.Error("Failed to fetch user", zap.Uint("id", id), zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return user.ToResponse(), nil
}

// Create は新しいユーザーを作成します
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	// ユーザー名の重複チェック
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Error("Failed to check username existence", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}
	if exists {
		return nil, util.NewConflictError(util.ErrCodeUserAlreadyExists, errors.New("username already exists"))
	}

	// メールアドレスの重複チェック
	exists, err = s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Failed to check email existence", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}
	if exists {
		return nil, util.NewConflictError(util.ErrCodeUserAlreadyExists, errors.New("email already exists"))
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodePasswordHashError, err)
	}

	// ユーザーモデルを作成
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		Department:   req.Department,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		Status:       model.UserStatusActive,
	}

	// データベースに保存
	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		// 監査ログに失敗を記録
		if s.auditLogService != nil {
			auditReq := &model.CreateAuditLogRequest{
				UserID:       1, // システムユーザー（実装時に認証ユーザーから取得）
				Action:       model.ActionCreate,
				ResourceType: model.ResourceTypeUser,
				ResourceID:   req.Username,
				Status:       model.AuditStatusFailed,
				ErrorMessage: err.Error(),
			}
			s.auditLogService.LogAction(ctx, auditReq)
		}
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("User created", zap.Uint("id", user.ID), zap.String("username", user.Username))

	// 監査ログに成功を記録
	if s.auditLogService != nil {
		auditReq := &model.CreateAuditLogRequest{
			UserID:       1, // システムユーザー（実装時に認証ユーザーから取得）
			Action:       model.ActionCreate,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   user.Username,
			Changes: model.AuditLogChanges{
				Before: map[string]interface{}{},
				After: map[string]interface{}{
					"id":         user.ID,
					"username":   user.Username,
					"email":      user.Email,
					"full_name":  user.FullName,
					"department": user.Department,
					"role":       user.Role,
					"status":     user.Status,
				},
			},
			Status: model.AuditStatusSuccess,
		}
		s.auditLogService.LogAction(ctx, auditReq)
	}

	return user.ToResponse(), nil
}

// Update はユーザー情報を更新します
func (s *UserService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	// 既存ユーザーを取得
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError(util.ErrCodeUserNotFound, err)
		}
		s.logger.Error("Failed to fetch user", zap.Uint("id", id), zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 監査ログ用に更新前の値を保存
	beforeChanges := map[string]interface{}{
		"email":      user.Email,
		"full_name":  user.FullName,
		"department": user.Department,
		"role":       user.Role,
		"status":     user.Status,
	}

	// 更新データを適用
	if req.Email != nil {
		// メールアドレスの重複チェック（自分以外）
		existingUser, err := s.repo.FindByEmail(ctx, *req.Email)
		if err == nil && existingUser.ID != id {
			return nil, util.NewConflictError(util.ErrCodeUserAlreadyExists, errors.New("email already exists"))
		}
		user.Email = *req.Email
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Department != nil {
		user.Department = *req.Department
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	// 監査ログ用に更新後の値を保存
	afterChanges := map[string]interface{}{
		"email":      user.Email,
		"full_name":  user.FullName,
		"department": user.Department,
		"role":       user.Role,
		"status":     user.Status,
	}

	// データベースを更新
	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("Failed to update user", zap.Uint("id", id), zap.Error(err))
		// 監査ログに失敗を記録
		if s.auditLogService != nil {
			auditReq := &model.CreateAuditLogRequest{
				UserID:       1, // システムユーザー（実装時に認証ユーザーから取得）
				Action:       model.ActionUpdate,
				ResourceType: model.ResourceTypeUser,
				ResourceID:   user.Username,
				Status:       model.AuditStatusFailed,
				ErrorMessage: err.Error(),
			}
			s.auditLogService.LogAction(ctx, auditReq)
		}
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("User updated", zap.Uint("id", user.ID))

	// 監査ログに成功を記録
	if s.auditLogService != nil {
		auditReq := &model.CreateAuditLogRequest{
			UserID:       1, // システムユーザー（実装時に認証ユーザーから取得）
			Action:       model.ActionUpdate,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   user.Username,
			Changes: model.AuditLogChanges{
				Before: beforeChanges,
				After:  afterChanges,
			},
			Status: model.AuditStatusSuccess,
		}
		s.auditLogService.LogAction(ctx, auditReq)
	}

	return user.ToResponse(), nil
}

// Delete はユーザーを削除します（ソフトデリート）
func (s *UserService) Delete(ctx context.Context, id uint) error {
	// 存在確認
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.NewNotFoundError(util.ErrCodeUserNotFound, err)
		}
		s.logger.Error("Failed to fetch user", zap.Uint("id", id), zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 削除実行
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete user", zap.Uint("id", id), zap.Error(err))
		// 監査ログに失敗を記録
		if s.auditLogService != nil {
			auditReq := &model.CreateAuditLogRequest{
				UserID:       1, // システムユーザー（実装時に認証ユーザーから取得）
				Action:       model.ActionDelete,
				ResourceType: model.ResourceTypeUser,
				ResourceID:   user.Username,
				Status:       model.AuditStatusFailed,
				ErrorMessage: err.Error(),
			}
			s.auditLogService.LogAction(ctx, auditReq)
		}
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("User deleted", zap.Uint("id", id))

	// 監査ログに成功を記録
	if s.auditLogService != nil {
		auditReq := &model.CreateAuditLogRequest{
			UserID:       1, // システムユーザー（実装時に認証ユーザーから取得）
			Action:       model.ActionDelete,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   user.Username,
			Changes: model.AuditLogChanges{
				Before: map[string]interface{}{
					"id":         user.ID,
					"username":   user.Username,
					"email":      user.Email,
					"full_name":  user.FullName,
					"department": user.Department,
					"role":       user.Role,
					"status":     user.Status,
				},
				After: map[string]interface{}{},
			},
			Status: model.AuditStatusSuccess,
		}
		s.auditLogService.LogAction(ctx, auditReq)
	}

	return nil
}
