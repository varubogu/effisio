package repository

import (
	"context"

	"github.com/varubogu/effisio/backend/internal/model"
	"gorm.io/gorm"
)

// TaskActivityRepository はタスクアクティビティのリポジトリインターフェースです
type TaskActivityRepository interface {
	FindByTaskID(ctx context.Context, taskID uint) ([]*model.TaskActivity, error)
	FindByID(ctx context.Context, id uint) (*model.TaskActivity, error)
	Create(ctx context.Context, activity *model.TaskActivity) error
}

type taskActivityRepository struct {
	db *gorm.DB
}

// NewTaskActivityRepository はTaskActivityRepositoryを作成します
func NewTaskActivityRepository(db *gorm.DB) TaskActivityRepository {
	return &taskActivityRepository{db: db}
}

// FindByTaskID はタスクIDでアクティビティ一覧を取得します
func (r *taskActivityRepository) FindByTaskID(ctx context.Context, taskID uint) ([]*model.TaskActivity, error) {
	var activities []*model.TaskActivity
	if err := r.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at DESC").
		Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// FindByID はIDでアクティビティを取得します
func (r *taskActivityRepository) FindByID(ctx context.Context, id uint) (*model.TaskActivity, error) {
	var activity model.TaskActivity
	if err := r.db.WithContext(ctx).
		Preload("User").
		First(&activity, id).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

// Create はアクティビティを作成します
func (r *taskActivityRepository) Create(ctx context.Context, activity *model.TaskActivity) error {
	return r.db.WithContext(ctx).Create(activity).Error
}
