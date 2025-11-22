package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
)

// RoleRepository はRoleのリポジトリインターフェースです
type RoleRepository interface {
	FindAll(ctx context.Context) ([]*model.Role, error)
	FindByID(ctx context.Context, id uint) (*model.Role, error)
	FindByIDWithPermissions(ctx context.Context, id uint) (*model.Role, error)
	FindByName(ctx context.Context, name string) (*model.Role, error)
	FindByNameWithPermissions(ctx context.Context, name string) (*model.Role, error)
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id uint) error
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	RemoveAllPermissions(ctx context.Context, roleID uint) error
}

// roleRepository はRoleRepositoryの実装です
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository は新しいRoleRepositoryを作成します
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// FindAll は全てのロールを取得します
func (r *roleRepository) FindAll(ctx context.Context) ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.db.WithContext(ctx).Order("name").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// FindByID はIDでロールを取得します
func (r *roleRepository) FindByID(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByIDWithPermissions はIDでロールと権限を取得します
func (r *roleRepository) FindByIDWithPermissions(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName は名前でロールを取得します
func (r *roleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByNameWithPermissions は名前でロールと権限を取得します
func (r *roleRepository) FindByNameWithPermissions(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Create は新しいロールを作成します
func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// Update はロールを更新します
func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// Delete はロールを削除します
func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// AssignPermissions はロールに権限を割り当てます
func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	role, err := r.FindByID(ctx, roleID)
	if err != nil {
		return err
	}

	// 権限を取得
	var permissions []*model.Permission
	if err := r.db.WithContext(ctx).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	// GORMの Association を使って権限を割り当て
	// Replace は既存の関連をすべて削除して新しい関連を作成
	return r.db.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions)
}

// RemoveAllPermissions はロールから全ての権限を削除します
func (r *roleRepository) RemoveAllPermissions(ctx context.Context, roleID uint) error {
	role, err := r.FindByID(ctx, roleID)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(role).Association("Permissions").Clear()
}
