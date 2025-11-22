package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/varubogu/effisio/backend/pkg/util"
	"go.uber.org/zap"
)

func getTestJWTService() *util.JWTService {
	return util.NewJWTService("test-secret-key", 15*time.Minute, 7*24*time.Hour)
}

func getTestLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func TestAuthMiddleware_RequireAuth_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	jwtService := getTestJWTService()
	authMiddleware := NewAuthMiddleware(jwtService, getTestLogger())

	// Generate a valid token
	token, err := jwtService.GenerateAccessToken(1, "testuser", "admin", []string{"users:read"})
	require.NoError(t, err)

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, uint(1), userID)

		username, exists := c.Get("username")
		assert.True(t, exists)
		assert.Equal(t, "testuser", username)

		role, exists := c.Get("role")
		assert.True(t, exists)
		assert.Equal(t, "admin", role)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_RequireAuth_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	jwtService := getTestJWTService()
	authMiddleware := NewAuthMiddleware(jwtService, getTestLogger())

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_RequireAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	jwtService := getTestJWTService()
	authMiddleware := NewAuthMiddleware(jwtService, getTestLogger())

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_RequireAuth_MalformedAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	jwtService := getTestJWTService()
	authMiddleware := NewAuthMiddleware(jwtService, getTestLogger())

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	testCases := []string{
		"InvalidToken",           // Missing Bearer prefix
		"Bearer",                 // Missing token
		"Basic dXNlcm5hbWU6cGFzcw==", // Wrong prefix
	}

	for _, authHeader := range testCases {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Failed for auth header: "+authHeader)
	}
}

func TestAuthMiddleware_OptionalAuth_WithValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	jwtService := getTestJWTService()
	authMiddleware := NewAuthMiddleware(jwtService, getTestLogger())

	// Generate a valid token
	token, err := jwtService.GenerateAccessToken(1, "testuser", "user", []string{"tasks:read"})
	require.NoError(t, err)

	router.GET("/optional", authMiddleware.OptionalAuth(), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, uint(1), userID)

		c.JSON(http.StatusOK, gin.H{"authenticated": true})
	})

	req := httptest.NewRequest("GET", "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_OptionalAuth_WithoutToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	jwtService := getTestJWTService()
	authMiddleware := NewAuthMiddleware(jwtService, getTestLogger())

	router.GET("/optional", authMiddleware.OptionalAuth(), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.False(t, exists) // Should not be set

		c.JSON(http.StatusOK, gin.H{"authenticated": false})
	})

	req := httptest.NewRequest("GET", "/optional", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_OptionalAuth_WithInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	jwtService := getTestJWTService()
	authMiddleware := NewAuthMiddleware(jwtService, getTestLogger())

	router.GET("/optional", authMiddleware.OptionalAuth(), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.False(t, exists) // Should not be set

		c.JSON(http.StatusOK, gin.H{"authenticated": false})
	})

	req := httptest.NewRequest("GET", "/optional", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// OptionalAuth should allow the request to proceed even with invalid token
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_ContextPropagation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	jwtService := getTestJWTService()
	authMiddleware := NewAuthMiddleware(jwtService, getTestLogger())

	permissions := []string{"users:read", "users:write", "tasks:delete"}
	token, err := jwtService.GenerateAccessToken(42, "john.doe", "manager", permissions)
	require.NoError(t, err)

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		role, _ := c.Get("role")
		perms, _ := c.Get("permissions")

		assert.Equal(t, uint(42), userID)
		assert.Equal(t, "john.doe", username)
		assert.Equal(t, "manager", role)
		assert.Equal(t, permissions, perms)

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_RequireAuth_WrongSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	correctService := getTestJWTService()
	wrongService := util.NewJWTService("wrong-secret", 15*time.Minute, 7*24*time.Hour)

	// Generate token with correct secret
	token, err := correctService.GenerateAccessToken(1, "testuser", "admin", []string{})
	require.NoError(t, err)

	// Create middleware with wrong secret
	authMiddleware := NewAuthMiddleware(wrongService, getTestLogger())

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail validation due to wrong secret
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
