package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// AuditLogHandler は監査ログ関連のHTTPハンドラーを提供します
type AuditLogHandler struct {
	service *service.AuditLogService
	logger  *zap.Logger
}

// NewAuditLogHandler は新しいAuditLogHandlerを作成します
func NewAuditLogHandler(service *service.AuditLogService, logger *zap.Logger) *AuditLogHandler {
	return &AuditLogHandler{
		service: service,
		logger:  logger,
	}
}

// List godoc
// @Summary 監査ログ一覧取得
// @Description 監査ログを取得します（フィルタリング、ページネーション対応）
// @Tags audit_logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query int false "ユーザーID"
// @Param action query string false "アクション"
// @Param resource query string false "リソース"
// @Param resource_id query string false "リソースID"
// @Param start_date query string false "開始日時 (RFC3339形式)"
// @Param end_date query string false "終了日時 (RFC3339形式)"
// @Param page query int false "ページ番号"
// @Param per_page query int false "1ページあたりの件数"
// @Success 200 {object} util.PaginatedResponse
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /audit-logs [get]
func (h *AuditLogHandler) List(c *gin.Context) {
	// フィルタリングパラメータの取得
	filter := &repository.AuditLogFilter{}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(userID)
			filter.UserID = &uid
		}
	}

	if action := c.Query("action"); action != "" {
		filter.Action = action
	}

	if resource := c.Query("resource"); resource != "" {
		filter.Resource = resource
	}

	if resourceID := c.Query("resource_id"); resourceID != "" {
		filter.ResourceID = resourceID
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	// ページネーションパラメータの取得
	pagination := util.GetPaginationParams(c)

	// 監査ログ取得
	response, err := h.service.List(c.Request.Context(), filter, pagination)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(200, response)
}

// GetByID godoc
// @Summary 監査ログ詳細取得
// @Description IDで監査ログを取得します
// @Tags audit_logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "監査ログID"
// @Success 200 {object} util.Response{data=model.AuditLogResponse}
// @Failure 400 {object} util.Response "不正なリクエスト"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 404 {object} util.Response "監査ログが見つかりません"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /audit-logs/{id} [get]
func (h *AuditLogHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ValidationError(c, err)
		return
	}

	log, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, log)
}
