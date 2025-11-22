package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
)

// RefreshTokenRepository はリフレッシュトークンのデータアクセスを提供します
type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository は新しいRefreshTokenRepositoryを作成します
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

// Create はリフレッシュトークンを作成します
func (r *RefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindByTokenID はトークンIDでリフレッシュトークンを取得します
func (r *RefreshTokenRepository) FindByTokenID(ctx context.Context, tokenID string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	if err := r.db.WithContext(ctx).Where("token_id = ?", tokenID).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

// FindByUserID はユーザーIDで有効なリフレッシュトークンを全て取得します
func (r *RefreshTokenRepository) FindByUserID(ctx context.Context, userID uint) ([]*model.RefreshToken, error) {
	var tokens []*model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Find(&tokens).Error
	return tokens, err
}

// Revoke はトークンを無効化します
func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenID string) error {
	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("token_id = ?", tokenID).
		Update("revoked", true).Error
}

// RevokeAllByUserID はユーザーの全リフレッシュトークンを無効化します
func (r *RefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

// DeleteExpired は期限切れのトークンを削除します
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&model.RefreshToken{}).Error
}
