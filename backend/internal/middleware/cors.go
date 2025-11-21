package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/varubogu/effisio/backend/internal/config"
)

// CORS はCORS設定を行うミドルウェアです
func CORS(cfg *config.Config) gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12時間
	}

	// 本番環境では許可するオリジンを制限
	if cfg.Server.Env == "production" {
		// 実際の本番環境のドメインに変更してください
		config.AllowOrigins = []string{"https://your-production-domain.com"}
	}

	return cors.New(config)
}
