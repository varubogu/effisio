package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// AuditLogHandler は監査ログ関連のHTTPハンドラを提供します
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

// List は監査ログ一覧を取得します
// @Summary 監査ログ一覧取得
// @Tags audit_logs
// @Security Bearer
// @Param page query int false "ページ番号（デフォルト: 1）"
// @Param per_page query int false "1ページあたりの件数（デフォルト: 10）"
// @Success 200 {object} util.PaginatedResponse
// @Failure 401 {object} util.ErrorResponse
// @Router /api/v1/audit-logs [get]
func (h *AuditLogHandler) List(c *gin.Context) {
	params := util.GetPaginationParams(c)

	response, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		h.logger.Error("Failed to list audit logs", zap.Error(err))
		util.HandleError(c, err)
		return
	}

	util.Paginated(c, response)
}

// GetByID はIDで監査ログを取得します
// @Summary 監査ログ詳細取得
// @Tags audit_logs
// @Security Bearer
// @Param id path int true "監査ログID"
// @Success 200 {object} model.AuditLogResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /api/v1/audit-logs/{id} [get]
func (h *AuditLogHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.Error(c, http.StatusBadRequest, util.ErrCodeBadRequest, "Invalid audit log ID", nil)
		return
	}

	auditLog, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, auditLog)
}

// ListByUserID はユーザーの監査ログを取得します
// @Summary ユーザーの監査ログ一覧取得
// @Tags audit_logs
// @Security Bearer
// @Param user_id path int true "ユーザーID"
// @Param page query int false "ページ番号（デフォルト: 1）"
// @Param per_page query int false "1ページあたりの件数（デフォルト: 10）"
// @Success 200 {object} util.PaginatedResponse
// @Failure 400 {object} util.ErrorResponse
// @Router /api/v1/audit-logs/user/{user_id} [get]
func (h *AuditLogHandler) ListByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		util.Error(c, http.StatusBadRequest, util.ErrCodeBadRequest, "Invalid user ID", nil)
		return
	}

	params := util.GetPaginationParams(c)

	response, err := h.service.ListByUserID(c.Request.Context(), uint(userID), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Paginated(c, response)
}

// ListByResource はリソースの監査ログを取得します
// @Summary リソースの監査ログ一覧取得
// @Tags audit_logs
// @Security Bearer
// @Param resource_type query string true "リソースタイプ（user, role, organization等）"
// @Param resource_id query string true "リソースID"
// @Param page query int false "ページ番号（デフォルト: 1）"
// @Param per_page query int false "1ページあたりの件数（デフォルト: 10）"
// @Success 200 {object} util.PaginatedResponse
// @Failure 400 {object} util.ErrorResponse
// @Router /api/v1/audit-logs/resource [get]
func (h *AuditLogHandler) ListByResource(c *gin.Context) {
	resourceType := c.Query("resource_type")
	resourceID := c.Query("resource_id")

	if resourceType == "" || resourceID == "" {
		util.Error(c, http.StatusBadRequest, util.ErrCodeBadRequest,
			"resource_type and resource_id query parameters are required", nil)
		return
	}

	params := util.GetPaginationParams(c)

	response, err := h.service.ListByResource(c.Request.Context(), resourceType, resourceID, params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Paginated(c, response)
}

// ListByAction はアクションで監査ログを取得します
// @Summary アクション別の監査ログ一覧取得
// @Tags audit_logs
// @Security Bearer
// @Param action query string true "アクション（create, read, update, delete, login, logout）"
// @Param page query int false "ページ番号（デフォルト: 1）"
// @Param per_page query int false "1ページあたりの件数（デフォルト: 10）"
// @Success 200 {object} util.PaginatedResponse
// @Failure 400 {object} util.ErrorResponse
// @Router /api/v1/audit-logs/action [get]
func (h *AuditLogHandler) ListByAction(c *gin.Context) {
	action := c.Query("action")
	if action == "" {
		util.Error(c, http.StatusBadRequest, util.ErrCodeBadRequest,
			"action query parameter is required", nil)
		return
	}

	params := util.GetPaginationParams(c)

	response, err := h.service.ListByAction(c.Request.Context(), action, params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Paginated(c, response)
}

// ListByDateRange は日付範囲で監査ログを取得します
// @Summary 日付範囲の監査ログ一覧取得
// @Tags audit_logs
// @Security Bearer
// @Param start_date query string true "開始日時（RFC3339形式）"
// @Param end_date query string true "終了日時（RFC3339形式）"
// @Param page query int false "ページ番号（デフォルト: 1）"
// @Param per_page query int false "1ページあたりの件数（デフォルト: 10）"
// @Success 200 {object} util.PaginatedResponse
// @Failure 400 {object} util.ErrorResponse
// @Router /api/v1/audit-logs/date-range [get]
func (h *AuditLogHandler) ListByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		util.Error(c, http.StatusBadRequest, util.ErrCodeBadRequest,
			"start_date and end_date query parameters are required", nil)
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		util.Error(c, http.StatusBadRequest, util.ErrCodeBadRequest,
			"Invalid start_date format (use RFC3339)", nil)
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		util.Error(c, http.StatusBadRequest, util.ErrCodeBadRequest,
			"Invalid end_date format (use RFC3339)", nil)
		return
	}

	params := util.GetPaginationParams(c)

	response, err := h.service.ListByDateRange(c.Request.Context(), startDate, endDate, params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Paginated(c, response)
}

// GetStatistics は監査ログの統計情報を取得します
// @Summary 監査ログ統計情報取得
// @Tags audit_logs
// @Security Bearer
// @Success 200 {object} service.AuditStatistics
// @Failure 401 {object} util.ErrorResponse
// @Router /api/v1/audit-logs/statistics [get]
func (h *AuditLogHandler) GetStatistics(c *gin.Context) {
	stats, err := h.service.GetStatistics(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get audit log statistics", zap.Error(err))
		util.HandleError(c, err)
		return
	}

	util.Success(c, stats)
}

// Create は監査ログを作成します（内部使用）
// @Summary 監査ログ作成
// @Tags audit_logs
// @Security Bearer
// @Param body body model.CreateAuditLogRequest true "監査ログデータ"
// @Success 201 {object} model.AuditLogResponse
// @Failure 400 {object} util.ErrorResponse
// @Router /api/v1/audit-logs [post]
func (h *AuditLogHandler) Create(c *gin.Context) {
	var req model.CreateAuditLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, util.ErrCodeBadRequest, "Invalid request body", nil)
		return
	}

	// クライアントIPを取得
	req.IPAddress = c.ClientIP()

	// ユーザーエージェントを取得
	req.UserAgent = c.GetHeader("User-Agent")

	auditLog, err := h.service.LogAction(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.SuccessResponse{
		Code:    http.StatusCreated,
		Message: "Audit log created",
		Data:    auditLog,
	})
}

// DeleteOldLogs は古い監査ログを削除します
// @Summary 古い監査ログ削除
// @Tags audit_logs
// @Security Bearer
// @Param days query int false "保持日数（デフォルト: 90）"
// @Success 200 {object} util.SuccessResponse
// @Failure 400 {object} util.ErrorResponse
// @Router /api/v1/audit-logs/delete-old [delete]
func (h *AuditLogHandler) DeleteOldLogs(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "90")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		util.Error(c, http.StatusBadRequest, util.ErrCodeBadRequest, "Invalid days value", nil)
		return
	}

	if err := h.service.DeleteOldLogs(c.Request.Context(), days); err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, gin.H{"message": "Old audit logs deleted"})
}
