package model

import "time"

// Role はシステムのロールを表します
type Role struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Many-to-Many関連
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"-"`
}

// TableName はテーブル名を指定します
func (Role) TableName() string {
	return "roles"
}

// RoleResponse はRoleのレスポンス用構造体です
type RoleResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Permissions []PermissionResponse  `json:"permissions,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// ToResponse はRoleをRoleResponseに変換します
func (r *Role) ToResponse() *RoleResponse {
	response := &RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	// Permissionsが読み込まれている場合は含める
	if len(r.Permissions) > 0 {
		response.Permissions = make([]PermissionResponse, len(r.Permissions))
		for i, perm := range r.Permissions {
			response.Permissions[i] = *perm.ToResponse()
		}
	}

	return response
}

// CreateRoleRequest はRole作成時のリクエスト構造体です
type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required,max=50"`
	Description   string `json:"description" binding:"max=255"`
	PermissionIDs []uint `json:"permission_ids" binding:"omitempty"`
}

// UpdateRoleRequest はRole更新時のリクエスト構造体です
type UpdateRoleRequest struct {
	Name          *string `json:"name" binding:"omitempty,max=50"`
	Description   *string `json:"description" binding:"omitempty,max=255"`
	PermissionIDs *[]uint `json:"permission_ids" binding:"omitempty"`
}
