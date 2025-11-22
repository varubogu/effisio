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

// OrganizationService はOrganization関連のビジネスロジックを提供します
type OrganizationService struct {
	repo   repository.OrganizationRepository
	logger *zap.Logger
}

// NewOrganizationService は新しいOrganizationServiceを作成します
func NewOrganizationService(
	repo repository.OrganizationRepository,
	logger *zap.Logger,
) *OrganizationService {
	return &OrganizationService{
		repo:   repo,
		logger: logger,
	}
}

// List は全ての組織を取得します
func (s *OrganizationService) List(ctx context.Context) ([]*model.OrganizationResponse, error) {
	organizations, err := s.repo.FindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch organizations", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.OrganizationResponse, len(organizations))
	for i, org := range organizations {
		responses[i] = org.ToResponse()
	}

	return responses, nil
}

// GetByID はIDで組織を取得します
func (s *OrganizationService) GetByID(ctx context.Context, id uint) (*model.OrganizationResponse, error) {
	org, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("ORG_001", err)
		}
		s.logger.Error("Failed to fetch organization", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return org.ToResponse(), nil
}

// GetTree は組織ツリー全体を取得します
func (s *OrganizationService) GetTree(ctx context.Context) (*model.OrganizationTreeResponse, error) {
	roots, err := s.repo.FindTree(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch organization tree", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]model.OrganizationResponse, len(roots))
	for i, root := range roots {
		responses[i] = *root.ToResponse()
	}

	return &model.OrganizationTreeResponse{
		Organizations: responses,
	}, nil
}

// GetChildren は指定組織の子組織を取得します
func (s *OrganizationService) GetChildren(ctx context.Context, parentID uint) ([]*model.OrganizationResponse, error) {
	// 親組織の存在確認
	if _, err := s.repo.FindByID(ctx, parentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("ORG_001", err)
		}
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	children, err := s.repo.FindChildren(ctx, parentID)
	if err != nil {
		s.logger.Error("Failed to fetch children organizations", zap.Error(err), zap.Uint("parent_id", parentID))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.OrganizationResponse, len(children))
	for i, child := range children {
		responses[i] = child.ToResponse()
	}

	return responses, nil
}

// Create は新しい組織を作成します
func (s *OrganizationService) Create(ctx context.Context, req *model.CreateOrganizationRequest) (*model.OrganizationResponse, error) {
	// コードの重複チェック
	if req.Code != "" {
		exists, err := s.repo.ExistsByCode(ctx, req.Code)
		if err != nil {
			s.logger.Error("Failed to check organization code existence", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		if exists {
			return nil, util.NewConflictError("ORG_002", errors.New("organization with this code already exists"))
		}
	}

	// 親組織の存在確認
	if req.ParentID != nil {
		if _, err := s.repo.FindByID(ctx, *req.ParentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, util.NewNotFoundError("ORG_003", errors.New("parent organization not found"))
			}
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
	}

	org := &model.Organization{
		Name:        req.Name,
		Code:        req.Code,
		ParentID:    req.ParentID,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, org); err != nil {
		s.logger.Error("Failed to create organization", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Organization created", zap.Uint("id", org.ID), zap.String("name", org.Name))
	return org.ToResponse(), nil
}

// Update は組織を更新します
func (s *OrganizationService) Update(ctx context.Context, id uint, req *model.UpdateOrganizationRequest) (*model.OrganizationResponse, error) {
	org, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("ORG_001", err)
		}
		s.logger.Error("Failed to fetch organization", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 更新内容を適用
	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Code != nil {
		// コードの重複チェック（自分以外）
		exists, err := s.repo.ExistsByCode(ctx, *req.Code)
		if err != nil {
			s.logger.Error("Failed to check organization code", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		if exists {
			existing, _ := s.repo.FindByCode(ctx, *req.Code)
			if existing != nil && existing.ID != id {
				return nil, util.NewConflictError("ORG_002", errors.New("organization with this code already exists"))
			}
		}
		org.Code = *req.Code
	}
	if req.ParentID != nil {
		// 自分自身を親にできない
		if *req.ParentID == id {
			return nil, util.NewValidationError("ORG_004", errors.New("cannot set self as parent"))
		}

		// 親組織の存在確認
		parent, err := s.repo.FindByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, util.NewNotFoundError("ORG_003", errors.New("parent organization not found"))
			}
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}

		// 子孫を親にできない（循環参照防止）
		descendants, err := s.repo.FindDescendants(ctx, id)
		if err != nil {
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		for _, desc := range descendants {
			if desc.ID == *req.ParentID {
				return nil, util.NewValidationError("ORG_005", errors.New("cannot set descendant as parent"))
			}
		}

		// パスの検証（親のパスに自分が含まれていないか）
		if parent.Path != "" {
			// 親のパスに自分のIDが含まれている場合は循環
			// この検証は上の子孫チェックで十分カバーされている
		}

		org.ParentID = req.ParentID
	}
	if req.Description != nil {
		org.Description = *req.Description
	}

	if err := s.repo.Update(ctx, org); err != nil {
		s.logger.Error("Failed to update organization", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Organization updated", zap.Uint("id", id))
	return org.ToResponse(), nil
}

// Delete は組織を削除します
func (s *OrganizationService) Delete(ctx context.Context, id uint) error {
	// 存在確認
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.NewNotFoundError("ORG_001", err)
		}
		s.logger.Error("Failed to fetch organization", zap.Error(err), zap.Uint("id", id))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 子組織がある場合は削除できない（CASCADE DELETEは設定されているが、明示的にチェック）
	children, err := s.repo.FindChildren(ctx, id)
	if err != nil {
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}
	if len(children) > 0 {
		return util.NewValidationError("ORG_006", errors.New("cannot delete organization with children"))
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete organization", zap.Error(err), zap.Uint("id", id))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Organization deleted", zap.Uint("id", id))
	return nil
}
