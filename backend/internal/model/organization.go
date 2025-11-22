package model

import "time"

// Organization は組織を表します
type Organization struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:50" json:"code"`
	ParentID    *uint          `gorm:"index" json:"parent_id"`
	Path        string         `gorm:"size:1000;index" json:"path"`
	Level       int            `gorm:"not null;default:0;index" json:"level"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// 関連
	Parent   *Organization   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Organization  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName はテーブル名を指定します
func (Organization) TableName() string {
	return "organizations"
}

// OrganizationResponse はOrganizationのレスポンス用構造体です
type OrganizationResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Code        string                  `json:"code"`
	ParentID    *uint                   `json:"parent_id"`
	Path        string                  `json:"path"`
	Level       int                     `json:"level"`
	Description string                  `json:"description"`
	Children    []OrganizationResponse  `json:"children,omitempty"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

// ToResponse はOrganizationをOrganizationResponseに変換します
func (o *Organization) ToResponse() *OrganizationResponse {
	response := &OrganizationResponse{
		ID:          o.ID,
		Name:        o.Name,
		Code:        o.Code,
		ParentID:    o.ParentID,
		Path:        o.Path,
		Level:       o.Level,
		Description: o.Description,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}

	// 子組織が読み込まれている場合は含める
	if len(o.Children) > 0 {
		response.Children = make([]OrganizationResponse, len(o.Children))
		for i, child := range o.Children {
			response.Children[i] = *child.ToResponse()
		}
	}

	return response
}

// CreateOrganizationRequest はOrganization作成時のリクエスト構造体です
type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Code        string `json:"code" binding:"omitempty,max=50"`
	ParentID    *uint  `json:"parent_id" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
}

// UpdateOrganizationRequest はOrganization更新時のリクエスト構造体です
type UpdateOrganizationRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Code        *string `json:"code" binding:"omitempty,max=50"`
	ParentID    *uint   `json:"parent_id" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
}

// OrganizationTreeResponse は組織ツリーのレスポンスです
type OrganizationTreeResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
}
