package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// AuditLogService は監査ログ関連のビジネスロジックを提供します
type AuditLogService struct {
	repo     repository.AuditLogRepository
	userRepo repository.UserRepository
	logger   *zap.Logger
}

// NewAuditLogService は新しいAuditLogServiceを作成します
func NewAuditLogService(
	repo repository.AuditLogRepository,
	userRepo repository.UserRepository,
	logger *zap.Logger,
) *AuditLogService {
	return &AuditLogService{
		repo:     repo,
		userRepo: userRepo,
		logger:   logger,
	}
}

// LogEntry は監査ログ記録用のエントリです
type LogEntry struct {
	UserID     *uint
	Action     string
	Resource   string
	ResourceID string
	IPAddress  string
	UserAgent  string
	Changes    map[string]interface{}
}

// Log は監査ログを記録します（非同期実行推奨）
func (s *AuditLogService) Log(ctx context.Context, entry *LogEntry) error {
	log := &model.AuditLog{
		UserID:     entry.UserID,
		Action:     entry.Action,
		Resource:   entry.Resource,
		ResourceID: entry.ResourceID,
		IPAddress:  entry.IPAddress,
		UserAgent:  entry.UserAgent,
		Changes:    entry.Changes,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, log); err != nil {
		s.logger.Error("Failed to create audit log",
			zap.Error(err),
			zap.String("action", entry.Action),
			zap.String("resource", entry.Resource),
		)
		return err
	}

	s.logger.Debug("Audit log created",
		zap.String("action", entry.Action),
		zap.String("resource", entry.Resource),
		zap.String("resource_id", entry.ResourceID),
	)

	return nil
}

// LogAsync は監査ログを非同期で記録します（パフォーマンス最適化）
func (s *AuditLogService) LogAsync(entry *LogEntry) {
	go func() {
		ctx := context.Background()
		if err := s.Log(ctx, entry); err != nil {
			// エラーログは既にLog内で記録されている
		}
	}()
}

