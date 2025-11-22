package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// PermissionHandler は権限関連のHTTPハンドラーを提供します
type PermissionHandler struct {
	service *service.PermissionService
	logger  *zap.Logger
}

// NewPermissionHandler は新しいPermissionHandlerを作成します
func NewPermissionHandler(service *service.PermissionService, logger *zap.Logger) *PermissionHandler {
	return &PermissionHandler{
		service: service,
		logger:  logger,
	}
}

// List godoc
// @Summary 権限一覧取得
// @Description 全ての権限を取得します（admin のみ）
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} util.Response{data=[]model.PermissionResponse}
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /permissions [get]
func (h *PermissionHandler) List(c *gin.Context) {
	permissions, err := h.service.List(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, permissions)
}

// GetByID godoc
// @Summary 権限詳細取得
// @Description IDで権限を取得します（admin のみ）
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "権限ID"
// @Success 200 {object} util.Response{data=model.PermissionResponse}
// @Failure 400 {object} util.Response "不正なリクエスト"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 404 {object} util.Response "権限が見つかりません"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /permissions/{id} [get]
func (h *PermissionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ValidationError(c, err)
		return
	}

	permission, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, permission)
}

// Create godoc
// @Summary 権限作成
// @Description 新しい権限を作成します（admin のみ）
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreatePermissionRequest true "権限作成リクエスト"
// @Success 201 {object} util.Response{data=model.PermissionResponse}
// @Failure 400 {object} util.Response "バリデーションエラー"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 409 {object} util.Response "権限名が既に存在"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /permissions [post]
func (h *PermissionHandler) Create(c *gin.Context) {
	var req model.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, err)
		return
	}

	permission, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Created(c, permission)
}

// Update godoc
// @Summary 権限更新
// @Description 権限情報を更新します（admin のみ）
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "権限ID"
// @Param request body model.UpdatePermissionRequest true "権限更新リクエスト"
// @Success 200 {object} util.Response{data=model.PermissionResponse}
// @Failure 400 {object} util.Response "バリデーションエラー"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 404 {object} util.Response "権限が見つかりません"
// @Failure 409 {object} util.Response "権限名が既に存在"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /permissions/{id} [put]
func (h *PermissionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ValidationError(c, err)
		return
	}

	var req model.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, err)
		return
	}

	permission, err := h.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, permission)
}

// Delete godoc
// @Summary 権限削除
// @Description 権限を削除します（admin のみ）
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "権限ID"
// @Success 204 "削除成功"
// @Failure 400 {object} util.Response "不正なリクエスト"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 404 {object} util.Response "権限が見つかりません"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /permissions/{id} [delete]
func (h *PermissionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ValidationError(c, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		util.HandleError(c, err)
		return
	}

	util.NoContent(c)
}
