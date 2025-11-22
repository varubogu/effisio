package model

import "time"

// TaskStatus はタスクのステータスを表します
type TaskStatus string

const (
	TaskStatusTODO        TaskStatus = "TODO"
	TaskStatusInProgress  TaskStatus = "IN_PROGRESS"
	TaskStatusInReview    TaskStatus = "IN_REVIEW"
	TaskStatusDone        TaskStatus = "DONE"
	TaskStatusCancelled   TaskStatus = "CANCELLED"
)

// TaskPriority はタスクの優先度を表します
type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "LOW"
	TaskPriorityMedium TaskPriority = "MEDIUM"
	TaskPriorityHigh   TaskPriority = "HIGH"
	TaskPriorityUrgent TaskPriority = "URGENT"
)

// Task はタスクエンティティを表します
type Task struct {
	ID             uint          `gorm:"primarykey" json:"id"`
	Title          string        `gorm:"not null;size:255" json:"title"`
	Description    string        `gorm:"type:text" json:"description"`
	Status         TaskStatus    `gorm:"not null;size:20;default:TODO" json:"status"`
	Priority       TaskPriority  `gorm:"not null;size:10;default:MEDIUM" json:"priority"`
	AssignedToID   *uint         `gorm:"column:assigned_to;index" json:"assigned_to_id"`
	CreatedByID    uint          `gorm:"column:created_by;not null;index" json:"created_by_id"`
	OrganizationID *uint         `gorm:"index" json:"organization_id"`
	DueDate        *time.Time    `json:"due_date"`
	CompletedAt    *time.Time    `json:"completed_at"`
	CreatedAt      time.Time     `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// リレーション
	AssignedTo   *User         `gorm:"foreignKey:AssignedToID" json:"assigned_to,omitempty"`
	CreatedBy    *User         `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	Organization *Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// TableName はテーブル名を指定します
func (Task) TableName() string {
	return "tasks"
}

// TaskResponse はタスクのレスポンス用構造体です
type TaskResponse struct {
	ID             uint                     `json:"id"`
	Title          string                   `json:"title"`
	Description    string                   `json:"description"`
	Status         TaskStatus               `json:"status"`
	Priority       TaskPriority             `json:"priority"`
	AssignedToID   *uint                    `json:"assigned_to_id"`
	AssignedTo     *UserResponse            `json:"assigned_to,omitempty"`
	CreatedByID    uint                     `json:"created_by_id"`
	CreatedBy      *UserResponse            `json:"created_by,omitempty"`
	OrganizationID *uint                    `json:"organization_id"`
	Organization   *OrganizationResponse    `json:"organization,omitempty"`
	DueDate        *time.Time               `json:"due_date"`
	CompletedAt    *time.Time               `json:"completed_at"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

// ToResponse はTaskをTaskResponseに変換します
func (t *Task) ToResponse() *TaskResponse {
	response := &TaskResponse{
		ID:             t.ID,
		Title:          t.Title,
		Description:    t.Description,
		Status:         t.Status,
		Priority:       t.Priority,
		AssignedToID:   t.AssignedToID,
		CreatedByID:    t.CreatedByID,
		OrganizationID: t.OrganizationID,
		DueDate:        t.DueDate,
		CompletedAt:    t.CompletedAt,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}

	if t.AssignedTo != nil {
		response.AssignedTo = t.AssignedTo.ToResponse()
	}

	if t.CreatedBy != nil {
		response.CreatedBy = t.CreatedBy.ToResponse()
	}

	if t.Organization != nil {
		response.Organization = t.Organization.ToResponse()
	}

	return response
}

// CreateTaskRequest はタスク作成リクエストを表します
type CreateTaskRequest struct {
	Title          string        `json:"title" binding:"required,max=255"`
	Description    string        `json:"description"`
	Status         TaskStatus    `json:"status" binding:"omitempty,oneof=TODO IN_PROGRESS IN_REVIEW DONE CANCELLED"`
	Priority       TaskPriority  `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	AssignedToID   *uint         `json:"assigned_to_id"`
	OrganizationID *uint         `json:"organization_id"`
	DueDate        *time.Time    `json:"due_date"`
}

// UpdateTaskRequest はタスク更新リクエストを表します
type UpdateTaskRequest struct {
	Title          *string       `json:"title" binding:"omitempty,max=255"`
	Description    *string       `json:"description"`
	Status         *TaskStatus   `json:"status" binding:"omitempty,oneof=TODO IN_PROGRESS IN_REVIEW DONE CANCELLED"`
	Priority       *TaskPriority `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	AssignedToID   *uint         `json:"assigned_to_id"`
	OrganizationID *uint         `json:"organization_id"`
	DueDate        *time.Time    `json:"due_date"`
}

// TaskFilter はタスクの検索条件を表します
type TaskFilter struct {
	Status         *TaskStatus   `form:"status"`
	Priority       *TaskPriority `form:"priority"`
	AssignedToID   *uint         `form:"assigned_to_id"`
	CreatedByID    *uint         `form:"created_by_id"`
	OrganizationID *uint         `form:"organization_id"`
	DueBefore      *time.Time    `form:"due_before"`
	DueAfter       *time.Time    `form:"due_after"`
	Page           int           `form:"page,default=1"`
	PageSize       int           `form:"page_size,default=20"`
	SortBy         string        `form:"sort_by,default=created_at"`
	SortOrder      string        `form:"sort_order,default=desc"`
}
