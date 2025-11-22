package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// AuditLogFilter は監査ログのフィルタリング条件です
type AuditLogFilter struct {
	UserID     *uint
	Action     string
	Resource   string
	ResourceID string
	StartDate  *time.Time
	EndDate    *time.Time
}

// AuditLogRepository はAuditLogのリポジトリインターフェースです
type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	FindAll(ctx context.Context, filter *AuditLogFilter, pagination *util.PaginationParams) ([]*model.AuditLog, int64, error)
	FindByID(ctx context.Context, id uint64) (*model.AuditLog, error)
	DeleteOlderThan(ctx context.Context, date time.Time) (int64, error)
}

// auditLogRepository はAuditLogRepositoryの実装です
type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository は新しいAuditLogRepositoryを作成します
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Create は新しい監査ログを作成します
func (r *auditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// FindAll は監査ログを検索します（フィルタリングとページネーション対応）
func (r *auditLogRepository) FindAll(
	ctx context.Context,
	filter *AuditLogFilter,
	pagination *util.PaginationParams,
) ([]*model.AuditLog, int64, error) {
	var logs []*model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	// フィルタリング
	if filter != nil {
		if filter.UserID != nil {
			query = query.Where("user_id = ?", *filter.UserID)
		}
		if filter.Action != "" {
			query = query.Where("action = ?", filter.Action)
		}
		if filter.Resource != "" {
			query = query.Where("resource = ?", filter.Resource)
		}
		if filter.ResourceID != "" {
			query = query.Where("resource_id = ?", filter.ResourceID)
		}
		if filter.StartDate != nil {
			query = query.Where("created_at >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			query = query.Where("created_at <= ?", *filter.EndDate)
		}
	}

	// 件数取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション適用
	if pagination != nil {
		query = query.Offset(pagination.Offset).Limit(pagination.PerPage)
	}

	// ソート（最新順）
	query = query.Order("created_at DESC")

	// データ取得
	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// FindByID はIDで監査ログを取得します
func (r *auditLogRepository) FindByID(ctx context.Context, id uint64) (*model.AuditLog, error) {
	var log model.AuditLog
	if err := r.db.WithContext(ctx).First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// DeleteOlderThan は指定日時より古い監査ログを削除します
func (r *auditLogRepository) DeleteOlderThan(ctx context.Context, date time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("created_at < ?", date).Delete(&model.AuditLog{})
	return result.RowsAffected, result.Error
}
