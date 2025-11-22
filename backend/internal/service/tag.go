package service

import (
	"context"
	"errors"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TagService はタグのサービスです
type TagService struct {
	repo   repository.TagRepository
	logger *zap.Logger
}

// NewTagService はTagServiceを作成します
func NewTagService(repo repository.TagRepository, logger *zap.Logger) *TagService {
	return &TagService{
		repo:   repo,
		logger: logger,
	}
}

// List はタグ一覧を取得します
func (s *TagService) List(ctx context.Context) ([]*model.TagResponse, error) {
	tags, err := s.repo.FindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to find tags", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = tag.ToResponse()
	}

	return responses, nil
}

// GetByID はタグを取得します
func (s *TagService) GetByID(ctx context.Context, id uint) (*model.TagResponse, error) {
	tag, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("TAG_001", errors.New("tag not found"))
		}
		s.logger.Error("Failed to find tag", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return tag.ToResponse(), nil
}

// Create はタグを作成します
func (s *TagService) Create(ctx context.Context, req *model.CreateTagRequest) (*model.TagResponse, error) {
	// タグ名の重複チェック
	exists, err := s.repo.ExistsByName(ctx, req.Name)
	if err != nil {
		s.logger.Error("Failed to check tag existence", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}
	if exists {
		return nil, util.NewValidationError("TAG_002", errors.New("tag name already exists"))
	}

	// タグの作成
	tag := &model.Tag{
		Name:        req.Name,
		Color:       "#6B7280",
		Description: req.Description,
	}

	// リクエストでカラーが指定されている場合は上書き
	if req.Color != "" {
		tag.Color = req.Color
	}

	if err := s.repo.Create(ctx, tag); err != nil {
		s.logger.Error("Failed to create tag", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return tag.ToResponse(), nil
}

// Update はタグを更新します
func (s *TagService) Update(ctx context.Context, id uint, req *model.UpdateTagRequest) (*model.TagResponse, error) {
	// 既存のタグを取得
	tag, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("TAG_001", errors.New("tag not found"))
		}
		s.logger.Error("Failed to find tag", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// タグ名の更新と重複チェック
	if req.Name != nil {
		exists, err := s.repo.ExistsByNameExcludingID(ctx, *req.Name, id)
		if err != nil {
			s.logger.Error("Failed to check tag existence", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		if exists {
			return nil, util.NewValidationError("TAG_002", errors.New("tag name already exists"))
		}
		tag.Name = *req.Name
	}

	// カラーの更新
	if req.Color != nil {
		tag.Color = *req.Color
	}

	// 説明の更新
	if req.Description != nil {
		tag.Description = *req.Description
	}

	if err := s.repo.Update(ctx, tag); err != nil {
		s.logger.Error("Failed to update tag", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return tag.ToResponse(), nil
}

// Delete はタグを削除します
func (s *TagService) Delete(ctx context.Context, id uint) error {
	// 既存のタグを取得
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.NewNotFoundError("TAG_001", errors.New("tag not found"))
		}
		s.logger.Error("Failed to find tag", zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete tag", zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return nil
}
