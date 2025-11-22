package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthHandler はヘルスチェックハンドラーです
type HealthHandler struct {
	logger *zap.Logger
}

// NewHealthHandler は新しいHealthHandlerを作成します
func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// Check はヘルスチェックエンドポイントです
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "effisio-api",
	})
}
