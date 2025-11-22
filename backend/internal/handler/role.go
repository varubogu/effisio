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

// RoleHandler はロール関連のHTTPハンドラーを提供します
type RoleHandler struct {
	service *service.RoleService
	logger  *zap.Logger
}

// NewRoleHandler は新しいRoleHandlerを作成します
func NewRoleHandler(service *service.RoleService, logger *zap.Logger) *RoleHandler {
	return &RoleHandler{
		service: service,
		logger:  logger,
	}
}

// List godoc
// @Summary ロール一覧取得
// @Description 全てのロールを取得します（admin のみ）
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param with_permissions query boolean false "権限情報を含める"
// @Success 200 {object} util.Response{data=[]model.RoleResponse}
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	withPermissions := c.DefaultQuery("with_permissions", "false") == "true"

	roles, err := h.service.List(c.Request.Context(), withPermissions)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, roles)
}

// GetByID godoc
// @Summary ロール詳細取得
// @Description IDでロールを取得します（admin のみ）
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ロールID"
// @Param with_permissions query boolean false "権限情報を含める"
// @Success 200 {object} util.Response{data=model.RoleResponse}
// @Failure 400 {object} util.Response "不正なリクエスト"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 404 {object} util.Response "ロールが見つかりません"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /roles/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ValidationError(c, err)
		return
	}

	withPermissions := c.DefaultQuery("with_permissions", "true") == "true"

	role, err := h.service.GetByID(c.Request.Context(), uint(id), withPermissions)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, role)
}

// Create godoc
// @Summary ロール作成
// @Description 新しいロールを作成します（admin のみ）
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateRoleRequest true "ロール作成リクエスト"
// @Success 201 {object} util.Response{data=model.RoleResponse}
// @Failure 400 {object} util.Response "バリデーションエラー"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 409 {object} util.Response "ロール名が既に存在"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req model.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, err)
		return
	}

	role, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Created(c, role)
}

// Update godoc
// @Summary ロール更新
// @Description ロール情報を更新します（admin のみ）
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ロールID"
// @Param request body model.UpdateRoleRequest true "ロール更新リクエスト"
// @Success 200 {object} util.Response{data=model.RoleResponse}
// @Failure 400 {object} util.Response "バリデーションエラー"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 404 {object} util.Response "ロールが見つかりません"
// @Failure 409 {object} util.Response "ロール名が既に存在"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ValidationError(c, err)
		return
	}

	var req model.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, err)
		return
	}

	role, err := h.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, role)
}

// Delete godoc
// @Summary ロール削除
// @Description ロールを削除します（admin のみ）
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ロールID"
// @Success 204 "削除成功"
// @Failure 400 {object} util.Response "不正なリクエスト"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 404 {object} util.Response "ロールが見つかりません"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
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
