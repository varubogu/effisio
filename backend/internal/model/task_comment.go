package model

import "time"

// TaskComment はタスクコメントエンティティを表します
type TaskComment struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	TaskID    uint      `gorm:"not null;index" json:"task_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// リレーション
	Task *Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName はテーブル名を指定します
func (TaskComment) TableName() string {
	return "task_comments"
}

// TaskCommentResponse はタスクコメントのレスポンス用構造体です
type TaskCommentResponse struct {
	ID        uint          `json:"id"`
	TaskID    uint          `json:"task_id"`
	UserID    uint          `json:"user_id"`
	User      *UserResponse `json:"user,omitempty"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// ToResponse はTaskCommentをTaskCommentResponseに変換します
func (tc *TaskComment) ToResponse() *TaskCommentResponse {
	response := &TaskCommentResponse{
		ID:        tc.ID,
		TaskID:    tc.TaskID,
		UserID:    tc.UserID,
		Content:   tc.Content,
		CreatedAt: tc.CreatedAt,
		UpdatedAt: tc.UpdatedAt,
	}

	if tc.User != nil {
		response.User = tc.User.ToResponse()
	}

	return response
}

// CreateTaskCommentRequest はタスクコメント作成リクエストを表します
type CreateTaskCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// UpdateTaskCommentRequest はタスクコメント更新リクエストを表します
type UpdateTaskCommentRequest struct {
	Content string `json:"content" binding:"required"`
}
