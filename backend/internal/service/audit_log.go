package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// AuditLogService は監査ログ関連のビジネスロジックを提供します
type AuditLogService struct {
	repo   *repository.AuditLogRepository
	logger *zap.Logger
}

// NewAuditLogService は新しいAuditLogServiceを作成します
func NewAuditLogService(repo *repository.AuditLogRepository, logger *zap.Logger) *AuditLogService {
	return &AuditLogService{
		repo:   repo,
		logger: logger,
	}
}

// LogAction はアクションを記録します
func (s *AuditLogService) LogAction(ctx context.Context, req *model.CreateAuditLogRequest) (*model.AuditLogResponse, error) {
	// リクエストの検証
	if err := s.validateCreateRequest(req); err != nil {
		s.logger.Warn("Invalid audit log request", zap.Error(err))
		return nil, util.NewBadRequestError(util.ErrCodeBadRequest, err)
	}

	// JSONB形式で変更内容を保存
	changesJSON, err := json.Marshal(req.Changes)
	if err != nil {
		s.logger.Error("Failed to marshal changes", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeInternalError, err)
	}

	// 監査ログモデルを作成
	auditLog := &model.AuditLog{
		UserID:       req.UserID,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Changes:      changesJSON,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		Status:       req.Status,
		ErrorMessage: req.ErrorMessage,
	}

	// データベースに保存
	if err := s.repo.Create(ctx, auditLog); err != nil {
		s.logger.Error("Failed to create audit log", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Audit log created",
		zap.Uint("user_id", auditLog.UserID),
		zap.String("action", auditLog.Action),
		zap.String("resource_type", auditLog.ResourceType),
		zap.String("resource_id", auditLog.ResourceID),
	)

	return auditLog.ToResponse(), nil
}

// GetByID はIDで監査ログを取得します
func (s *AuditLogService) GetByID(ctx context.Context, id uint) (*model.AuditLogResponse, error) {
	auditLog, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError(util.ErrCodeNotFound, err)
		}
		s.logger.Error("Failed to fetch audit log", zap.Uint("id", id), zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return auditLog.ToResponse(), nil
}

// List は監査ログ一覧を取得します
func (s *AuditLogService) List(ctx context.Context, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	auditLogs, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		s.logger.Error("Failed to fetch audit logs", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.AuditLogResponse, len(auditLogs))
	for i, log := range auditLogs {
		responses[i] = log.ToResponse()
	}

	return util.NewPaginatedResponse(responses, total, params), nil
}

// ListByUserID はユーザーIDで監査ログ一覧を取得します
func (s *AuditLogService) ListByUserID(ctx context.Context, userID uint, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	auditLogs, total, err := s.repo.FindByUserID(ctx, userID, params)
	if err != nil {
		s.logger.Error("Failed to fetch audit logs by user", zap.Uint("user_id", userID), zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.AuditLogResponse, len(auditLogs))
	for i, log := range auditLogs {
		responses[i] = log.ToResponse()
	}

	return util.NewPaginatedResponse(responses, total, params), nil
}

// ListByResource はリソースで監査ログ一覧を取得します
func (s *AuditLogService) ListByResource(ctx context.Context, resourceType, resourceID string, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	if resourceType == "" || resourceID == "" {
		return nil, util.NewBadRequestError(util.ErrCodeBadRequest, errors.New("resourceType and resourceID are required"))
	}

	auditLogs, total, err := s.repo.FindByResourceID(ctx, resourceType, resourceID, params)
	if err != nil {
		s.logger.Error("Failed to fetch audit logs by resource",
			zap.String("resource_type", resourceType),
			zap.String("resource_id", resourceID),
			zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.AuditLogResponse, len(auditLogs))
	for i, log := range auditLogs {
		responses[i] = log.ToResponse()
	}

	return util.NewPaginatedResponse(responses, total, params), nil
}

// ListByAction はアクションで監査ログ一覧を取得します
func (s *AuditLogService) ListByAction(ctx context.Context, action string, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	if action == "" {
		return nil, util.NewBadRequestError(util.ErrCodeBadRequest, errors.New("action is required"))
	}

	auditLogs, total, err := s.repo.FindByAction(ctx, action, params)
	if err != nil {
		s.logger.Error("Failed to fetch audit logs by action",
			zap.String("action", action),
			zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.AuditLogResponse, len(auditLogs))
	for i, log := range auditLogs {
		responses[i] = log.ToResponse()
	}

	return util.NewPaginatedResponse(responses, total, params), nil
}

// ListByDateRange は日付範囲で監査ログ一覧を取得します
func (s *AuditLogService) ListByDateRange(ctx context.Context, startDate, endDate time.Time, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	if startDate.After(endDate) {
		return nil, util.NewBadRequestError(util.ErrCodeBadRequest, errors.New("startDate must be before endDate"))
	}

	auditLogs, total, err := s.repo.FindByDateRange(ctx, startDate, endDate, params)
	if err != nil {
		s.logger.Error("Failed to fetch audit logs by date range",
			zap.Time("start_date", startDate),
			zap.Time("end_date", endDate),
			zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.AuditLogResponse, len(auditLogs))
	for i, log := range auditLogs {
		responses[i] = log.ToResponse()
	}

	return util.NewPaginatedResponse(responses, total, params), nil
}

// GetStatistics は監査ログの統計情報を取得します
type AuditStatistics struct {
	TotalLogs     int64            `json:"total_logs"`
	ByAction      map[string]int64 `json:"by_action"`
	ByStatus      map[string]int64 `json:"by_status"`
	SuccessRate   float64          `json:"success_rate"`
}

// GetStatistics は監査ログの統計情報を取得します
func (s *AuditLogService) GetStatistics(ctx context.Context) (*AuditStatistics, error) {
	// アクション別集計
	byAction, err := s.repo.CountByAction(ctx)
	if err != nil {
		s.logger.Error("Failed to get action statistics", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// ステータス別集計
	byStatus, err := s.repo.CountByStatus(ctx)
	if err != nil {
		s.logger.Error("Failed to get status statistics", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 成功率を計算
	successRate := 0.0
	total := int64(0)
	for _, count := range byStatus {
		total += count
	}
	if total > 0 {
		if successCount, ok := byStatus[model.AuditStatusSuccess]; ok {
			successRate = float64(successCount) / float64(total)
		}
	}

	return &AuditStatistics{
		TotalLogs:   total,
		ByAction:    byAction,
		ByStatus:    byStatus,
		SuccessRate: successRate,
	}, nil
}

// DeleteOldLogs は古い監査ログを削除します
func (s *AuditLogService) DeleteOldLogs(ctx context.Context, days int) error {
	if days < 1 {
		return util.NewBadRequestError(util.ErrCodeBadRequest, errors.New("days must be at least 1"))
	}

	if err := s.repo.DeleteOldLogs(ctx, days); err != nil {
		s.logger.Error("Failed to delete old audit logs", zap.Int("days", days), zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Old audit logs deleted", zap.Int("days", days))
	return nil
}

// validateCreateRequest はリクエストを検証します
func (s *AuditLogService) validateCreateRequest(req *model.CreateAuditLogRequest) error {
	if req.UserID == 0 {
		return errors.New("userID is required")
	}

	if req.Action == "" {
		return errors.New("action is required")
	}

	if req.ResourceType == "" {
		return errors.New("resourceType is required")
	}

	if req.ResourceID == "" {
		return errors.New("resourceID is required")
	}

	if req.Status == "" {
		return errors.New("status is required")
	}

	// アクションの有効性チェック
	validActions := map[string]bool{
		model.ActionCreate: true,
		model.ActionRead:   true,
		model.ActionUpdate: true,
		model.ActionDelete: true,
		model.ActionLogin:  true,
		model.ActionLogout: true,
	}
	if !validActions[req.Action] {
		return errors.New("invalid action")
	}

	// ステータスの有効性チェック
	validStatuses := map[string]bool{
		model.AuditStatusSuccess: true,
		model.AuditStatusFailed:  true,
	}
	if !validStatuses[req.Status] {
		return errors.New("invalid status")
	}

	return nil
}
