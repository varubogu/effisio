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

// TagHandler はタグのハンドラーです
type TagHandler struct {
	service *service.TagService
	logger  *zap.Logger
}

// NewTagHandler はTagHandlerを作成します
func NewTagHandler(service *service.TagService, logger *zap.Logger) *TagHandler {
	return &TagHandler{
		service: service,
		logger:  logger,
	}
}

// List godoc
// @Summary タグ一覧取得
// @Description タグ一覧を取得します
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.TagResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tags [get]
func (h *TagHandler) List(c *gin.Context) {
	response, err := h.service.List(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, response)
}

// GetByID godoc
// @Summary タグ詳細取得
// @Description IDを指定してタグの詳細を取得します
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "タグID"
// @Security BearerAuth
// @Success 200 {object} model.TagResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tags/{id} [get]
func (h *TagHandler) GetByID(c *gin.Context) {
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
// @Summary タグ作成
// @Description 新しいタグを作成します（adminのみ）
// @Tags tags
// @Accept json
// @Produce json
// @Param request body model.CreateTagRequest true "タグ作成リクエスト"
// @Security BearerAuth
// @Success 201 {object} model.TagResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 403 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tags [post]
func (h *TagHandler) Create(c *gin.Context) {
	var req model.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	response, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// Update godoc
// @Summary タグ更新
// @Description タグを更新します（adminのみ）
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "タグID"
// @Param request body model.UpdateTagRequest true "タグ更新リクエスト"
// @Security BearerAuth
// @Success 200 {object} model.TagResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 403 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tags/{id} [put]
func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	var req model.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	response, err := h.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, response)
}

// Delete godoc
// @Summary タグ削除
// @Description タグを削除します（adminのみ）
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "タグID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 403 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /tags/{id} [delete]
func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.HandleError(c, util.NewValidationError(util.ErrCodeValidationFailed, err))
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		util.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
