package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
	"go.uber.org/zap"
)

// TaskCommentHandler はタスクコメントのハンドラーです
type TaskCommentHandler struct {
	service *service.TaskCommentService
	logger  *zap.Logger
}

// NewTaskCommentHandler はTaskCommentHandlerを作成します
func NewTaskCommentHandler(service *service.TaskCommentService, logger *zap.Logger) *TaskCommentHandler {
	return &TaskCommentHandler{
		service: service,
		logger:  logger,
	}
}

// ListByTaskID godoc
// @Summary タスクコメント一覧取得
// @Description タスクIDを指定してコメント一覧を取得します
// @Tags task-comments
// @Accept json
// @Produce json
// @Param task_id path int true "タスクID"
// @Security BearerAuth
// @Success 200 {array} model.TaskCommentResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{task_id}/comments [get]
func (h *TaskCommentHandler) ListByTaskID(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	response, err := h.service.ListByTaskID(c.Request.Context(), uint(taskID))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, response)
}

// GetByID godoc
// @Summary タスクコメント詳細取得
// @Description コメントIDを指定してコメントの詳細を取得します
// @Tags task-comments
// @Accept json
// @Produce json
// @Param task_id path int true "タスクID"
// @Param id path int true "コメントID"
// @Security BearerAuth
// @Success 200 {object} model.TaskCommentResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{task_id}/comments/{id} [get]
func (h *TaskCommentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	response, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, response)
}

// Create godoc
// @Summary タスクコメント作成
// @Description タスクに新しいコメントを作成します
// @Tags task-comments
// @Accept json
// @Produce json
// @Param task_id path int true "タスクID"
// @Param request body model.CreateTaskCommentRequest true "コメント作成リクエスト"
// @Security BearerAuth
// @Success 201 {object} model.TaskCommentResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{task_id}/comments [post]
func (h *TaskCommentHandler) Create(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	var req model.CreateTaskCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	// 現在のユーザーIDを取得
	currentUserID, exists := c.Get("user_id")
	if !exists {
		util.HandleError(c, util.NewUnauthorizedError(util.ErrCodeUnauthorized, nil))
		return
	}

	response, err := h.service.Create(c.Request.Context(), uint(taskID), &req, currentUserID.(uint))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// Update godoc
// @Summary タスクコメント更新
// @Description コメントを更新します（作成者のみ）
// @Tags task-comments
// @Accept json
// @Produce json
// @Param task_id path int true "タスクID"
// @Param id path int true "コメントID"
// @Param request body model.UpdateTaskCommentRequest true "コメント更新リクエスト"
// @Security BearerAuth
// @Success 200 {object} model.TaskCommentResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 403 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{task_id}/comments/{id} [put]
func (h *TaskCommentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	var req model.UpdateTaskCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	// 現在のユーザーIDを取得
	currentUserID, exists := c.Get("user_id")
	if !exists {
		util.HandleError(c, util.NewUnauthorizedError(util.ErrCodeUnauthorized, nil))
		return
	}

	response, err := h.service.Update(c.Request.Context(), uint(id), &req, currentUserID.(uint))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, response)
}

// Delete godoc
// @Summary タスクコメント削除
// @Description コメントを削除します（作成者のみ）
// @Tags task-comments
// @Accept json
// @Produce json
// @Param task_id path int true "タスクID"
// @Param id path int true "コメントID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 403 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{task_id}/comments/{id} [delete]
func (h *TaskCommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	// 現在のユーザーIDを取得
	currentUserID, exists := c.Get("user_id")
	if !exists {
		util.HandleError(c, util.NewUnauthorizedError(util.ErrCodeUnauthorized, nil))
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id), currentUserID.(uint)); err != nil {
		util.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
