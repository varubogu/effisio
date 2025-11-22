package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
	"go.uber.org/zap"
)

// TaskActivityHandler はタスクアクティビティのハンドラーです
type TaskActivityHandler struct {
	service *service.TaskActivityService
	logger  *zap.Logger
}

// NewTaskActivityHandler はTaskActivityHandlerを作成します
func NewTaskActivityHandler(service *service.TaskActivityService, logger *zap.Logger) *TaskActivityHandler {
	return &TaskActivityHandler{
		service: service,
		logger:  logger,
	}
}

// ListByTaskID godoc
// @Summary タスクアクティビティ一覧取得
// @Description タスクIDを指定してアクティビティ一覧を取得します
// @Tags task-activities
// @Accept json
// @Produce json
// @Param task_id path int true "タスクID"
// @Security BearerAuth
// @Success 200 {array} model.TaskActivityResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tasks/{task_id}/activities [get]
func (h *TaskActivityHandler) ListByTaskID(c *gin.Context) {
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
