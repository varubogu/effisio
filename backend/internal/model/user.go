package model

import (
	"time"

	"gorm.io/gorm"
)

// User はユーザーモデルです
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string         `gorm:"not null;size:255" json:"-"` // JSONには含めない
	Role      string         `gorm:"not null;size:20;default:'user'" json:"role"`
	IsActive  bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // ソフトデリート
}

// TableName はテーブル名を指定します
func (User) TableName() string {
	return "users"
}

// UserRole は許可されたユーザーロール
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleUser    UserRole = "user"
	RoleViewer  UserRole = "viewer"
)

// IsValidRole はロールが有効かチェックします
func IsValidRole(role string) bool {
	switch UserRole(role) {
	case RoleAdmin, RoleManager, RoleUser, RoleViewer:
		return true
	default:
		return false
	}
}

// CreateUserRequest はユーザー作成リクエストです
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"omitempty"`
}

// UpdateUserRequest はユーザー更新リクエストです
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"omitempty"`
	IsActive *bool  `json:"is_active" binding:"omitempty"`
}

// UserResponse はユーザーレスポンスです（パスワードを除外）
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse はUserをUserResponseに変換します
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
