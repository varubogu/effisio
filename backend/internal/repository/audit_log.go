package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// AuditLogRepository は監査ログのデータアクセスを提供します
type AuditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository は新しいAuditLogRepositoryを作成します
func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{
		db: db,
	}
}

// Create は監査ログを作成します
func (r *AuditLogRepository) Create(ctx context.Context, auditLog *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(auditLog).Error
}

// FindByID はIDで監査ログを取得します
func (r *AuditLogRepository) FindByID(ctx context.Context, id uint) (*model.AuditLog, error) {
	var auditLog model.AuditLog
	if err := r.db.WithContext(ctx).First(&auditLog, id).Error; err != nil {
		return nil, err
	}
	return &auditLog, nil
}

// FindAll は監査ログ一覧を取得します（ページネーション付き）
func (r *AuditLogRepository) FindAll(ctx context.Context, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	var auditLogs []*model.AuditLog
	var total int64

	// 件数を取得
	if err := r.db.WithContext(ctx).Model(&model.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーションでデータを取得
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(params.Offset).
		Limit(params.PerPage).
		Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}

	return auditLogs, total, nil
}

// FindByUserID はユーザーIDで監査ログを取得します
func (r *AuditLogRepository) FindByUserID(ctx context.Context, userID uint, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	var auditLogs []*model.AuditLog
	var total int64

	// 件数を取得
	if err := r.db.WithContext(ctx).Model(&model.AuditLog{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーションでデータを取得
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(params.Offset).
		Limit(params.PerPage).
		Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}

	return auditLogs, total, nil
}

// FindByResourceID はリソースIDで監査ログを取得します
func (r *AuditLogRepository) FindByResourceID(ctx context.Context, resourceType, resourceID string, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	var auditLogs []*model.AuditLog
	var total int64

	// 件数を取得
	if err := r.db.WithContext(ctx).Model(&model.AuditLog{}).
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーションでデータを取得
	if err := r.db.WithContext(ctx).
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("created_at DESC").
		Offset(params.Offset).
		Limit(params.PerPage).
		Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}

	return auditLogs, total, nil
}

// FindByAction はアクションで監査ログを取得します
func (r *AuditLogRepository) FindByAction(ctx context.Context, action string, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	var auditLogs []*model.AuditLog
	var total int64

	// 件数を取得
	if err := r.db.WithContext(ctx).Model(&model.AuditLog{}).
		Where("action = ?", action).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーションでデータを取得
	if err := r.db.WithContext(ctx).
		Where("action = ?", action).
		Order("created_at DESC").
		Offset(params.Offset).
		Limit(params.PerPage).
		Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}

	return auditLogs, total, nil
}

// FindByDateRange は日付範囲で監査ログを取得します
func (r *AuditLogRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	var auditLogs []*model.AuditLog
	var total int64

	// 件数を取得
	if err := r.db.WithContext(ctx).Model(&model.AuditLog{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーションでデータを取得
	if err := r.db.WithContext(ctx).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Offset(params.Offset).
		Limit(params.PerPage).
		Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}

	return auditLogs, total, nil
}

// DeleteOldLogs は古い監査ログを削除します（指定日数より古いもの）
func (r *AuditLogRepository) DeleteOldLogs(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return r.db.WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&model.AuditLog{}).Error
}

// CountByAction はアクション別の集計を取得します
func (r *AuditLogRepository) CountByAction(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		Action string
		Count  int64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.AuditLog{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Action] = r.Count
	}
	return counts, nil
}

// CountByStatus はステータス別の集計を取得します
func (r *AuditLogRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.AuditLog{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}
