package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/pkg/util"
)

// RBACMiddleware はロールベースアクセス制御ミドルウェアを提供します
type RBACMiddleware struct {
	logger *zap.Logger
}

// NewRBACMiddleware は新しいRBACMiddlewareを作成します
func NewRBACMiddleware(logger *zap.Logger) *RBACMiddleware {
	return &RBACMiddleware{
		logger: logger,
	}
}

// RequirePermission は指定された権限を持つユーザーのみアクセスを許可します
// このミドルウェアは RequireAuth の後に使用する必要があります
func (m *RBACMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			m.logger.Warn("Permissions not found in context")
			util.Error(c, http.StatusUnauthorized, util.ErrCodeUnauthorized, "authentication required", nil)
			c.Abort()
			return
		}

		permList, ok := permissions.([]string)
		if !ok {
			m.logger.Error("Invalid permissions format in context")
			util.Error(c, http.StatusInternalServerError, util.ErrCodeInternalError, "internal server error", nil)
			c.Abort()
			return
		}

		// 権限をチェック
		if !contains(permList, permission) {
			userID, _ := c.Get("user_id")
			m.logger.Warn("Insufficient permissions",
				zap.Uint("user_id", userID.(uint)),
				zap.String("required_permission", permission),
				zap.Strings("user_permissions", permList),
			)
			util.Error(c, http.StatusForbidden, util.ErrCodeInsufficientPermission, "insufficient permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission は指定された権限のいずれかを持つユーザーのみアクセスを許可します
func (m *RBACMiddleware) RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPermissions, exists := c.Get("permissions")
		if !exists {
			m.logger.Warn("Permissions not found in context")
			util.Error(c, http.StatusUnauthorized, util.ErrCodeUnauthorized, "authentication required", nil)
			c.Abort()
			return
		}

		permList, ok := userPermissions.([]string)
		if !ok {
			m.logger.Error("Invalid permissions format in context")
			util.Error(c, http.StatusInternalServerError, util.ErrCodeInternalError, "internal server error", nil)
			c.Abort()
			return
		}

		// いずれかの権限を持っているかチェック
		hasPermission := false
		for _, perm := range permissions {
			if contains(permList, perm) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			userID, _ := c.Get("user_id")
			m.logger.Warn("Insufficient permissions",
				zap.Uint("user_id", userID.(uint)),
				zap.Strings("required_permissions", permissions),
				zap.Strings("user_permissions", permList),
			)
			util.Error(c, http.StatusForbidden, util.ErrCodeInsufficientPermission, "insufficient permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole は指定されたロールを持つユーザーのみアクセスを許可します
func (m *RBACMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			m.logger.Warn("Role not found in context")
			util.Error(c, http.StatusUnauthorized, util.ErrCodeUnauthorized, "authentication required", nil)
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			m.logger.Error("Invalid role format in context")
			util.Error(c, http.StatusInternalServerError, util.ErrCodeInternalError, "internal server error", nil)
			c.Abort()
			return
		}

		if roleStr != role {
			userID, _ := c.Get("user_id")
			m.logger.Warn("Insufficient role",
				zap.Uint("user_id", userID.(uint)),
				zap.String("required_role", role),
				zap.String("user_role", roleStr),
			)
			util.Error(c, http.StatusForbidden, util.ErrCodeInsufficientPermission, "insufficient permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole は指定されたロールのいずれかを持つユーザーのみアクセスを許可します
func (m *RBACMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			m.logger.Warn("Role not found in context")
			util.Error(c, http.StatusUnauthorized, util.ErrCodeUnauthorized, "authentication required", nil)
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			m.logger.Error("Invalid role format in context")
			util.Error(c, http.StatusInternalServerError, util.ErrCodeInternalError, "internal server error", nil)
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if roleStr == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			userID, _ := c.Get("user_id")
			m.logger.Warn("Insufficient role",
				zap.Uint("user_id", userID.(uint)),
				zap.Strings("required_roles", roles),
				zap.String("user_role", roleStr),
			)
			util.Error(c, http.StatusForbidden, util.ErrCodeInsufficientPermission, "insufficient permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// contains はスライスに指定された文字列が含まれているかチェックします
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
