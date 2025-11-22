package repository

import (
	"context"

	"github.com/varubogu/effisio/backend/internal/model"
	"gorm.io/gorm"
)

// TagRepository はタグのリポジトリインターフェースです
type TagRepository interface {
	FindAll(ctx context.Context) ([]*model.Tag, error)
	FindByID(ctx context.Context, id uint) (*model.Tag, error)
	FindByName(ctx context.Context, name string) (*model.Tag, error)
	FindByIDs(ctx context.Context, ids []uint) ([]*model.Tag, error)
	Create(ctx context.Context, tag *model.Tag) error
	Update(ctx context.Context, tag *model.Tag) error
	Delete(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByNameExcludingID(ctx context.Context, name string, excludeID uint) (bool, error)
}

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository はTagRepositoryを作成します
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

// FindAll はタグ一覧を取得します
func (r *tagRepository) FindAll(ctx context.Context) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindByID はIDでタグを取得します
func (r *tagRepository) FindByID(ctx context.Context, id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.WithContext(ctx).First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindByName は名前でタグを取得します
func (r *tagRepository) FindByName(ctx context.Context, name string) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindByIDs は複数のIDでタグを取得します
func (r *tagRepository) FindByIDs(ctx context.Context, ids []uint) ([]*model.Tag, error) {
	if len(ids) == 0 {
		return []*model.Tag{}, nil
	}

	var tags []*model.Tag
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// Create はタグを作成します
func (r *tagRepository) Create(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// Update はタグを更新します
func (r *tagRepository) Update(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// Delete はタグを削除します
func (r *tagRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Tag{}, id).Error
}

// ExistsByName は名前でタグの存在確認をします
func (r *tagRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Tag{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByNameExcludingID は指定IDを除外して名前でタグの存在確認をします
func (r *tagRepository) ExistsByNameExcludingID(ctx context.Context, name string, excludeID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Tag{}).
		Where("name = ? AND id != ?", name, excludeID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
