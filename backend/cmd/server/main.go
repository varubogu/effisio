package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/config"
	"github.com/varubogu/effisio/backend/internal/handler"
	"github.com/varubogu/effisio/backend/internal/middleware"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
)

func main() {
	// .envファイルを読み込み（開発環境用）
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .envファイルが見つかりません。環境変数を使用します")
	}

	// ロガーの初期化
	logger, err := initLogger()
	if err != nil {
		log.Fatalf("❌ ロガーの初期化に失敗しました: %v", err)
	}
	defer logger.Sync()

	// 設定の読み込み
	cfg := config.Load()

	// データベース接続
	db, err := initDB(cfg)
	if err != nil {
		logger.Fatal("❌ データベース接続に失敗しました", zap.Error(err))
	}
	logger.Info("✅ データベースに接続しました")

	// Redis接続（将来的に実装）
	// redisClient := initRedis(cfg)

	// ユーティリティの初期化
	jwtService := util.NewJWTService(
		[]byte(cfg.JWT.Secret),
		cfg.JWT.AccessTokenExpiration,
		cfg.JWT.RefreshTokenExpiration,
	)

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	// サービスの初期化
	userService := service.NewUserService(userRepo, logger)
	permissionService := service.NewPermissionService(permissionRepo, logger)
	roleService := service.NewRoleService(roleRepo, permissionRepo, logger)
	auditLogService := service.NewAuditLogService(auditLogRepo, userRepo, logger)
	authService := service.NewAuthService(userRepo, refreshTokenRepo, jwtService, roleService, logger)

	// ハンドラーの初期化
	healthHandler := handler.NewHealthHandler(logger)
	userHandler := handler.NewUserHandler(userService, auditLogService, logger)
	authHandler := handler.NewAuthHandler(authService, auditLogService, logger)
	permissionHandler := handler.NewPermissionHandler(permissionService, logger)
	roleHandler := handler.NewRoleHandler(roleService, logger)
	auditLogHandler := handler.NewAuditLogHandler(auditLogService, logger)

	// ミドルウェアの初期化
	authMiddleware := middleware.NewAuthMiddleware(jwtService, logger)
	rbacMiddleware := middleware.NewRBACMiddleware(logger)

	// Ginルーターの設定
	router := setupRouter(cfg, logger, healthHandler, userHandler, authHandler, permissionHandler, roleHandler, auditLogHandler, authMiddleware, rbacMiddleware)

	// HTTPサーバーの設定
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// グレースフルシャットダウンの設定
	go func() {
		logger.Info("🚀 サーバーを起動しています",
			zap.String("port", cfg.Server.Port),
			zap.String("env", cfg.Server.Env),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("❌ サーバーの起動に失敗しました", zap.Error(err))
		}
	}()

	// シグナルを待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 サーバーをシャットダウンしています...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("❌ サーバーのシャットダウンに失敗しました", zap.Error(err))
	}

	logger.Info("✅ サーバーが正常に終了しました")
}

// initLogger はロガーを初期化します
func initLogger() (*zap.Logger, error) {
	env := os.Getenv("ENV")
	if env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

// initDB はデータベース接続を初期化します
func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// コネクションプールの設定
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	return db, nil
}

// setupRouter はGinルーターを設定します
func setupRouter(
	cfg *config.Config,
	logger *zap.Logger,
	healthHandler *handler.HealthHandler,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	permissionHandler *handler.PermissionHandler,
	roleHandler *handler.RoleHandler,
	auditLogHandler *handler.AuditLogHandler,
	authMiddleware *middleware.AuthMiddleware,
	rbacMiddleware *middleware.RBACMiddleware,
) *gin.Engine {
	// 本番環境ではリリースモードに設定
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// ミドルウェアの設定
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS(cfg))

	// ヘルスチェックエンドポイント
	router.GET("/health", healthHandler.Check)

	// APIルート
	api := router.Group("/api/v1")
	{
		// 認証不要のエンドポイント
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		// 認証関連
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)

			// 認証が必要なエンドポイント
			auth.POST("/logout-all", authMiddleware.RequireAuth(), authHandler.LogoutAll)
		}

		// ユーザー関連（認証と権限が必要）
		users := api.Group("/users")
		users.Use(authMiddleware.RequireAuth()) // 全てのユーザーエンドポイントで認証が必要
		{
			// 一覧取得と詳細取得は全ての認証済みユーザーが可能
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.GetByID)

			// 作成は admin のみ
			users.POST("", rbacMiddleware.RequireRole("admin"), userHandler.Create)

			// 更新は admin と manager のみ
			users.PUT("/:id", rbacMiddleware.RequireAnyRole("admin", "manager"), userHandler.Update)

			// 削除は admin のみ
			users.DELETE("/:id", rbacMiddleware.RequireRole("admin"), userHandler.Delete)
		}

		// 権限管理（admin のみ）
		permissions := api.Group("/permissions")
		permissions.Use(authMiddleware.RequireAuth())
		permissions.Use(rbacMiddleware.RequirePermission("permissions:read"))
		{
			permissions.GET("", permissionHandler.List)
			permissions.GET("/:id", permissionHandler.GetByID)

			// 作成・更新・削除は permissions:write が必要
			permissions.POST("", rbacMiddleware.RequirePermission("permissions:write"), permissionHandler.Create)
			permissions.PUT("/:id", rbacMiddleware.RequirePermission("permissions:write"), permissionHandler.Update)
			permissions.DELETE("/:id", rbacMiddleware.RequirePermission("permissions:write"), permissionHandler.Delete)
		}

		// ロール管理（admin のみ）
		roles := api.Group("/roles")
		roles.Use(authMiddleware.RequireAuth())
		roles.Use(rbacMiddleware.RequirePermission("roles:read"))
		{
			roles.GET("", roleHandler.List)
			roles.GET("/:id", roleHandler.GetByID)

			// 作成・更新・削除は roles:write が必要
			roles.POST("", rbacMiddleware.RequirePermission("roles:write"), roleHandler.Create)
			roles.PUT("/:id", rbacMiddleware.RequirePermission("roles:write"), roleHandler.Update)
			roles.DELETE("/:id", rbacMiddleware.RequirePermission("roles:write"), roleHandler.Delete)
		}

		// 監査ログ（audit:read 権限が必要）
		auditLogs := api.Group("/audit-logs")
		auditLogs.Use(authMiddleware.RequireAuth())
		auditLogs.Use(rbacMiddleware.RequirePermission("audit:read"))
		{
			auditLogs.GET("", auditLogHandler.List)
			auditLogs.GET("/:id", auditLogHandler.GetByID)
		}
	}

	return router
}
