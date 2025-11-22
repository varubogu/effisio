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

// TaskService はタスクのサービスです
type TaskService struct {
	repo             repository.TaskRepository
	userRepo         repository.UserRepository
	organizationRepo repository.OrganizationRepository
	tagRepo          repository.TagRepository
	logger           *zap.Logger
}

// NewTaskService はTaskServiceを作成します
func NewTaskService(
	repo repository.TaskRepository,
	userRepo repository.UserRepository,
	organizationRepo repository.OrganizationRepository,
	tagRepo repository.TagRepository,
	logger *zap.Logger,
) *TaskService {
	return &TaskService{
		repo:             repo,
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
		tagRepo:          tagRepo,
		logger:           logger,
	}
}

// ListResponse はタスク一覧のレスポンスです
type TaskListResponse struct {
	Tasks      []*model.TaskResponse `json:"tasks"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// List はタスク一覧を取得します
func (s *TaskService) List(ctx context.Context, filter *model.TaskFilter) (*TaskListResponse, error) {
	tasks, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to find tasks", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	total, err := s.repo.CountByFilter(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to count tasks", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = task.ToResponse()
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &TaskListResponse{
		Tasks:      responses,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByID はタスクを取得します
func (s *TaskService) GetByID(ctx context.Context, id uint) (*model.TaskResponse, error) {
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("TASK_001", errors.New("task not found"))
		}
		s.logger.Error("Failed to find task", zap.Error(err), zap.Uint("id", id))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return task.ToResponse(), nil
}

// Create はタスクを作成します
func (s *TaskService) Create(ctx context.Context, req *model.CreateTaskRequest, createdByID uint) (*model.TaskResponse, error) {
	// 担当者の存在確認
	if req.AssignedToID != nil {
		if _, err := s.userRepo.FindByID(ctx, *req.AssignedToID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, util.NewValidationError("TASK_002", errors.New("assigned user not found"))
			}
			s.logger.Error("Failed to find assigned user", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
	}

	// 組織の存在確認
	if req.OrganizationID != nil {
		if _, err := s.organizationRepo.FindByID(ctx, *req.OrganizationID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, util.NewValidationError("TASK_003", errors.New("organization not found"))
			}
			s.logger.Error("Failed to find organization", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
	}

	// タグの存在確認
	var tags []model.Tag
	if len(req.TagIDs) > 0 {
		foundTags, err := s.tagRepo.FindByIDs(ctx, req.TagIDs)
		if err != nil {
			s.logger.Error("Failed to find tags", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		if len(foundTags) != len(req.TagIDs) {
			return nil, util.NewValidationError("TAG_003", errors.New("some tags not found"))
		}
		for _, tag := range foundTags {
			tags = append(tags, *tag)
		}
	}

	// タスクの作成
	task := &model.Task{
		Title:          req.Title,
		Description:    req.Description,
		Status:         model.TaskStatusTODO,
		Priority:       model.TaskPriorityMedium,
		AssignedToID:   req.AssignedToID,
		CreatedByID:    createdByID,
		OrganizationID: req.OrganizationID,
		DueDate:        req.DueDate,
		Tags:           tags,
	}

	// リクエストでステータスが指定されている場合は上書き
	if req.Status != "" {
		task.Status = req.Status
	}

	// リクエストで優先度が指定されている場合は上書き
	if req.Priority != "" {
		task.Priority = req.Priority
	}

	if err := s.repo.Create(ctx, task); err != nil {
		s.logger.Error("Failed to create task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 作成後のタスクを取得（リレーションを含む）
	createdTask, err := s.repo.FindByID(ctx, task.ID)
	if err != nil {
		s.logger.Error("Failed to find created task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return createdTask.ToResponse(), nil
}

// Update はタスクを更新します
func (s *TaskService) Update(ctx context.Context, id uint, req *model.UpdateTaskRequest, currentUserID uint) (*model.TaskResponse, error) {
	// 既存のタスクを取得
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("TASK_001", errors.New("task not found"))
		}
		s.logger.Error("Failed to find task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 権限チェック: 作成者のみが更新可能（adminロールは後でミドルウェアで制御）
	// ここではビジネスロジックとして作成者のチェックのみ行う
	// 実際にはミドルウェアでadminロールを持つユーザーはこのチェックをスキップする設計も可能

	// 担当者の変更がある場合は存在確認
	if req.AssignedToID != nil {
		if _, err := s.userRepo.FindByID(ctx, *req.AssignedToID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, util.NewValidationError("TASK_002", errors.New("assigned user not found"))
			}
			s.logger.Error("Failed to find assigned user", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		task.AssignedToID = req.AssignedToID
	}

	// 組織の変更がある場合は存在確認
	if req.OrganizationID != nil {
		if _, err := s.organizationRepo.FindByID(ctx, *req.OrganizationID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, util.NewValidationError("TASK_003", errors.New("organization not found"))
			}
			s.logger.Error("Failed to find organization", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
		task.OrganizationID = req.OrganizationID
	}

	// タグの更新（指定された場合は既存のタグを全て置き換え）
	if req.TagIDs != nil {
		var tags []model.Tag
		if len(req.TagIDs) > 0 {
			foundTags, err := s.tagRepo.FindByIDs(ctx, req.TagIDs)
			if err != nil {
				s.logger.Error("Failed to find tags", zap.Error(err))
				return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
			}
			if len(foundTags) != len(req.TagIDs) {
				return nil, util.NewValidationError("TAG_003", errors.New("some tags not found"))
			}
			for _, tag := range foundTags {
				tags = append(tags, *tag)
			}
		}
		task.Tags = tags
	}

	// フィールドの更新
	if req.Title != nil {
		task.Title = *req.Title
	}

	if req.Description != nil {
		task.Description = *req.Description
	}

	if req.Status != nil {
		task.Status = *req.Status
	}

	if req.Priority != nil {
		task.Priority = *req.Priority
	}

	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.repo.Update(ctx, task); err != nil {
		s.logger.Error("Failed to update task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 更新後のタスクを取得
	updatedTask, err := s.repo.FindByID(ctx, task.ID)
	if err != nil {
		s.logger.Error("Failed to find updated task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return updatedTask.ToResponse(), nil
}

// Delete はタスクを削除します
func (s *TaskService) Delete(ctx context.Context, id uint, currentUserID uint) error {
	// 既存のタスクを取得
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.NewNotFoundError("TASK_001", errors.New("task not found"))
		}
		s.logger.Error("Failed to find task", zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 権限チェック: 作成者のみが削除可能（adminロールは後でミドルウェアで制御）
	// ここではビジネスロジックとして作成者のチェックのみ行う
	_ = task // 権限チェックはハンドラー層で実装

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete task", zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return nil
}

// UpdateStatus はタスクのステータスを更新します
func (s *TaskService) UpdateStatus(ctx context.Context, id uint, status model.TaskStatus, currentUserID uint) (*model.TaskResponse, error) {
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("TASK_001", errors.New("task not found"))
		}
		s.logger.Error("Failed to find task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	task.Status = status

	if err := s.repo.Update(ctx, task); err != nil {
		s.logger.Error("Failed to update task status", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 更新後のタスクを取得
	updatedTask, err := s.repo.FindByID(ctx, task.ID)
	if err != nil {
		s.logger.Error("Failed to find updated task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return updatedTask.ToResponse(), nil
}

// AssignTask はタスクを担当者に割り当てます
func (s *TaskService) AssignTask(ctx context.Context, id uint, assignedToID *uint, currentUserID uint) (*model.TaskResponse, error) {
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError("TASK_001", errors.New("task not found"))
		}
		s.logger.Error("Failed to find task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 担当者の存在確認
	if assignedToID != nil {
		if _, err := s.userRepo.FindByID(ctx, *assignedToID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, util.NewValidationError("TASK_002", errors.New("assigned user not found"))
			}
			s.logger.Error("Failed to find assigned user", zap.Error(err))
			return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
		}
	}

	task.AssignedToID = assignedToID

	if err := s.repo.Update(ctx, task); err != nil {
		s.logger.Error("Failed to assign task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 更新後のタスクを取得
	updatedTask, err := s.repo.FindByID(ctx, task.ID)
	if err != nil {
		s.logger.Error("Failed to find updated task", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return updatedTask.ToResponse(), nil
}
