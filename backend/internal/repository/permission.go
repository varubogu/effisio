package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
)

// PermissionRepository はPermissionのリポジトリインターフェースです
type PermissionRepository interface {
	FindAll(ctx context.Context) ([]*model.Permission, error)
	FindByID(ctx context.Context, id uint) (*model.Permission, error)
	FindByName(ctx context.Context, name string) (*model.Permission, error)
	FindByIDs(ctx context.Context, ids []uint) ([]*model.Permission, error)
	Create(ctx context.Context, permission *model.Permission) error
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id uint) error
}

// permissionRepository はPermissionRepositoryの実装です
type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository は新しいPermissionRepositoryを作成します
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// FindAll は全ての権限を取得します
func (r *permissionRepository) FindAll(ctx context.Context) ([]*model.Permission, error) {
	var permissions []*model.Permission
	if err := r.db.WithContext(ctx).Order("resource, action").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// FindByID はIDで権限を取得します
func (r *permissionRepository) FindByID(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	if err := r.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByName は名前で権限を取得します
func (r *permissionRepository) FindByName(ctx context.Context, name string) (*model.Permission, error) {
	var permission model.Permission
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByIDs は複数のIDで権限を取得します
func (r *permissionRepository) FindByIDs(ctx context.Context, ids []uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// Create は新しい権限を作成します
func (r *permissionRepository) Create(ctx context.Context, permission *model.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

// Update は権限を更新します
func (r *permissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

// Delete は権限を削除します
func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Permission{}, id).Error
}
