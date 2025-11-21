# Go バックエンド開発ガイド

## 目次
1. [開発環境セットアップ](#開発環境セットアップ)
2. [プロジェクト構成](#プロジェクト構成)
3. [開発の流れ](#開発の流れ)
4. [Gin フレームワーク](#ginフレームワーク)
5. [GORM でのデータアクセス](#gormでのデータアクセス)
6. [認証・認可実装](#認証認可実装)
7. [テスト](#テスト)
8. [コード品質](#コード品質)

---

## 開発環境セットアップ

### 前提条件
- Go 1.21 以上
- PostgreSQL 13.0 以上
- Redis 6.0 以上
- Docker & Docker Compose
- Git

### セットアップ手順

```bash
# リポジトリクローン
git clone <repository-url> internalsystem
cd internalsystem/backend

# Go モジュール初期化（新規プロジェクト）
go mod init github.com/yourusername/internalsystem

# 依存関係インストール
go mod download

# Docker Compose で DB・Redis 起動
docker-compose up -d

# マイグレーション実行
migrate -path ./migrations -database "postgres://user:password@localhost:5432/internalsystem" up

# サーバー起動
go run cmd/server/main.go
```

### 環境変数設定

`.env.local` ファイルを作成：

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=internalsystem

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=3600

# Server
SERVER_PORT=8080
ENV=development
LOG_LEVEL=debug
```

---

## プロジェクト構成

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # エントリーポイント
├── internal/
│   ├── config/
│   │   └── config.go              # 設定管理
│   ├── models/                    # GORM モデル
│   │   ├── user.go
│   │   ├── role.go
│   │   └── ...
│   ├── repository/                # データアクセス層
│   │   ├── user_repository.go
│   │   ├── role_repository.go
│   │   └── ...
│   ├── service/                   # ビジネスロジック層
│   │   ├── user_service.go
│   │   ├── auth_service.go
│   │   └── ...
│   ├── handler/                   # HTTP ハンドラ
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   └── ...
│   ├── middleware/                # ミドルウェア
│   │   ├── auth.go               # JWT 検証
│   │   ├── cors.go
│   │   └── logger.go
│   └── utils/
│       ├── errors.go             # エラーハンドリング
│       ├── response.go           # レスポンス整形
│       └── validators.go         # バリデーション
├── migrations/                    # DB マイグレーション
│   ├── 001_initial_schema.up.sql
│   └── 001_initial_schema.down.sql
├── tests/
│   ├── unit/
│   ├── integration/
│   └── fixtures/
├── go.mod
├── go.sum
├── Dockerfile
└── main.go (→ cmd/server/main.go)
```

---

## 開発の流れ

### 1. モデル定義（models/）

```go
// internal/models/user.go
package models

import "time"

type User struct {
    ID        int64     `gorm:"primaryKey"`
    Username  string    `gorm:"uniqueIndex;not null"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    FullName  string
    Department string
    RoleID    int64
    Role      *Role     `gorm:"foreignKey:RoleID"`
    Status    string    `gorm:"default:'active'"`
    LastLogin *time.Time
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
    DeletedAt *time.Time `gorm:"index"`
}

func (User) TableName() string {
    return "users"
}
```

### 2. リポジトリ実装（repository/）

```go
// internal/repository/user_repository.go
package repository

import "gorm.io/gorm"

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// GetByEmail ユーザーをメールアドレスで取得
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    var user models.User
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// Create ユーザーを作成
func (r *UserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

// List ユーザー一覧を取得
func (r *UserRepository) List(offset, limit int) ([]models.User, error) {
    var users []models.User
    if err := r.db.
        Preload("Role").
        Offset(offset).
        Limit(limit).
        Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
```

### 3. サービス実装（service/）

```go
// internal/service/auth_service.go
package service

type AuthService struct {
    userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
    return &AuthService{userRepo: userRepo}
}

// Login ユーザーをログイン
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
    // ユーザーを取得
    user, err := s.userRepo.GetByEmail(email)
    if err != nil {
        return nil, "", errors.New("user not found")
    }

    // パスワード検証
    if !s.verifyPassword(user.Password, password) {
        return nil, "", errors.New("invalid password")
    }

    // JWT トークン生成
    token, err := s.generateToken(user)
    if err != nil {
        return nil, "", err
    }

    return user, token, nil
}

// トークン生成
func (s *AuthService) generateToken(user *models.User) (string, error) {
    // JWT ライブラリを使用してトークン生成
    // ...
}
```

### 4. ハンドラ実装（handler/）

```go
// internal/handler/auth_handler.go
package handler

import "github.com/gin-gonic/gin"

type AuthHandler struct {
    authService *service.AuthService
}

// POST /api/v1/auth/login
// @Summary ユーザーログイン
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "ログイン情報"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }

    user, token, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(401, gin.H{"error": "authentication failed"})
        return
    }

    c.JSON(200, gin.H{
        "user": user,
        "token": token,
    })
}
```

### 5. ルーティング設定（main.go）

```go
// cmd/server/main.go
package main

import "github.com/gin-gonic/gin"

func main() {
    // 初期化
    config := config.LoadConfig()
    db := database.Connect(config)

    // リポジトリ・サービス初期化
    userRepo := repository.NewUserRepository(db)
    authService := service.NewAuthService(userRepo)
    authHandler := handler.NewAuthHandler(authService)

    // Gin ルーター
    router := gin.Default()

    // ミドルウェア
    router.Use(middleware.CORSMiddleware())
    router.Use(middleware.LoggerMiddleware())

    // ルート定義
    api := router.Group("/api/v1")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/login", authHandler.Login)
            auth.POST("/logout", authHandler.Logout)
        }
    }

    // サーバー起動
    router.Run(":8080")
}
```

---

## Gin フレームワーク

### 基本的な使い方

```go
// ルーター初期化
router := gin.Default()

// ルート登録
router.GET("/users/:id", GetUser)
router.POST("/users", CreateUser)
router.PUT("/users/:id", UpdateUser)
router.DELETE("/users/:id", DeleteUser)

// グループ化
api := router.Group("/api/v1")
api.GET("/users", ListUsers)

// ミドルウェア
router.Use(middleware.AuthMiddleware())

// 実行
router.Run(":8080")
```

### リクエスト処理

```go
func CreateUser(c *gin.Context) {
    // JSON リクエスト解析
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // クエリパラメータ
    page := c.Query("page")

    // パスパラメータ
    id := c.Param("id")

    // レスポンス
    c.JSON(201, gin.H{"user": user})
}
```

---

## GORM でのデータアクセス

### モデル関連付け

```go
// 1:N 関連
type Role struct {
    ID    int64
    Name  string
    Users []User
}

type User struct {
    ID     int64
    RoleID int64
    Role   *Role
}

// N:N 関連
type User struct {
    ID          int64
    Permissions []Permission `gorm:"many2many:user_permissions;"`
}

// Preload（関連データ一括取得）
var users []User
db.Preload("Role").Preload("Permissions").Find(&users)
```

### クエリ

```go
// 単一取得
var user User
db.First(&user, id)

// 一覧取得
var users []User
db.Where("status = ?", "active").Find(&users)

// ページネーション
db.Offset(offset).Limit(limit).Find(&users)

// トランザクション
tx := db.BeginTx(ctx, nil)
// ...
tx.Commit()
```

---

## 認証・認可実装

### JWT ミドルウェア

```go
// internal/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")

        claims, err := jwt.ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

### RBAC チェック

```go
// internal/middleware/rbac.go
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("role")

        allowed := false
        for _, r := range roles {
            if userRole == r {
                allowed = true
                break
            }
        }

        if !allowed {
            c.JSON(403, gin.H{"error": "forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// 使用例
auth.DELETE("/users/:id",
    middleware.AuthMiddleware(),
    middleware.RequireRole("admin"),
    DeleteUser,
)
```

---

## テスト

### ユニットテスト

```go
// internal/service/auth_service_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := NewAuthService(mockRepo)

    // Act
    user, token, err := service.Login("test@example.com", "password")

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.NotEmpty(t, token)
}
```

### 統合テスト

```go
// tests/integration/auth_test.go
func TestLoginEndpoint(t *testing.T) {
    router := setupRouter()

    req := httptest.NewRequest(
        "POST",
        "/api/v1/auth/login",
        bytes.NewBufferString(`{"email":"test@example.com","password":"password"}`),
    )

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

---

## コード品質

### 実行

```bash
# テスト
go test ./... -v

# カバレッジ
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# リント
golangci-lint run ./...

# フォーマット
go fmt ./...

# Vet（静的解析）
go vet ./...
```

### ベストプラクティス

1. **エラーハンドリング**: すべてのエラーを適切に処理
2. **ロギング**: 重要な操作はログに記録
3. **バリデーション**: ユーザー入力を厳密にバリデーション
4. **セキュリティ**: パスワードハッシング、SQL インジェクション対策
5. **テスト**: 重要ロジックは必ずテスト
6. **ドキュメント**: 公開関数には godoc コメント

---

## よくある問題と解決方法

### データベース接続エラー
```bash
# 原因：DB サーバーが起動していない
docker-compose up -d

# または接続情報が間違っている
# .env.local を確認
```

### JWT トークン無効エラー
```bash
# 原因：トークンの署名キーが異なる
# JWT_SECRET が同じ値に設定されているか確認
```

### CORS エラー
```bash
# 原因：フロントエンドからのリクエストが拒否されている
# CORS ミドルウェアが正しく設定されているか確認
```

---

このガイドを参考に、チーム内で開発を進めてください。不明点があれば、docs/ に追加ドキュメントを作成してください。
