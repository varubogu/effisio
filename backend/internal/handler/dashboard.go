package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// DashboardHandler はダッシュボード関連のハンドラーです
type DashboardHandler struct {
	service *service.DashboardService
	logger  *zap.Logger
}

// NewDashboardHandler は新しいDashboardHandlerを作成します
func NewDashboardHandler(service *service.DashboardService, logger *zap.Logger) *DashboardHandler {
	return &DashboardHandler{
		service: service,
		logger:  logger,
	}
}

// Overview はダッシュボード概要を取得します
// @Summary ダッシュボード概要取得
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} util.Response{data=service.DashboardOverview}
// @Failure 500 {object} util.ErrorResponse
// @Router /api/v1/dashboard/overview [get]
func (h *DashboardHandler) Overview(c *gin.Context) {
	overview, err := h.service.GetOverview(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get dashboard overview", zap.Error(err))
		util.HandleError(c, err)
		return
	}

	util.Success(c, overview)
}
