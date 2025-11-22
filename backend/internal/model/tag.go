package model

import "time"

// Tag はタグエンティティを表します
type Tag struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"not null;unique;size:50" json:"name"`
	Color       string    `gorm:"not null;size:7;default:#6B7280" json:"color"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// リレーション（多対多）
	Tasks []Task `gorm:"many2many:task_tags;" json:"tasks,omitempty"`
}

// TableName はテーブル名を指定します
func (Tag) TableName() string {
	return "tags"
}

// TagResponse はタグのレスポンス用構造体です
type TagResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse はTagをTagResponseに変換します
func (t *Tag) ToResponse() *TagResponse {
	return &TagResponse{
		ID:          t.ID,
		Name:        t.Name,
		Color:       t.Color,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// CreateTagRequest はタグ作成リクエストを表します
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Color       string `json:"color" binding:"omitempty,len=7"`
	Description string `json:"description"`
}

// UpdateTagRequest はタグ更新リクエストを表します
type UpdateTagRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=50"`
	Color       *string `json:"color" binding:"omitempty,len=7"`
	Description *string `json:"description"`
}

// TaskTag はタスクタグ中間テーブルを表します
type TaskTag struct {
	TaskID    uint      `gorm:"primaryKey" json:"task_id"`
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName はテーブル名を指定します
func (TaskTag) TableName() string {
	return "task_tags"
}
