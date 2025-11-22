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

// CountAll は全ユーザー数を取得します
func (r *UserRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByStatus はステータス別ユーザー数を取得します
func (r *UserRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByRole はロール別ユーザー数を取得します
func (r *UserRepository) CountByRole(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		Role  string
		Count int64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Select("role, COUNT(*) as count").
		Group("role").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	roleCount := make(map[string]int64)
	for _, result := range results {
		roleCount[result.Role] = result.Count
	}

	return roleCount, nil
}

// CountByDepartment は部門別ユーザー数を取得します
func (r *UserRepository) CountByDepartment(ctx context.Context) ([]struct {
	Department string
	Count      int64
}, error) {
	var results []struct {
		Department string
		Count      int64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Select("COALESCE(department, '未設定') as department, COUNT(*) as count").
		Group("department").
		Order("count DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// GetLastLoginStats は過去N日間のログイン統計を取得します
func (r *UserRepository) GetLastLoginStats(ctx context.Context, days int) ([]struct {
	Date  string
	Count int64
}, error) {
	var results []struct {
		Date  string
		Count int64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Select("DATE(last_login) as date, COUNT(*) as count").
		Where("last_login IS NOT NULL AND last_login >= NOW() - INTERVAL '?' DAY", days).
		Group("DATE(last_login)").
		Order("date DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
