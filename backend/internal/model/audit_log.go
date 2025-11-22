package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// AuditLog は監査ログモデルです
type AuditLog struct {
	ID            uint            `gorm:"primarykey" json:"id"`
	UserID        uint            `gorm:"not null;index" json:"user_id"`
	Action        string          `gorm:"not null;size:50;index" json:"action"`
	ResourceType  string          `gorm:"not null;size:50;index" json:"resource_type"`
	ResourceID    string          `gorm:"not null;size:50;index" json:"resource_id"`
	Changes       datatypes.JSONType `gorm:"type:jsonb;not null" json:"changes"`
	IPAddress     string          `gorm:"size:45" json:"ip_address"`
	UserAgent     string          `gorm:"type:text" json:"user_agent"`
	Status        string          `gorm:"not null;size:20;default:'success';index" json:"status"`
	ErrorMessage  string          `gorm:"type:text" json:"error_message"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// TableName はテーブル名を指定します
func (AuditLog) TableName() string {
	return "audit_logs"
}

// アクション定数
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionLogin  = "login"
	ActionLogout = "logout"
)

// リソースタイプ定数
const (
	ResourceTypeUser         = "user"
	ResourceTypeRole         = "role"
	ResourceTypeOrganization = "organization"
	ResourceTypeAuditLog     = "audit_log"
)

// ステータス定数
const (
	AuditStatusSuccess = "success"
	AuditStatusFailed  = "failed"
)

// AuditLogChanges は変更内容を表現します
type AuditLogChanges struct {
	Before map[string]interface{} `json:"before"`
	After  map[string]interface{} `json:"after"`
}

// MarshalJSON は AuditLogChanges を JSON にマーシャルします
func (a AuditLogChanges) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"before": a.Before,
		"after":  a.After,
	})
}

// UnmarshalJSON は JSON から AuditLogChanges をアンマーシャルします
func (a *AuditLogChanges) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if before, ok := raw["before"].(map[string]interface{}); ok {
		a.Before = before
	}
	if after, ok := raw["after"].(map[string]interface{}); ok {
		a.After = after
	}

	return nil
}

// CreateAuditLogRequest は監査ログ作成リクエストです
type CreateAuditLogRequest struct {
	UserID       uint                   `json:"user_id" binding:"required"`
	Action       string                 `json:"action" binding:"required,oneof=create read update delete login logout"`
	ResourceType string                 `json:"resource_type" binding:"required"`
	ResourceID   string                 `json:"resource_id" binding:"required"`
	Changes      AuditLogChanges        `json:"changes"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
	Status       string                 `json:"status" binding:"required,oneof=success failed"`
	ErrorMessage string                 `json:"error_message"`
}

// AuditLogResponse は監査ログレスポンスです
type AuditLogResponse struct {
	ID            uint            `json:"id"`
	UserID        uint            `json:"user_id"`
	Action        string          `json:"action"`
	ResourceType  string          `json:"resource_type"`
	ResourceID    string          `json:"resource_id"`
	Changes       AuditLogChanges `json:"changes"`
	IPAddress     string          `json:"ip_address"`
	UserAgent     string          `json:"user_agent"`
	Status        string          `json:"status"`
	ErrorMessage  string          `json:"error_message"`
	CreatedAt     time.Time       `json:"created_at"`
}

// ToResponse は AuditLog をレスポンスに変換します
func (a *AuditLog) ToResponse() *AuditLogResponse {
	var changes AuditLogChanges
	if err := json.Unmarshal(a.Changes, &changes); err != nil {
		// JSONのアンマーシャルに失敗した場合は空の変更を返す
		changes = AuditLogChanges{
			Before: make(map[string]interface{}),
			After:  make(map[string]interface{}),
		}
	}

	return &AuditLogResponse{
		ID:           a.ID,
		UserID:       a.UserID,
		Action:       a.Action,
		ResourceType: a.ResourceType,
		ResourceID:   a.ResourceID,
		Changes:      changes,
		IPAddress:    a.IPAddress,
		UserAgent:    a.UserAgent,
		Status:       a.Status,
		ErrorMessage: a.ErrorMessage,
		CreatedAt:    a.CreatedAt,
	}
}
