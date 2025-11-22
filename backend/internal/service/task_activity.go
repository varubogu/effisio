package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TaskActivityService はタスクアクティビティのサービスです
type TaskActivityService struct {
	repo     repository.TaskActivityRepository
	taskRepo repository.TaskRepository
	logger   *zap.Logger
}

// NewTaskActivityService はTaskActivityServiceを作成します
func NewTaskActivityService(
	repo repository.TaskActivityRepository,
	taskRepo repository.TaskRepository,
	logger *zap.Logger,
) *TaskActivityService {
	return &TaskActivityService{
		repo:     repo,
		taskRepo: taskRepo,
		logger:   logger,
	}
}

// ListByTaskID はタスクIDでアクティビティ一覧を取得します
func (s *TaskActivityService) ListByTaskID(ctx context.Context, taskID uint) ([]*model.TaskActivityResponse, error) {
	// タスクの存在確認
	if _, err := s.taskRepo.FindByID(ctx, taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("TASK_001", errors.New("task not found"))
		}
		s.logger.Error("Failed to find task", zap.Error(err), zap.Uint("task_id", taskID))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	activities, err := s.repo.FindByTaskID(ctx, taskID)
	if err != nil {
		s.logger.Error("Failed to find task activities", zap.Error(err), zap.Uint("task_id", taskID))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.TaskActivityResponse, len(activities))
	for i, activity := range activities {
		responses[i] = activity.ToResponse()
	}

	return responses, nil
}

// LogActivity はアクティビティを記録します
func (s *TaskActivityService) LogActivity(ctx context.Context, activity *model.TaskActivity) error {
	if err := s.repo.Create(ctx, activity); err != nil {
		s.logger.Error("Failed to create task activity", zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}
	return nil
}

// LogCreated はタスク作成アクティビティを記録します
func (s *TaskActivityService) LogCreated(ctx context.Context, taskID uint, userID uint) error {
	activity := &model.TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: model.ActivityTypeCreated,
	}
	return s.LogActivity(ctx, activity)
}

// LogFieldChange はフィールド変更アクティビティを記録します
func (s *TaskActivityService) LogFieldChange(ctx context.Context, taskID uint, userID uint, fieldName string, oldValue string, newValue string) error {
	activity := &model.TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: model.ActivityTypeUpdated,
		FieldName:    fieldName,
		OldValue:     oldValue,
		NewValue:     newValue,
	}
	return s.LogActivity(ctx, activity)
}

// LogStatusChange はステータス変更アクティビティを記録します
func (s *TaskActivityService) LogStatusChange(ctx context.Context, taskID uint, userID uint, oldStatus string, newStatus string) error {
	activity := &model.TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: model.ActivityTypeStatusChanged,
		FieldName:    "status",
		OldValue:     oldStatus,
		NewValue:     newStatus,
	}
	return s.LogActivity(ctx, activity)
}

// LogAssigned は担当者変更アクティビティを記録します
func (s *TaskActivityService) LogAssigned(ctx context.Context, taskID uint, userID uint, oldAssignedToID *uint, newAssignedToID *uint) error {
	oldValue := "未割当"
	if oldAssignedToID != nil {
		oldValue = fmt.Sprintf("User ID: %d", *oldAssignedToID)
	}

	newValue := "未割当"
	if newAssignedToID != nil {
		newValue = fmt.Sprintf("User ID: %d", *newAssignedToID)
	}

	activity := &model.TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: model.ActivityTypeAssigned,
		FieldName:    "assigned_to",
		OldValue:     oldValue,
		NewValue:     newValue,
	}
	return s.LogActivity(ctx, activity)
}

// LogCommented はコメント追加アクティビティを記録します
func (s *TaskActivityService) LogCommented(ctx context.Context, taskID uint, userID uint, commentID uint) error {
	metadata := model.ActivityMetadata{
		"comment_id": commentID,
	}

	activity := &model.TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: model.ActivityTypeCommented,
		Metadata:     metadata,
	}
	return s.LogActivity(ctx, activity)
}

// LogTagAdded はタグ追加アクティビティを記録します
func (s *TaskActivityService) LogTagAdded(ctx context.Context, taskID uint, userID uint, tagID uint, tagName string) error {
	metadata := model.ActivityMetadata{
		"tag_id":   tagID,
		"tag_name": tagName,
	}

	activity := &model.TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: model.ActivityTypeTagAdded,
		Metadata:     metadata,
	}
	return s.LogActivity(ctx, activity)
}

// LogTagRemoved はタグ削除アクティビティを記録します
func (s *TaskActivityService) LogTagRemoved(ctx context.Context, taskID uint, userID uint, tagID uint, tagName string) error {
	metadata := model.ActivityMetadata{
		"tag_id":   tagID,
		"tag_name": tagName,
	}

	activity := &model.TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: model.ActivityTypeTagRemoved,
		Metadata:     metadata,
	}
	return s.LogActivity(ctx, activity)
}
