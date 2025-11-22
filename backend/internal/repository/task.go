package repository

import (
	"context"

	"github.com/varubogu/effisio/backend/internal/model"
	"gorm.io/gorm"
)

// TaskRepository はタスクのリポジトリインターフェースです
type TaskRepository interface {
	FindAll(ctx context.Context, filter *model.TaskFilter) ([]*model.Task, error)
	FindByID(ctx context.Context, id uint) (*model.Task, error)
	Create(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id uint) error
	CountByFilter(ctx context.Context, filter *model.TaskFilter) (int64, error)
}

type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository はTaskRepositoryを作成します
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

// FindAll はフィルタリングとページネーション付きでタスク一覧を取得します
func (r *taskRepository) FindAll(ctx context.Context, filter *model.TaskFilter) ([]*model.Task, error) {
	var tasks []*model.Task
	query := r.db.WithContext(ctx)

	// フィルタリング
	query = r.applyFilter(query, filter)

	// ソート
	sortColumn := "created_at"
	if filter.SortBy != "" {
		sortColumn = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	query = query.Order(sortColumn + " " + sortOrder)

	// ページネーション
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// リレーションをプリロード
	query = query.Preload("AssignedTo").Preload("CreatedBy").Preload("Organization")

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// FindByID はIDでタスクを取得します
func (r *taskRepository) FindByID(ctx context.Context, id uint) (*model.Task, error) {
	var task model.Task
	if err := r.db.WithContext(ctx).
		Preload("AssignedTo").
		Preload("CreatedBy").
		Preload("Organization").
		First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// Create はタスクを作成します
func (r *taskRepository) Create(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// Update はタスクを更新します
func (r *taskRepository) Update(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// Delete はタスクを削除します
func (r *taskRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Task{}, id).Error
}

// CountByFilter はフィルタに一致するタスク数を取得します
func (r *taskRepository) CountByFilter(ctx context.Context, filter *model.TaskFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Task{})

	// フィルタリング
	query = r.applyFilter(query, filter)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// applyFilter はクエリにフィルタを適用します
func (r *taskRepository) applyFilter(query *gorm.DB, filter *model.TaskFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}

	if filter.AssignedToID != nil {
		query = query.Where("assigned_to = ?", *filter.AssignedToID)
	}

	if filter.CreatedByID != nil {
		query = query.Where("created_by = ?", *filter.CreatedByID)
	}

	if filter.OrganizationID != nil {
		query = query.Where("organization_id = ?", *filter.OrganizationID)
	}

	if filter.DueBefore != nil {
		query = query.Where("due_date <= ?", *filter.DueBefore)
	}

	if filter.DueAfter != nil {
		query = query.Where("due_date >= ?", *filter.DueAfter)
	}

	return query
}