// List は監査ログ一覧を取得します
func (s *AuditLogService) List(
	ctx context.Context,
	filter *repository.AuditLogFilter,
	pagination *util.PaginationParams,
) (*util.PaginatedResponse, error) {
	logs, total, err := s.repo.FindAll(ctx, filter, pagination)
	if err != nil {
		s.logger.Error("Failed to fetch audit logs", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// ユーザー名を取得してレスポンスに含める
	responses := make([]*model.AuditLogResponse, len(logs))
	for i, log := range logs {
		response := log.ToResponse()

		// ユーザー名を取得
		if log.UserID != nil {
			user, err := s.userRepo.FindByID(ctx, *log.UserID)
			if err == nil {
				response.Username = user.Username
			}
		}

		responses[i] = response
	}

	paginationInfo := util.CalculatePaginationInfo(total, pagination.Page, pagination.PerPage)

	return &util.PaginatedResponse{
		Code:       200,
		Message:    "success",
		Data:       responses,
		Pagination: paginationInfo,
	}, nil
}

// GetByID はIDで監査ログを取得します
func (s *AuditLogService) GetByID(ctx context.Context, id uint64) (*model.AuditLogResponse, error) {
	log, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("AUDIT_001", err)
		}
		s.logger.Error("Failed to fetch audit log", zap.Error(err), zap.Uint64("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	response := log.ToResponse()

	// ユーザー名を取得
	if log.UserID != nil {
		user, err := s.userRepo.FindByID(ctx, *log.UserID)
		if err == nil {
			response.Username = user.Username
		}
	}

	return response, nil
}

// ヘルパー関数: 変更を記録するためのChangesマップを作成

// MakeChanges は変更前後の値をChangesマップとして作成します
func MakeChanges(before, after interface{}) map[string]interface{} {
	return map[string]interface{}{
		"before": before,
		"after":  after,
	}
}

// MakeChangesFromStructs は構造体から変更内容を抽出します
func MakeChangesFromStructs(before, after interface{}) map[string]interface{} {
	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(after)

	var beforeMap, afterMap map[string]interface{}
	json.Unmarshal(beforeJSON, &beforeMap)
	json.Unmarshal(afterJSON, &afterMap)

	return map[string]interface{}{
		"before": beforeMap,
		"after":  afterMap,
	}
}

// 便利な監査ログ記録メソッド

// LogLoginSuccess はログイン成功を記録します
func (s *AuditLogService) LogLoginSuccess(userID uint, ipAddress, userAgent string) {
	s.LogAsync(&LogEntry{
		UserID:    &userID,
		Action:    model.ActionLoginSuccess,
		Resource:  model.ResourceAuth,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}

// LogLoginFailed はログイン失敗を記録します
func (s *AuditLogService) LogLoginFailed(username, ipAddress, userAgent string) {
	s.LogAsync(&LogEntry{
		Action:    model.ActionLoginFailed,
		Resource:  model.ResourceAuth,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Changes: map[string]interface{}{
			"username": username,
		},
	})
}

// LogLogout はログアウトを記録します
func (s *AuditLogService) LogLogout(userID uint, ipAddress, userAgent string) {
	s.LogAsync(&LogEntry{
		UserID:    &userID,
		Action:    model.ActionLogout,
		Resource:  model.ResourceAuth,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}

// LogUserCreate はユーザー作成を記録します
func (s *AuditLogService) LogUserCreate(userID uint, newUser *model.User, ipAddress, userAgent string) {
	s.LogAsync(&LogEntry{
		UserID:     &userID,
		Action:     model.ActionUserCreate,
		Resource:   model.ResourceUser,
		ResourceID: fmt.Sprintf("%d", newUser.ID),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Changes: map[string]interface{}{
			"after": map[string]interface{}{
				"username": newUser.Username,
				"email":    newUser.Email,
				"role":     newUser.Role,
				"status":   newUser.Status,
			},
		},
	})
}

// LogUserUpdate はユーザー更新を記録します
func (s *AuditLogService) LogUserUpdate(userID uint, targetUserID uint, before, after *model.User, ipAddress, userAgent string) {
	s.LogAsync(&LogEntry{
		UserID:     &userID,
		Action:     model.ActionUserUpdate,
		Resource:   model.ResourceUser,
		ResourceID: fmt.Sprintf("%d", targetUserID),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Changes:    MakeChangesFromStructs(before, after),
	})
}

// LogUserDelete はユーザー削除を記録します
func (s *AuditLogService) LogUserDelete(userID uint, deletedUser *model.User, ipAddress, userAgent string) {
	s.LogAsync(&LogEntry{
		UserID:     &userID,
		Action:     model.ActionUserDelete,
		Resource:   model.ResourceUser,
		ResourceID: fmt.Sprintf("%d", deletedUser.ID),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Changes: map[string]interface{}{
			"before": map[string]interface{}{
				"username": deletedUser.Username,
				"email":    deletedUser.Email,
				"role":     deletedUser.Role,
			},
		},
	})
}

// LogRoleCreate はロール作成を記録します
func (s *AuditLogService) LogRoleCreate(userID uint, role *model.Role, ipAddress, userAgent string) {
	s.LogAsync(&LogEntry{
		UserID:     &userID,
		Action:     model.ActionRoleCreate,
		Resource:   model.ResourceRole,
		ResourceID: fmt.Sprintf("%d", role.ID),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Changes: map[string]interface{}{
			"after": map[string]interface{}{
				"name":        role.Name,
				"description": role.Description,
			},
		},
	})
}

// LogRoleUpdate はロール更新を記録します
func (s *AuditLogService) LogRoleUpdate(userID uint, roleID uint, before, after *model.Role, ipAddress, userAgent string) {
	s.LogAsync(&LogEntry{
		UserID:     &userID,
		Action:     model.ActionRoleUpdate,
		Resource:   model.ResourceRole,
		ResourceID: strconv.Itoa(int(roleID)),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Changes:    MakeChangesFromStructs(before, after),
	})
}

// LogRoleDelete はロール削除を記録します
func (s *AuditLogService) LogRoleDelete(userID uint, role *model.Role, ipAddress, userAgent string) {
	s.LogAsync(&LogEntry{
		UserID:     &userID,
		Action:     model.ActionRoleDelete,
		Resource:   model.ResourceRole,
		ResourceID: fmt.Sprintf("%d", role.ID),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Changes: map[string]interface{}{
			"before": map[string]interface{}{
				"name": role.Name,
			},
		},
	})
}
