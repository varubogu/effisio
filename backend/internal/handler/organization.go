package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// OrganizationHandler は組織関連のHTTPハンドラーを提供します
type OrganizationHandler struct {
	service *service.OrganizationService
	logger  *zap.Logger
}

// NewOrganizationHandler は新しいOrganizationHandlerを作成します
func NewOrganizationHandler(service *service.OrganizationService, logger *zap.Logger) *OrganizationHandler {
	return &OrganizationHandler{
		service: service,
		logger:  logger,
	}
}

// List godoc
// @Summary 組織一覧取得
// @Description 全ての組織を取得します
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} util.Response{data=[]model.OrganizationResponse}
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /organizations [get]
func (h *OrganizationHandler) List(c *gin.Context) {
	organizations, err := h.service.List(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, organizations)
}

// GetTree godoc
// @Summary 組織ツリー取得
// @Description 組織の階層構造をツリー形式で取得します
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} util.Response{data=model.OrganizationTreeResponse}
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /organizations/tree [get]
func (h *OrganizationHandler) GetTree(c *gin.Context) {
	tree, err := h.service.GetTree(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, tree)
}

// GetByID godoc
// @Summary 組織詳細取得
// @Description IDで組織を取得します
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "組織ID"
// @Success 200 {object} util.Response{data=model.OrganizationResponse}
// @Failure 400 {object} util.Response "不正なリクエスト"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 404 {object} util.Response "組織が見つかりません"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /organizations/{id} [get]
func (h *OrganizationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ValidationError(c, err)
		return
	}

	org, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, org)
}

// GetChildren godoc
// @Summary 子組織取得
// @Description 指定組織の直接の子組織を取得します
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "親組織ID"
// @Success 200 {object} util.Response{data=[]model.OrganizationResponse}
// @Failure 400 {object} util.Response "不正なリクエスト"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 404 {object} util.Response "組織が見つかりません"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /organizations/{id}/children [get]
func (h *OrganizationHandler) GetChildren(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ValidationError(c, err)
		return
	}

	children, err := h.service.GetChildren(c.Request.Context(), uint(id))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, children)
}

// Create godoc
// @Summary 組織作成
// @Description 新しい組織を作成します
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateOrganizationRequest true "組織作成リクエスト"
// @Success 201 {object} util.Response{data=model.OrganizationResponse}
// @Failure 400 {object} util.Response "バリデーションエラー"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 409 {object} util.Response "コードが既に存在"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /organizations [post]
func (h *OrganizationHandler) Create(c *gin.Context) {
	var req model.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, err)
		return
	}

	org, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Created(c, org)
}

// Update godoc
// @Summary 組織更新
// @Description 組織情報を更新します
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "組織ID"
// @Param request body model.UpdateOrganizationRequest true "組織更新リクエスト"
// @Success 200 {object} util.Response{data=model.OrganizationResponse}
// @Failure 400 {object} util.Response "バリデーションエラー"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 404 {object} util.Response "組織が見つかりません"
// @Failure 409 {object} util.Response "コードが既に存在"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /organizations/{id} [put]
func (h *OrganizationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ValidationError(c, err)
		return
	}

	var req model.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, err)
		return
	}

	org, err := h.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, org)
}

// Delete godoc
// @Summary 組織削除
// @Description 組織を削除します（子組織がある場合は削除できません）
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "組織ID"
// @Success 204 "削除成功"
// @Failure 400 {object} util.Response "不正なリクエスト"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 403 {object} util.Response "権限不足"
// @Failure 404 {object} util.Response "組織が見つかりません"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /organizations/{id} [delete]
func (h *OrganizationHandler) Delete(c *gin.Context) {
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
