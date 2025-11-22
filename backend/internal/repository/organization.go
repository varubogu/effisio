package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
)

// OrganizationRepository はOrganizationのリポジトリインターフェースです
type OrganizationRepository interface {
	FindAll(ctx context.Context) ([]*model.Organization, error)
	FindByID(ctx context.Context, id uint) (*model.Organization, error)
	FindByCode(ctx context.Context, code string) (*model.Organization, error)
	FindRoots(ctx context.Context) ([]*model.Organization, error)
	FindChildren(ctx context.Context, parentID uint) ([]*model.Organization, error)
	FindDescendants(ctx context.Context, id uint) ([]*model.Organization, error)
	FindTree(ctx context.Context) ([]*model.Organization, error)
	Create(ctx context.Context, org *model.Organization) error
	Update(ctx context.Context, org *model.Organization) error
	Delete(ctx context.Context, id uint) error
	ExistsByCode(ctx context.Context, code string) (bool, error)
}

// organizationRepository はOrganizationRepositoryの実装です
type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository は新しいOrganizationRepositoryを作成します
func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

// FindAll は全ての組織を取得します
func (r *organizationRepository) FindAll(ctx context.Context) ([]*model.Organization, error) {
	var organizations []*model.Organization
	if err := r.db.WithContext(ctx).Order("path").Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

// FindByID はIDで組織を取得します
func (r *organizationRepository) FindByID(ctx context.Context, id uint) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.WithContext(ctx).First(&org, id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// FindByCode はコードで組織を取得します
func (r *organizationRepository) FindByCode(ctx context.Context, code string) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// FindRoots はルート組織（親がない組織）を取得します
func (r *organizationRepository) FindRoots(ctx context.Context) ([]*model.Organization, error) {
	var organizations []*model.Organization
	if err := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Order("name").
		Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

// FindChildren は指定組織の直接の子組織を取得します
func (r *organizationRepository) FindChildren(ctx context.Context, parentID uint) ([]*model.Organization, error) {
	var organizations []*model.Organization
	if err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("name").
		Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

// FindDescendants は指定組織の全ての子孫組織を取得します（パスを使用）
func (r *organizationRepository) FindDescendants(ctx context.Context, id uint) ([]*model.Organization, error) {
	// まず対象組織を取得してパスを確認
	org, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var descendants []*model.Organization
	if err := r.db.WithContext(ctx).
		Where("path LIKE ?", org.Path+"/%").
		Order("path").
		Find(&descendants).Error; err != nil {
		return nil, err
	}
	return descendants, nil
}

// FindTree は組織ツリー全体を取得します（階層構造）
func (r *organizationRepository) FindTree(ctx context.Context) ([]*model.Organization, error) {
	var organizations []*model.Organization

	// 全ての組織を取得（パスでソート）
	if err := r.db.WithContext(ctx).
		Order("path").
		Find(&organizations).Error; err != nil {
		return nil, err
	}

	// マップを作成してIDから組織を検索できるようにする
	orgMap := make(map[uint]*model.Organization)
	for i := range organizations {
		orgMap[organizations[i].ID] = organizations[i]
		organizations[i].Children = []model.Organization{} // 初期化
	}

	// 階層構造を構築
	var roots []*model.Organization
	for i := range organizations {
		org := organizations[i]
		if org.ParentID == nil {
			// ルート組織
			roots = append(roots, &org)
		} else {
			// 親組織の子として追加
			if parent, exists := orgMap[*org.ParentID]; exists {
				parent.Children = append(parent.Children, org)
			}
		}
	}

	return roots, nil
}

// Create は新しい組織を作成します
func (r *organizationRepository) Create(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

// Update は組織を更新します
func (r *organizationRepository) Update(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

// Delete は組織を削除します
func (r *organizationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Organization{}, id).Error
}

// ExistsByCode はコードの存在チェックを行います
func (r *organizationRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Organization{}).
		Where("code = ?", code).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
