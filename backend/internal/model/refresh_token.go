package model

import (
	"time"
)

// RefreshToken はリフレッシュトークンモデルです
type RefreshToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	TokenID   string    `gorm:"uniqueIndex;not null;size:255" json:"token_id"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	Revoked   bool      `gorm:"not null;default:false" json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName はテーブル名を指定します
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsValid はトークンが有効かチェックします
func (rt *RefreshToken) IsValid() bool {
	return !rt.Revoked && time.Now().Before(rt.ExpiresAt)
}
