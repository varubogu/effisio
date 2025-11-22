package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// TaskActivityType はタスクアクティビティのタイプを表します
type TaskActivityType string

const (
	ActivityTypeCreated       TaskActivityType = "created"
	ActivityTypeUpdated       TaskActivityType = "updated"
	ActivityTypeStatusChanged TaskActivityType = "status_changed"
	ActivityTypeAssigned      TaskActivityType = "assigned"
	ActivityTypeCommented     TaskActivityType = "commented"
	ActivityTypeTagAdded      TaskActivityType = "tag_added"
	ActivityTypeTagRemoved    TaskActivityType = "tag_removed"
)

// ActivityMetadata はアクティビティのメタデータを表します
type ActivityMetadata map[string]interface{}

// Value はdriver.Valuerインターフェースを実装します
func (m ActivityMetadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan はsql.Scannerインターフェースを実装します
func (m *ActivityMetadata) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

// TaskActivity はタスクアクティビティログエンティティを表します
type TaskActivity struct {
	ID           uint             `gorm:"primarykey" json:"id"`
	TaskID       uint             `gorm:"not null;index" json:"task_id"`
	UserID       uint             `gorm:"not null;index" json:"user_id"`
	ActivityType TaskActivityType `gorm:"not null;size:50;index" json:"activity_type"`
	FieldName    string           `gorm:"size:100" json:"field_name"`
	OldValue     string           `gorm:"type:text" json:"old_value"`
	NewValue     string           `gorm:"type:text" json:"new_value"`
	Metadata     ActivityMetadata `gorm:"type:jsonb" json:"metadata"`
	CreatedAt    time.Time        `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"created_at"`

	// リレーション
	Task *Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName はテーブル名を指定します
func (TaskActivity) TableName() string {
	return "task_activities"
}

// TaskActivityResponse はタスクアクティビティのレスポンス用構造体です
type TaskActivityResponse struct {
	ID           uint             `json:"id"`
	TaskID       uint             `json:"task_id"`
	UserID       uint             `json:"user_id"`
	User         *UserResponse    `json:"user,omitempty"`
	ActivityType TaskActivityType `json:"activity_type"`
	FieldName    string           `json:"field_name,omitempty"`
	OldValue     string           `json:"old_value,omitempty"`
	NewValue     string           `json:"new_value,omitempty"`
	Metadata     ActivityMetadata `json:"metadata,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
}

// ToResponse はTaskActivityをTaskActivityResponseに変換します
func (ta *TaskActivity) ToResponse() *TaskActivityResponse {
	response := &TaskActivityResponse{
		ID:           ta.ID,
		TaskID:       ta.TaskID,
		UserID:       ta.UserID,
		ActivityType: ta.ActivityType,
		FieldName:    ta.FieldName,
		OldValue:     ta.OldValue,
		NewValue:     ta.NewValue,
		Metadata:     ta.Metadata,
		CreatedAt:    ta.CreatedAt,
	}

	if ta.User != nil {
		response.User = ta.User.ToResponse()
	}

	return response
}
