package model

import (
	"time"

	"gorm.io/gorm"
)

// User はユーザーモデルです
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	FullName     string         `gorm:"size:100" json:"full_name"`
	Department   string         `gorm:"size:100" json:"department"`
	PasswordHash string         `gorm:"not null;size:255;column:password_hash" json:"-"` // JSONには含めない
	Role         string         `gorm:"not null;size:20;default:'user'" json:"role"`
	Status       string         `gorm:"not null;size:20;default:'active'" json:"status"`
	LastLogin    *time.Time     `json:"last_login"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // ソフトデリート
}

// TableName はテーブル名を指定します
func (User) TableName() string {
	return "users"
}

// ステータス定数
const (
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
)

// ロール定数
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleUser    = "user"
	RoleViewer  = "viewer"
)

// IsValidStatus はステータスが有効かチェックします
func IsValidStatus(status string) bool {
	return status == UserStatusActive || status == UserStatusInactive || status == UserStatusSuspended
}

// IsValidRole はロールが有効かチェックします
func IsValidRole(role string) bool {
	return role == RoleAdmin || role == RoleManager || role == RoleUser || role == RoleViewer
}

// CreateUserRequest はユーザー作成リクエストです
type CreateUserRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Email      string `json:"email" binding:"required,email"`
	FullName   string `json:"full_name" binding:"max=100"`
	Department string `json:"department" binding:"max=100"`
	Password   string `json:"password" binding:"required,min=8,max=72"`
	Role       string `json:"role" binding:"required,oneof=admin manager user viewer"`
}

// UpdateUserRequest はユーザー更新リクエストです
type UpdateUserRequest struct {
	Email      *string `json:"email" binding:"omitempty,email"`
	FullName   *string `json:"full_name" binding:"omitempty,max=100"`
	Department *string `json:"department" binding:"omitempty,max=100"`
	Role       *string `json:"role" binding:"omitempty,oneof=admin manager user viewer"`
	Status     *string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// UserResponse はユーザーレスポンスです（パスワードを除外）
type UserResponse struct {
	ID         uint       `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	FullName   string     `json:"full_name"`
	Department string     `json:"department"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	LastLogin  *time.Time `json:"last_login"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ToResponse はUserをUserResponseに変換します
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		FullName:   u.FullName,
		Department: u.Department,
		Role:       u.Role,
		Status:     u.Status,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
