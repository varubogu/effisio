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

// TaskHandler はタスクのハンドラーです
type TaskHandler struct {
	service *service.TaskService
	logger  *zap.Logger
}

// NewTaskHandler はTaskHandlerを作成します
func NewTaskHandler(service *service.TaskService, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		logger:  logger,
	}
}

// List godoc
// @Summary タスク一覧取得
// @Description タスク一覧を取得します。フィルタリングとページネーションに対応しています。
// @Tags tasks
// @Accept json
// @Produce json
// @Param status query string false "ステータスでフィルタ" Enums(TODO, IN_PROGRESS, IN_REVIEW, DONE, CANCELLED)
// @Param priority query string false "優先度でフィルタ" Enums(LOW, MEDIUM, HIGH, URGENT)
// @Param assigned_to_id query int false "担当者IDでフィルタ"
// @Param created_by_id query int false "作成者IDでフィルタ"
// @Param organization_id query int false "組織IDでフィルタ"
// @Param due_before query string false "期限前でフィルタ（RFC3339形式）"
// @Param due_after query string false "期限後でフィルタ（RFC3339形式）"
// @Param page query int false "ページ番号" default(1)
// @Param page_size query int false "ページサイズ" default(20)
// @Param sort_by query string false "ソートカラム" default(created_at)
// @Param sort_order query string false "ソート順" Enums(asc, desc) default(desc)
// @Security BearerAuth
// @Success 200 {object} service.TaskListResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	var filter model.TaskFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	// デフォルト値の設定
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 20
	}
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	response, err := h.service.List(c.Request.Context(), &filter)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, response)
}

// GetByID godoc
// @Summary タスク詳細取得
// @Description IDを指定してタスクの詳細を取得します
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Security BearerAuth
// @Success 200 {object} model.TaskResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
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
// @Summary タスク作成
// @Description 新しいタスクを作成します
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body model.CreateTaskRequest true "タスク作成リクエスト"
// @Security BearerAuth
// @Success 201 {object} model.TaskResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	var req model.CreateTaskRequest
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

	response, err := h.service.Create(c.Request.Context(), &req, currentUserID.(uint))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// Update godoc
// @Summary タスク更新
// @Description タスクを更新します
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Param request body model.UpdateTaskRequest true "タスク更新リクエスト"
// @Security BearerAuth
// @Success 200 {object} model.TaskResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 403 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	var req model.UpdateTaskRequest
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
// @Summary タスク削除
// @Description タスクを削除します
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 403 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
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

// UpdateStatusRequest はステータス更新リクエストです
type UpdateStatusRequest struct {
	Status model.TaskStatus `json:"status" binding:"required,oneof=TODO IN_PROGRESS IN_REVIEW DONE CANCELLED"`
}

// UpdateStatus godoc
// @Summary タスクステータス更新
// @Description タスクのステータスを更新します
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Param request body UpdateStatusRequest true "ステータス更新リクエスト"
// @Security BearerAuth
// @Success 200 {object} model.TaskResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{id}/status [patch]
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	var req UpdateStatusRequest
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

	response, err := h.service.UpdateStatus(c.Request.Context(), uint(id), req.Status, currentUserID.(uint))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, response)
}

// AssignTaskRequest はタスク割り当てリクエストです
type AssignTaskRequest struct {
	AssignedToID *uint `json:"assigned_to_id"`
}

// AssignTask godoc
// @Summary タスク割り当て
// @Description タスクを担当者に割り当てます
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "タスクID"
// @Param request body AssignTaskRequest true "タスク割り当てリクエスト"
// @Security BearerAuth
// @Success 200 {object} model.TaskResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{id}/assign [patch]
func (h *TaskHandler) AssignTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	var req AssignTaskRequest
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

	response, err := h.service.AssignTask(c.Request.Context(), uint(id), req.AssignedToID, currentUserID.(uint))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, response)
}
