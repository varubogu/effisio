package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// AuditLog は監査ログを表します
type AuditLog struct {
	ID         uint64    `gorm:"primarykey" json:"id"`
	UserID     *uint     `gorm:"index" json:"user_id"`
	Action     string    `gorm:"not null;size:100;index" json:"action"`
	Resource   string    `gorm:"not null;size:50;index" json:"resource"`
	ResourceID string    `gorm:"size:255;index:idx_audit_logs_resource_id" json:"resource_id"`
	IPAddress  string    `gorm:"size:45" json:"ip_address"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	Changes    JSONB     `gorm:"type:jsonb" json:"changes"`
	CreatedAt  time.Time `gorm:"index:,sort:desc" json:"created_at"`
}

// TableName はテーブル名を指定します
func (AuditLog) TableName() string {
	return "audit_logs"
}

// JSONB はPostgreSQLのJSONB型を扱うためのカスタム型です
type JSONB map[string]interface{}

// Value はJSONBをデータベースに保存する際の値を返します
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan はデータベースからJSONBを読み込みます
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// AuditLogResponse はAuditLogのレスポンス用構造体です
type AuditLogResponse struct {
	ID         uint64                 `json:"id"`
	UserID     *uint                  `json:"user_id"`
	Username   string                 `json:"username,omitempty"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id,omitempty"`
	IPAddress  string                 `json:"ip_address,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	Changes    map[string]interface{} `json:"changes,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// ToResponse はAuditLogをAuditLogResponseに変換します
func (a *AuditLog) ToResponse() *AuditLogResponse {
	response := &AuditLogResponse{
		ID:         a.ID,
		UserID:     a.UserID,
		Action:     a.Action,
		Resource:   a.Resource,
		ResourceID: a.ResourceID,
		IPAddress:  a.IPAddress,
		UserAgent:  a.UserAgent,
		Changes:    a.Changes,
		CreatedAt:  a.CreatedAt,
	}
	return response
}

// 監査ログのアクション定数
const (
	// 認証関連
	ActionLoginSuccess  = "LOGIN_SUCCESS"
	ActionLoginFailed   = "LOGIN_FAILED"
	ActionLogout        = "LOGOUT"
	ActionLogoutAll     = "LOGOUT_ALL"
	ActionTokenRefresh  = "TOKEN_REFRESH"

	// ユーザー管理
	ActionUserCreate = "USER_CREATE"
	ActionUserUpdate = "USER_UPDATE"
	ActionUserDelete = "USER_DELETE"
	ActionUserView   = "USER_VIEW"

	// ロール管理
	ActionRoleCreate = "ROLE_CREATE"
	ActionRoleUpdate = "ROLE_UPDATE"
	ActionRoleDelete = "ROLE_DELETE"

	// 権限管理
	ActionPermissionCreate = "PERMISSION_CREATE"
	ActionPermissionUpdate = "PERMISSION_UPDATE"
	ActionPermissionDelete = "PERMISSION_DELETE"

	// 設定変更
	ActionSettingUpdate = "SETTING_UPDATE"
)

// リソース定数
const (
	ResourceAuth       = "auth"
	ResourceUser       = "users"
	ResourceRole       = "roles"
	ResourcePermission = "permissions"
	ResourceSetting    = "settings"
)
