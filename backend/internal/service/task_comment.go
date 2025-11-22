package service

import (
	"context"
	"errors"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TaskCommentService はタスクコメントのサービスです
type TaskCommentService struct {
	repo                repository.TaskCommentRepository
	taskRepo            repository.TaskRepository
	taskActivityService *TaskActivityService
	logger              *zap.Logger
}

// NewTaskCommentService はTaskCommentServiceを作成します
func NewTaskCommentService(
	repo repository.TaskCommentRepository,
	taskRepo repository.TaskRepository,
	taskActivityService *TaskActivityService,
	logger *zap.Logger,
) *TaskCommentService {
	return &TaskCommentService{
		repo:                repo,
		taskRepo:            taskRepo,
		taskActivityService: taskActivityService,
		logger:              logger,
	}
}

// ListByTaskID はタスクIDでコメント一覧を取得します
func (s *TaskCommentService) ListByTaskID(ctx context.Context, taskID uint) ([]*model.TaskCommentResponse, error) {
	// タスクの存在確認
	if _, err := s.taskRepo.FindByID(ctx, taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("TASK_001", errors.New("task not found"))
		}
		s.logger.Error("Failed to find task", zap.Error(err), zap.Uint("task_id", taskID))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	comments, err := s.repo.FindByTaskID(ctx, taskID)
	if err != nil {
		s.logger.Error("Failed to find task comments", zap.Error(err), zap.Uint("task_id", taskID))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.TaskCommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = comment.ToResponse()
	}

	return responses, nil
}

// GetByID はコメントIDでコメントを取得します
func (s *TaskCommentService) GetByID(ctx context.Context, id uint) (*model.TaskCommentResponse, error) {
	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("COMMENT_001", errors.New("comment not found"))
		}
		s.logger.Error("Failed to find comment", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return comment.ToResponse(), nil
}

// Create はコメントを作成します
func (s *TaskCommentService) Create(ctx context.Context, taskID uint, req *model.CreateTaskCommentRequest, userID uint) (*model.TaskCommentResponse, error) {
	// タスクの存在確認
	if _, err := s.taskRepo.FindByID(ctx, taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("TASK_001", errors.New("task not found"))
		}
		s.logger.Error("Failed to find task", zap.Error(err), zap.Uint("task_id", taskID))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// コメントの作成
	comment := &model.TaskComment{
		TaskID:  taskID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		s.logger.Error("Failed to create comment", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 作成後のコメントを取得（ユーザー情報を含む）
	createdComment, err := s.repo.FindByID(ctx, comment.ID)
	if err != nil {
		s.logger.Error("Failed to find created comment", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// アクティビティログの記録
	if err := s.taskActivityService.LogCommented(ctx, taskID, userID, comment.ID); err != nil {
		s.logger.Error("Failed to log commented activity", zap.Error(err))
		// アクティビティログの失敗はエラーとして返さない
	}

	return createdComment.ToResponse(), nil
}

// Update はコメントを更新します
func (s *TaskCommentService) Update(ctx context.Context, id uint, req *model.UpdateTaskCommentRequest, userID uint) (*model.TaskCommentResponse, error) {
	// 既存のコメントを取得
	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("COMMENT_001", errors.New("comment not found"))
		}
		s.logger.Error("Failed to find comment", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 権限チェック: コメント作成者のみが更新可能
	if comment.UserID != userID {
		return nil, util.NewForbiddenError("COMMENT_002", errors.New("only comment author can update"))
	}

	// コメント内容の更新
	comment.Content = req.Content

	if err := s.repo.Update(ctx, comment); err != nil {
		s.logger.Error("Failed to update comment", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 更新後のコメントを取得
	updatedComment, err := s.repo.FindByID(ctx, comment.ID)
	if err != nil {
		s.logger.Error("Failed to find updated comment", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return updatedComment.ToResponse(), nil
}

// Delete はコメントを削除します
func (s *TaskCommentService) Delete(ctx context.Context, id uint, userID uint) error {
	// 既存のコメントを取得
	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.NewNotFoundError("COMMENT_001", errors.New("comment not found"))
		}
		s.logger.Error("Failed to find comment", zap.Error(err), zap.Uint("id", id))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 権限チェック: コメント作成者のみが削除可能
	if comment.UserID != userID {
		return util.NewForbiddenError("COMMENT_002", errors.New("only comment author can delete"))
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete comment", zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return nil
}
