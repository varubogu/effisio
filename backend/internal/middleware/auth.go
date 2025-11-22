package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/pkg/util"
)

// AuthMiddleware は認証ミドルウェアを提供します
type AuthMiddleware struct {
	jwtService *util.JWTService
	logger     *zap.Logger
}

// NewAuthMiddleware は新しいAuthMiddlewareを作成します
func NewAuthMiddleware(jwtService *util.JWTService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// RequireAuth は認証が必要なエンドポイントで使用するミドルウェアです
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization ヘッダーからトークンを取得
		tokenString, err := m.jwtService.ExtractTokenFromAuthHeader(c.GetHeader("Authorization"))
		if err != nil {
			m.logger.Warn("Missing or invalid authorization header", zap.Error(err))
			util.Error(c, http.StatusUnauthorized, util.ErrCodeUnauthorized, "authentication required", nil)
			c.Abort()
			return
		}

		// トークンを検証
		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			m.logger.Warn("Invalid access token", zap.Error(err))
			util.Error(c, http.StatusUnauthorized, util.ErrCodeInvalidToken, "invalid or expired token", nil)
			c.Abort()
			return
		}

		// ユーザー情報をコンテキストに設定
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)

		m.logger.Debug("User authenticated",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("role", claims.Role),
		)

		c.Next()
	}
}

// OptionalAuth は認証がオプションのエンドポイントで使用するミドルウェアです
// トークンがあれば検証してコンテキストに設定し、なければそのまま次に進みます
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 認証ヘッダーがない場合はスキップ
			c.Next()
			return
		}

		tokenString, err := m.jwtService.ExtractTokenFromAuthHeader(authHeader)
		if err != nil {
			// トークンの形式が不正な場合はスキップ
			c.Next()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			// トークンが無効な場合はスキップ
			c.Next()
			return
		}

		// ユーザー情報をコンテキストに設定
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}
