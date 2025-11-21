package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
)

// UserService はユーザー関連のビジネスロジックを提供します
type UserService struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

// NewUserService は新しいUserServiceを作成します
func NewUserService(repo *repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// List はユーザー一覧を取得します
func (s *UserService) List(ctx context.Context) ([]*model.UserResponse, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// GetByID はIDでユーザーを取得します
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user.ToResponse(), nil
}

// Create は新しいユーザーを作成します
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// ロールのバリデーション
	role := req.Role
	if role == "" {
		role = string(model.RoleUser) // デフォルトロール
	} else if !model.IsValidRole(role) {
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("User created successfully",
		zap.Uint("id", user.ID),
		zap.String("username", user.Username),
	)

	return user.ToResponse(), nil
}

// Update はユーザー情報を更新します
func (s *UserService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// 更新可能なフィールドを更新
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		if !model.IsValidRole(req.Role) {
			return nil, fmt.Errorf("invalid role: %s", req.Role)
		}
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info("User updated successfully",
		zap.Uint("id", user.ID),
		zap.String("username", user.Username),
	)

	return user.ToResponse(), nil
}

// Delete はユーザーを削除します（ソフトデリート）
func (s *UserService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Info("User deleted successfully", zap.Uint("id", id))

	return nil
}
