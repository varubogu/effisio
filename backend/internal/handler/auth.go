package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// AuthHandler は認証関連のHTTPハンドラーを提供します
type AuthHandler struct {
	authService     *service.AuthService
	auditLogService *service.AuditLogService
	logger          *zap.Logger
}

// NewAuthHandler は新しいAuthHandlerを作成します
func NewAuthHandler(
	authService *service.AuthService,
	auditLogService *service.AuditLogService,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService:     authService,
		auditLogService: auditLogService,
		logger:          logger,
	}
}

// Login godoc
// @Summary ログイン
// @Description ユーザー名とパスワードで認証し、アクセストークンとリフレッシュトークンを発行します
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "ログインリクエスト"
// @Success 200 {object} util.Response{data=service.LoginResponse} "ログイン成功"
// @Failure 400 {object} util.Response "バリデーションエラー"
// @Failure 401 {object} util.Response "認証エラー"
// @Failure 403 {object} util.Response "アカウントが無効"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, err)
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		// ログイン失敗を記録
		h.auditLogService.LogLoginFailed(req.Username, c.ClientIP(), c.Request.UserAgent())
		util.HandleError(c, err)
		return
	}

	// ログイン成功を記録
	h.auditLogService.LogLoginSuccess(response.User.ID, c.ClientIP(), c.Request.UserAgent())

	util.Success(c, response)
}

// RefreshToken godoc
// @Summary トークンをリフレッシュ
// @Description リフレッシュトークンを使用して新しいアクセストークンとリフレッシュトークンを発行します
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.RefreshTokenRequest true "リフレッシュトークンリクエスト"
// @Success 200 {object} util.Response{data=service.RefreshTokenResponse} "トークンリフレッシュ成功"
// @Failure 400 {object} util.Response "バリデーションエラー"
// @Failure 401 {object} util.Response "トークンが無効または期限切れ"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req service.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, err)
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, response)
}

// Logout godoc
// @Summary ログアウト
// @Description リフレッシュトークンを無効化してログアウトします
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.RefreshTokenRequest true "リフレッシュトークンリクエスト"
// @Success 200 {object} util.Response "ログアウト成功"
// @Failure 400 {object} util.Response "バリデーションエラー"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req service.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, err)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		util.HandleError(c, err)
		return
	}

	// ログアウトを記録（user_idはコンテキストから取得、無ければnull）
	if userID, exists := c.Get("user_id"); exists {
		h.auditLogService.LogLogout(userID.(uint), c.ClientIP(), c.Request.UserAgent())
	}

	util.Success(c, gin.H{"message": "logged out successfully"})
}

// LogoutAll godoc
// @Summary 全セッションからログアウト
// @Description ユーザーの全リフレッシュトークンを無効化します（認証が必要）
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} util.Response "全セッションからログアウト成功"
// @Failure 401 {object} util.Response "認証が必要"
// @Failure 500 {object} util.Response "サーバーエラー"
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	// ミドルウェアでセットされたユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		util.Error(c, http.StatusUnauthorized, util.ErrCodeUnauthorized, "authentication required", nil)
		return
	}

	if err := h.authService.LogoutAll(c.Request.Context(), userID.(uint)); err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, gin.H{"message": "all sessions logged out successfully"})
}
