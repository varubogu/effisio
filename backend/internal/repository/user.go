package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// UserRepository はユーザーデータアクセスを提供します
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository は新しいUserRepositoryを作成します
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindAll は全てのユーザーを取得します（ページネーション付き）
func (r *UserRepository) FindAll(ctx context.Context, params *util.PaginationParams) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 総件数を取得
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション付きで取得
	err := r.db.WithContext(ctx).
		Offset(params.Offset).
		Limit(params.PerPage).
		Order("id ASC").
		Find(&users).Error

	return users, total, err
}

// FindByID はIDでユーザーを取得します
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail はメールアドレスでユーザーを取得します
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername はユーザー名でユーザーを取得します
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create は新しいユーザーを作成します
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update はユーザー情報を更新します
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete はユーザーを削除します（ソフトデリート）
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// ExistsByEmail はメールアドレスの存在確認をします
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByUsername はユーザー名の存在確認をします
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
