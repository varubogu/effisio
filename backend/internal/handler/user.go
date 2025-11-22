package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// UserHandler はユーザー関連のハンドラーです
type UserHandler struct {
	service         *service.UserService
	auditLogService *service.AuditLogService
	logger          *zap.Logger
}

// NewUserHandler は新しいUserHandlerを作成します
func NewUserHandler(
	service *service.UserService,
	auditLogService *service.AuditLogService,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		service:         service,
		auditLogService: auditLogService,
		logger:          logger,
	}
}

// List はユーザー一覧を取得します
// @Summary ユーザー一覧取得
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param per_page query int false "1ページあたりの件数" default(10)
// @Success 200 {object} util.PaginatedResponse
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	params := util.GetPaginationParams(c)
	result, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Paginated(c, result)
}

// GetByID はIDでユーザーを取得します
// @Summary ユーザー詳細取得
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Success 200 {object} model.UserResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.Error(c, 400, util.ErrCodeInvalidParameter, "Invalid user ID", nil)
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, gin.H{"user": user})
}

// Create は新しいユーザーを作成します
// @Summary ユーザー作成
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "ユーザー作成リクエスト"
// @Success 201 {object} model.UserResponse
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, util.ParseValidationErrors(err))
		return
	}

	user, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// ユーザー作成を記録
	if currentUserID, exists := c.Get("user_id"); exists {
		uid := currentUserID.(uint)
		h.auditLogService.LogAsync(&service.LogEntry{
			UserID:     &uid,
			Action:     model.ActionUserCreate,
			Resource:   model.ResourceUser,
			ResourceID: strconv.Itoa(int(user.ID)),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Changes: map[string]interface{}{
				"after": map[string]interface{}{
					"username": user.Username,
					"email":    user.Email,
					"role":     user.Role,
				},
			},
		})
	}

	util.Created(c, gin.H{"user": user})
}

// Update はユーザー情報を更新します
// @Summary ユーザー更新
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Param request body model.UpdateUserRequest true "ユーザー更新リクエスト"
// @Success 200 {object} model.UserResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.Error(c, 400, util.ErrCodeInvalidParameter, "Invalid user ID", nil)
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, util.ParseValidationErrors(err))
		return
	}

	user, err := h.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, gin.H{"user": user})
}

// Delete はユーザーを削除します
// @Summary ユーザー削除
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Success 204
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.Error(c, 400, util.ErrCodeInvalidParameter, "Invalid user ID", nil)
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		util.HandleError(c, err)
		return
	}

	util.NoContent(c)
}
