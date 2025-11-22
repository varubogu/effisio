package repository

import (
	"context"

	"github.com/varubogu/effisio/backend/internal/model"
	"gorm.io/gorm"
)

// TaskCommentRepository はタスクコメントのリポジトリインターフェースです
type TaskCommentRepository interface {
	FindByTaskID(ctx context.Context, taskID uint) ([]*model.TaskComment, error)
	FindByID(ctx context.Context, id uint) (*model.TaskComment, error)
	Create(ctx context.Context, comment *model.TaskComment) error
	Update(ctx context.Context, comment *model.TaskComment) error
	Delete(ctx context.Context, id uint) error
}

type taskCommentRepository struct {
	db *gorm.DB
}

// NewTaskCommentRepository はTaskCommentRepositoryを作成します
func NewTaskCommentRepository(db *gorm.DB) TaskCommentRepository {
	return &taskCommentRepository{db: db}
}

// FindByTaskID はタスクIDでコメント一覧を取得します
func (r *taskCommentRepository) FindByTaskID(ctx context.Context, taskID uint) ([]*model.TaskComment, error) {
	var comments []*model.TaskComment
	if err := r.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByID はIDでコメントを取得します
func (r *taskCommentRepository) FindByID(ctx context.Context, id uint) (*model.TaskComment, error) {
	var comment model.TaskComment
	if err := r.db.WithContext(ctx).
		Preload("User").
		First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Create はコメントを作成します
func (r *taskCommentRepository) Create(ctx context.Context, comment *model.TaskComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// Update はコメントを更新します
func (r *taskCommentRepository) Update(ctx context.Context, comment *model.TaskComment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// Delete はコメントを削除します
func (r *taskCommentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.TaskComment{}, id).Error
}
