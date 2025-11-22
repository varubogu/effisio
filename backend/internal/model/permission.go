package model

import "time"

// Permission はシステムの権限を表します
type Permission struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Resource    string    `gorm:"not null;size:50;index" json:"resource"`
	Action      string    `gorm:"not null;size:50;index" json:"action"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Many-to-Many関連
	Roles []Role `gorm:"many2many:role_permissions;" json:"-"`
}

// TableName はテーブル名を指定します
func (Permission) TableName() string {
	return "permissions"
}

// PermissionResponse はPermissionのレスポンス用構造体です
type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse はPermissionをPermissionResponseに変換します
func (p *Permission) ToResponse() *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// CreatePermissionRequest はPermission作成時のリクエスト構造体です
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Resource    string `json:"resource" binding:"required,max=50"`
	Action      string `json:"action" binding:"required,max=50"`
	Description string `json:"description" binding:"max=255"`
}

// UpdatePermissionRequest はPermission更新時のリクエスト構造体です
type UpdatePermissionRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Resource    *string `json:"resource" binding:"omitempty,max=50"`
	Action      *string `json:"action" binding:"omitempty,max=50"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}
