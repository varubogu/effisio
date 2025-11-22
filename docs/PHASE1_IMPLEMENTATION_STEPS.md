# Phase 1実装手順書

このドキュメントでは、Phase 1の実装を順を追って説明します。各ステップを完了させることで、基本的なCRUD機能が完全に動作するようになります。

## 目次

- [Phase 1の概要](#phase-1の概要)
- [前提条件](#前提条件)
- [実装ステップ](#実装ステップ)
  - [Step 1: スキーマ移行](#step-1-スキーマ移行)
  - [Step 2: モデル層の更新](#step-2-モデル層の更新)
  - [Step 3: ユーティリティ関数の実装](#step-3-ユーティリティ関数の実装)
  - [Step 4: リポジトリ層の完全実装](#step-4-リポジトリ層の完全実装)
  - [Step 5: サービス層の完全実装](#step-5-サービス層の完全実装)
  - [Step 6: ハンドラー層の完全実装](#step-6-ハンドラー層の完全実装)
  - [Step 7: バックエンドテストの実装](#step-7-バックエンドテストの実装)
  - [Step 8: フロントエンドの更新](#step-8-フロントエンドの更新)
  - [Step 9: 統合テスト](#step-9-統合テスト)
- [完了チェックリスト](#完了チェックリスト)

---

## Phase 1の概要

**目標:**
- ユーザーCRUD機能の完全実装
- 設計仕様に完全準拠したスキーマ
- ページネーション実装
- 統一されたエラーハンドリング
- バリデーション実装
- 単体テストのカバレッジ70%以上

**所要時間:** 1-2週間

**成果物:**
- 動作する基本CRUD API
- テストコード
- 更新されたフロントエンド
- ドキュメント

---

## 前提条件

以下が完了していることを確認してください:

- [ ] **[QUICK_START.md](QUICK_START.md)** に従って開発環境をセットアップ済み
- [ ] `make dev` で開発環境が起動する
- [ ] Phase 0のコードが動作している
- [ ] Git作業ブランチを作成済み（例: `feature/phase1-implementation`）

```bash
# 作業ブランチを作成
git checkout main
git pull origin main
git checkout -b feature/phase1-implementation
```

---

## 実装ステップ

### Step 1: スキーマ移行

**目的:** Phase 0の簡略化スキーマを設計仕様に準拠したスキーマに移行する

**詳細:** [MIGRATION_FROM_PHASE0.md](MIGRATION_FROM_PHASE0.md) を参照

#### 1-1. マイグレーションファイルの作成

```bash
# マイグレーションファイルを作成
touch backend/migrations/000002_update_users_table.up.sql
touch backend/migrations/000002_update_users_table.down.sql
```

**ファイル: `backend/migrations/000002_update_users_table.up.sql`**

```sql
-- 新しいカラムを追加
ALTER TABLE users ADD COLUMN full_name VARCHAR(100);
ALTER TABLE users ADD COLUMN department VARCHAR(100);
ALTER TABLE users ADD COLUMN last_login TIMESTAMP;

-- パスワードカラムをリネーム
ALTER TABLE users RENAME COLUMN password TO password_hash;

-- ステータスカラムを追加（一時的にNULL許可）
ALTER TABLE users ADD COLUMN status VARCHAR(20);

-- 既存データのステータスを設定
UPDATE users SET status = CASE
    WHEN is_active = true THEN 'active'
    WHEN is_active = false THEN 'inactive'
    ELSE 'active'
END;

-- statusをNOT NULLに変更してデフォルト値を設定
ALTER TABLE users ALTER COLUMN status SET NOT NULL;
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active';

-- is_activeカラムを削除
ALTER TABLE users DROP COLUMN is_active;

-- CHECK制約を追加
ALTER TABLE users ADD CONSTRAINT users_status_check
    CHECK (status IN ('active', 'inactive', 'suspended'));

-- インデックスを追加
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_department ON users(department);
```

**ファイル: `backend/migrations/000002_update_users_table.down.sql`**

```sql
-- ロールバック用
DROP INDEX IF EXISTS idx_users_department;
DROP INDEX IF EXISTS idx_users_status;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT true;
UPDATE users SET is_active = CASE WHEN status = 'active' THEN true ELSE false END;
ALTER TABLE users DROP COLUMN status;
ALTER TABLE users RENAME COLUMN password_hash TO password;
ALTER TABLE users DROP COLUMN last_login;
ALTER TABLE users DROP COLUMN department;
ALTER TABLE users DROP COLUMN full_name;
```

#### 1-2. マイグレーション実行

```bash
# データベースバックアップ（念のため）
docker-compose exec postgres pg_dump -U postgres effisio_dev > backup_before_migration.sql

# マイグレーション実行
make migrate-up

# 結果確認
docker-compose exec postgres psql -U postgres -d effisio_dev -c "\d users"

# 期待される出力: full_name, department, password_hash, status, last_login が存在する
```

#### 1-3. シードデータの更新

**ファイル: `backend/scripts/seed.sh`**

```bash
#!/bin/bash

set -e

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-effisio_dev}

echo "🌱 シードデータを投入しています..."

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF

DELETE FROM users WHERE username IN ('admin', 'manager', 'testuser', 'viewer', 'suspended_user');

INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at)
VALUES
  ('admin', 'admin@example.com', '管理者 太郎', 'IT部',
   '\$2a\$10\$X8yI6qZvKZH5mP3nR4tVH.YqJ5mN6oP7qR8sT9uV0wX1yZ2aB3cD4',
   'admin', 'active', NOW(), NOW()),
  ('manager', 'manager@example.com', '管理 次郎', '営業部',
   '\$2a\$10\$Y9zJ7rAvLAI6nQ4oS5uWI.ZrK6nO7pQ8rS9tU0vW1xY2zA3bC4dE5',
   'manager', 'active', NOW(), NOW()),
  ('testuser', 'testuser@example.com', 'テスト 三郎', '開発部',
   '\$2a\$10\$Z0aK8sBwMBJ7oR5pT6vXJ.AsL7oP8qR9sT0uV1wX2yZ3aB4cD5eF6',
   'user', 'active', NOW(), NOW()),
  ('viewer', 'viewer@example.com', '閲覧 四郎', '総務部',
   '\$2a\$10\$A1bL9tCxNCK8pS6qU7wYK.BtM8pQ9rS0tU1vW2xY3zA4bC5dD6eG7',
   'viewer', 'active', NOW(), NOW()),
  ('suspended_user', 'suspended@example.com', '停止 五郎', 'なし',
   '\$2a\$10\$B2cM0uDyODL9qT7rV8xZL.CuN9qR0sT1uV2wX3yZ4aB5cD6eE7fH8',
   'user', 'suspended', NOW(), NOW());

EOF

echo "✅ シードデータの投入が完了しました"
```

```bash
# シードデータを投入
make seed

# 結果確認
docker-compose exec postgres psql -U postgres -d effisio_dev -c "SELECT id, username, full_name, department, status FROM users;"
```

**✅ Checkpoint:** スキーマが更新され、新しいフィールドが追加されている

---

### Step 2: モデル層の更新

**目的:** Goのモデル定義を新しいスキーマに対応させる

#### 2-1. Userモデルの更新

**ファイル: `backend/internal/model/user.go`**

```go
package model

import (
	"time"
	"gorm.io/gorm"
)

// User ユーザーモデル
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	FullName     string         `gorm:"size:100" json:"full_name"`
	Department   string         `gorm:"size:100" json:"department"`
	PasswordHash string         `gorm:"not null;size:255;column:password_hash" json:"-"`
	Role         string         `gorm:"not null;size:20;default:'user'" json:"role"`
	Status       string         `gorm:"not null;size:20;default:'active'" json:"status"`
	LastLogin    *time.Time     `json:"last_login"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName テーブル名を指定
func (User) TableName() string {
	return "users"
}

// ステータス定数
const (
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
)

// ロール定数
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleUser    = "user"
	RoleViewer  = "viewer"
)

// UserResponse APIレスポンス用
type UserResponse struct {
	ID         uint       `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	FullName   string     `json:"full_name"`
	Department string     `json:"department"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	LastLogin  *time.Time `json:"last_login"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ToResponse UserからUserResponseへ変換
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		FullName:   u.FullName,
		Department: u.Department,
		Role:       u.Role,
		Status:     u.Status,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

// CreateUserRequest ユーザー作成リクエスト
type CreateUserRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Email      string `json:"email" binding:"required,email"`
	FullName   string `json:"full_name" binding:"max=100"`
	Department string `json:"department" binding:"max=100"`
	Password   string `json:"password" binding:"required,min=8,max=72"`
	Role       string `json:"role" binding:"required,oneof=admin manager user viewer"`
}

// UpdateUserRequest ユーザー更新リクエスト
type UpdateUserRequest struct {
	Email      *string `json:"email" binding:"omitempty,email"`
	FullName   *string `json:"full_name" binding:"omitempty,max=100"`
	Department *string `json:"department" binding:"omitempty,max=100"`
	Role       *string `json:"role" binding:"omitempty,oneof=admin manager user viewer"`
	Status     *string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// IsValidStatus ステータスの妥当性チェック
func IsValidStatus(status string) bool {
	return status == UserStatusActive || status == UserStatusInactive || status == UserStatusSuspended
}

// IsValidRole ロールの妥当性チェック
func IsValidRole(role string) bool {
	return role == RoleAdmin || role == RoleManager || role == RoleUser || role == RoleViewer
}
```

#### 2-2. ビルド確認

```bash
cd backend
go mod tidy
go build -o bin/server cmd/server/main.go

# エラーがないことを確認
```

**✅ Checkpoint:** モデルが更新され、ビルドエラーがない

---

### Step 3: ユーティリティ関数の実装

**目的:** 統一されたレスポンス形式、エラーハンドリング、ページネーションを実装

#### 3-1. エラー定義

**ファイル: `backend/pkg/util/error.go`**

```go
package util

import (
	"fmt"
	"net/http"
)

// エラーコード定数
const (
	// 認証エラー (AUTH_xxx)
	ErrCodeUnauthorized      = "AUTH_001"
	ErrCodeInvalidToken      = "AUTH_002"
	ErrCodeTokenExpired      = "AUTH_003"
	ErrCodeInsufficientPermission = "AUTH_004"

	// ユーザーエラー (USER_xxx)
	ErrCodeUserNotFound      = "USER_001"
	ErrCodeUserAlreadyExists = "USER_002"
	ErrCodeInvalidCredentials = "USER_003"

	// バリデーションエラー (VAL_xxx)
	ErrCodeValidationError   = "VAL_001"
	ErrCodeInvalidParameter  = "VAL_002"

	// データベースエラー (DB_xxx)
	ErrCodeDatabaseError     = "DB_001"
	ErrCodeRecordNotFound    = "DB_002"

	// システムエラー (SYS_xxx)
	ErrCodeInternalError     = "SYS_001"
	ErrCodePasswordHashError = "SYS_002"
)

// AppError アプリケーションエラー
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// エラーコンストラクタ
func NewBadRequestError(code string, err error) *AppError {
	return &AppError{Code: code, Message: "Bad request", StatusCode: http.StatusBadRequest, Err: err}
}

func NewUnauthorizedError(code string, err error) *AppError {
	return &AppError{Code: code, Message: "Unauthorized", StatusCode: http.StatusUnauthorized, Err: err}
}

func NewForbiddenError(code string, err error) *AppError {
	return &AppError{Code: code, Message: "Forbidden", StatusCode: http.StatusForbidden, Err: err}
}

func NewNotFoundError(code string, err error) *AppError {
	return &AppError{Code: code, Message: "Resource not found", StatusCode: http.StatusNotFound, Err: err}
}

func NewConflictError(code string, err error) *AppError {
	return &AppError{Code: code, Message: "Resource conflict", StatusCode: http.StatusConflict, Err: err}
}

func NewInternalError(code string, err error) *AppError {
	return &AppError{Code: code, Message: "Internal server error", StatusCode: http.StatusInternalServerError, Err: err}
}
```

#### 3-2. レスポンスヘルパー

**ファイル: `backend/pkg/util/response.go`**

```go
package util

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 統一レスポンス形式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success 成功レスポンス (200 OK)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// Created 作成成功レスポンス (201 Created)
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "created",
		Data:    data,
	})
}

// NoContent 内容なし (204 No Content)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error エラーレスポンス
func Error(c *gin.Context, statusCode int, code string, message string, details interface{}) {
	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: "error",
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ValidationError バリデーションエラー
func ValidationError(c *gin.Context, errors map[string]string) {
	Error(c, http.StatusBadRequest, ErrCodeValidationError, "Validation failed", errors)
}

// HandleError AppErrorからレスポンスを生成
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		Error(c, appErr.StatusCode, appErr.Code, appErr.Message, nil)
	} else {
		Error(c, http.StatusInternalServerError, ErrCodeInternalError, "Internal server error", nil)
	}
}

// ParseValidationErrors バリデーションエラーをパース
func ParseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			errors[field] = fmt.Sprintf("validation failed on tag '%s'", tag)
		}
	}
	return errors
}
```

#### 3-3. ページネーション

**ファイル: `backend/pkg/util/pagination.go`**

```go
package util

import (
	"math"
	"strconv"
	"github.com/gin-gonic/gin"
)

// PaginationParams ページネーションパラメータ
type PaginationParams struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Offset  int `json:"-"`
}

// PaginationInfo ページネーション情報
type PaginationInfo struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse ページネーション付きレスポンス
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// GetPaginationParams リクエストからページネーションパラメータを取得
func GetPaginationParams(c *gin.Context) *PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	return &PaginationParams{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
	}
}

// NewPaginationInfo ページネーション情報を生成
func NewPaginationInfo(total int64, params *PaginationParams) *PaginationInfo {
	totalPages := int(math.Ceil(float64(total) / float64(params.PerPage)))
	return &PaginationInfo{
		Page:       params.Page,
		PerPage:    params.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

// NewPaginatedResponse ページネーション付きレスポンスを生成
func NewPaginatedResponse(data interface{}, total int64, params *PaginationParams) *PaginatedResponse {
	return &PaginatedResponse{
		Data:       data,
		Pagination: *NewPaginationInfo(total, params),
	}
}

// Paginated ページネーション付きレスポンスを返す
func Paginated(c *gin.Context, response *PaginatedResponse) {
	c.JSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"message":    "success",
		"data":       response.Data,
		"pagination": response.Pagination,
	})
}
```

**✅ Checkpoint:** ユーティリティ関数が実装され、ビルドエラーがない

---

### Step 4: リポジトリ層の完全実装

**目的:** データベース操作を完全に実装

**ファイル: `backend/internal/repository/user.go`**

```go
package repository

import (
	"context"
	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/pkg/util"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAll 全ユーザーを取得（ページネーション付き）
func (r *UserRepository) FindAll(ctx context.Context, params *util.PaginationParams) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 総件数を取得
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション付きで取得
	err := r.db.WithContext(ctx).
		Offset(params.Offset).
		Limit(params.PerPage).
		Order("id ASC").
		Find(&users).Error

	return users, total, err
}

// FindByID IDでユーザーを取得
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail メールアドレスでユーザーを取得
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername ユーザー名でユーザーを取得
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create ユーザーを作成
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update ユーザーを更新
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete ユーザーを削除（ソフトデリート）
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// ExistsByEmail メールアドレスの存在確認
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByUsername ユーザー名の存在確認
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
```

**✅ Checkpoint:** リポジトリ層が実装され、ビルドエラーがない

---

### Step 5: サービス層の完全実装

**目的:** ビジネスロジックを実装

**ファイル: `backend/internal/service/user.go`**

```go
package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
)

type UserService struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo *repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// List ユーザー一覧を取得
func (s *UserService) List(ctx context.Context, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	users, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		s.logger.Error("Failed to fetch users", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return util.NewPaginatedResponse(responses, total, params), nil
}

// GetByID IDでユーザーを取得
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError(util.ErrCodeUserNotFound, err)
		}
		s.logger.Error("Failed to fetch user", zap.Uint("id", id), zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return user.ToResponse(), nil
}

// Create ユーザーを作成
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	// ユーザー名の重複チェック
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}
	if exists {
		return nil, util.NewConflictError(util.ErrCodeUserAlreadyExists, errors.New("username already exists"))
	}

	// メールアドレスの重複チェック
	exists, err = s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}
	if exists {
		return nil, util.NewConflictError(util.ErrCodeUserAlreadyExists, errors.New("email already exists"))
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodePasswordHashError, err)
	}

	// ユーザーモデルを作成
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		Department:   req.Department,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		Status:       model.UserStatusActive,
	}

	// データベースに保存
	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("User created", zap.Uint("id", user.ID), zap.String("username", user.Username))
	return user.ToResponse(), nil
}

// Update ユーザーを更新
func (s *UserService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	// 既存ユーザーを取得
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewNotFoundError(util.ErrCodeUserNotFound, err)
		}
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 更新データを適用
	if req.Email != nil {
		// メールアドレスの重複チェック（自分以外）
		existingUser, err := s.repo.FindByEmail(ctx, *req.Email)
		if err == nil && existingUser.ID != id {
			return nil, util.NewConflictError(util.ErrCodeUserAlreadyExists, errors.New("email already exists"))
		}
		user.Email = *req.Email
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Department != nil {
		user.Department = *req.Department
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	// データベースを更新
	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("Failed to update user", zap.Uint("id", id), zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("User updated", zap.Uint("id", user.ID))
	return user.ToResponse(), nil
}

// Delete ユーザーを削除
func (s *UserService) Delete(ctx context.Context, id uint) error {
	// 存在確認
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.NewNotFoundError(util.ErrCodeUserNotFound, err)
		}
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 削除実行
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete user", zap.Uint("id", id), zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("User deleted", zap.Uint("id", id))
	return nil
}
```

**✅ Checkpoint:** サービス層が実装され、ビルドエラーがない

---

### Step 6: ハンドラー層の完全実装

**目的:** HTTP APIエンドポイントを実装

**ファイル: `backend/internal/handler/user.go`**

```go
package handler

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/service"
	"github.com/varubogu/effisio/backend/pkg/util"
)

type UserHandler struct {
	service *service.UserService
	logger  *zap.Logger
}

func NewUserHandler(service *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// List ユーザー一覧を取得
// @Summary ユーザー一覧取得
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param per_page query int false "1ページあたりの件数" default(10)
// @Success 200 {object} util.PaginatedResponse
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	params := util.GetPaginationParams(c)
	result, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Paginated(c, result)
}

// GetByID IDでユーザーを取得
// @Summary ユーザー詳細取得
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Success 200 {object} model.UserResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.Error(c, 400, util.ErrCodeInvalidParameter, "Invalid user ID", nil)
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, gin.H{"user": user})
}

// Create ユーザーを作成
// @Summary ユーザー作成
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "ユーザー作成リクエスト"
// @Success 201 {object} model.UserResponse
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, util.ParseValidationErrors(err))
		return
	}

	user, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Created(c, gin.H{"user": user})
}

// Update ユーザーを更新
// @Summary ユーザー更新
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Param request body model.UpdateUserRequest true "ユーザー更新リクエスト"
// @Success 200 {object} model.UserResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.Error(c, 400, util.ErrCodeInvalidParameter, "Invalid user ID", nil)
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, util.ParseValidationErrors(err))
		return
	}

	user, err := h.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, gin.H{"user": user})
}

// Delete ユーザーを削除
// @Summary ユーザー削除
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ユーザーID"
// @Success 204
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.Error(c, 400, util.ErrCodeInvalidParameter, "Invalid user ID", nil)
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		util.HandleError(c, err)
		return
	}

	util.NoContent(c)
}
```

#### ルーティングの更新

**ファイル: `backend/cmd/server/main.go`** の一部を更新

```go
func setupRouter(cfg *config.Config, logger *zap.Logger, healthHandler *handler.HealthHandler, userHandler *handler.UserHandler) *gin.Engine {
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", healthHandler.Check)

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		users := v1.Group("/users")
		{
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.GetByID)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}

	return r
}
```

#### ビルドとテスト

```bash
cd backend

# ビルド
go mod tidy
make build

# 起動
docker-compose restart backend

# ログ確認
docker-compose logs -f backend

# APIテスト
curl http://localhost:8080/api/v1/ping
curl http://localhost:8080/api/v1/users | jq
```

**✅ Checkpoint:** ハンドラー層が実装され、APIが動作する

---

### Step 7: バックエンドテストの実装

**目的:** 単体テストを実装してカバレッジ70%以上を達成

#### 7-1. リポジトリテスト

**ファイル: `backend/internal/repository/user_test.go`**

```go
package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/pkg/util"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.User{})
	require.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// テストデータ作成
	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	repo.Create(context.Background(), user)

	// 取得テスト
	found, err := repo.FindByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, found.Username)
}

func TestUserRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// テストデータ作成
	users := []*model.User{
		{Username: "user1", Email: "user1@example.com", PasswordHash: "hash", Role: model.RoleUser, Status: model.UserStatusActive},
		{Username: "user2", Email: "user2@example.com", PasswordHash: "hash", Role: model.RoleUser, Status: model.UserStatusActive},
		{Username: "user3", Email: "user3@example.com", PasswordHash: "hash", Role: model.RoleUser, Status: model.UserStatusActive},
	}
	for _, u := range users {
		repo.Create(context.Background(), u)
	}

	// ページネーション付き取得
	params := &util.PaginationParams{Page: 1, PerPage: 10, Offset: 0}
	result, total, err := repo.FindAll(context.Background(), params)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, result, 3)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// テストデータ作成
	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	repo.Create(context.Background(), user)

	// 更新
	user.Email = "updated@example.com"
	err := repo.Update(context.Background(), user)
	assert.NoError(t, err)

	// 確認
	updated, _ := repo.FindByID(context.Background(), user.ID)
	assert.Equal(t, "updated@example.com", updated.Email)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// テストデータ作成
	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	repo.Create(context.Background(), user)

	// 削除
	err := repo.Delete(context.Background(), user.ID)
	assert.NoError(t, err)

	// 確認（ソフトデリートなので Unscoped が必要）
	var deleted model.User
	result := db.Unscoped().First(&deleted, user.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, deleted.DeletedAt)
}
```

#### 7-2. テスト実行

```bash
cd backend

# テスト実行
go test ./... -v

# カバレッジ付き
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# カバレッジ確認
open coverage.html
```

**✅ Checkpoint:** テストが実装され、カバレッジが70%以上

---

### Step 8: フロントエンドの更新

**目的:** フロントエンドを新しいAPIレスポンス形式に対応させる

#### 8-1. 型定義の更新

**ファイル: `frontend/src/types/user.ts`**

```typescript
export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  department: string;
  role: 'admin' | 'manager' | 'user' | 'viewer';
  status: 'active' | 'inactive' | 'suspended';
  last_login: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  full_name?: string;
  department?: string;
  password: string;
  role: 'admin' | 'manager' | 'user' | 'viewer';
}

export interface UpdateUserRequest {
  email?: string;
  full_name?: string;
  department?: string;
  role?: 'admin' | 'manager' | 'user' | 'viewer';
  status?: 'active' | 'inactive' | 'suspended';
}

export interface PaginatedResponse<T> {
  code: number;
  message: string;
  data: T[];
  pagination: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export type UserStatus = 'active' | 'inactive' | 'suspended';
export type UserRole = 'admin' | 'manager' | 'user' | 'viewer';
```

#### 8-2. API関数の更新

**ファイル: `frontend/src/lib/users.ts`**

```typescript
import { api } from './api';
import type {
  User,
  CreateUserRequest,
  UpdateUserRequest,
  PaginatedResponse,
  ApiResponse,
} from '@/types/user';

export const usersApi = {
  async getUsers(page = 1, perPage = 10): Promise<PaginatedResponse<User>> {
    const response = await api.get<PaginatedResponse<User>>('/users', {
      params: { page, per_page: perPage },
    });
    return response.data;
  },

  async getUserById(id: number): Promise<User> {
    const response = await api.get<ApiResponse<{ user: User }>>(`/users/${id}`);
    return response.data.data.user;
  },

  async createUser(data: CreateUserRequest): Promise<User> {
    const response = await api.post<ApiResponse<{ user: User }>>('/users', data);
    return response.data.data.user;
  },

  async updateUser(id: number, data: UpdateUserRequest): Promise<User> {
    const response = await api.put<ApiResponse<{ user: User }>>(`/users/${id}`, data);
    return response.data.data.user;
  },

  async deleteUser(id: number): Promise<void> {
    await api.delete(`/users/${id}`);
  },
};
```

#### 8-3. カスタムフックの更新

**ファイル: `frontend/src/hooks/useUsers.ts`**

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { usersApi } from '@/lib/users';
import type { CreateUserRequest, UpdateUserRequest } from '@/types/user';

export function useUsers(page = 1, perPage = 10) {
  return useQuery({
    queryKey: ['users', page, perPage],
    queryFn: () => usersApi.getUsers(page, perPage),
  });
}

export function useUser(id: number) {
  return useQuery({
    queryKey: ['users', id],
    queryFn: () => usersApi.getUserById(id),
    enabled: !!id,
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateUserRequest) => usersApi.createUser(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
}

export function useUpdateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserRequest }) =>
      usersApi.updateUser(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      queryClient.invalidateQueries({ queryKey: ['users', variables.id] });
    },
  });
}

export function useDeleteUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => usersApi.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
}
```

#### 8-4. ビルドとテスト

```bash
cd frontend

# 型チェック
npm run type-check

# ビルド
npm run build

# テスト
npm test
```

**✅ Checkpoint:** フロントエンドが更新され、新しいAPIに対応している

---

### Step 9: 統合テスト

**目的:** 全体が正しく動作することを確認

#### 9-1. APIテスト

```bash
# Ping
curl http://localhost:8080/api/v1/ping
# {"message":"pong"}

# ユーザー一覧（ページネーション付き）
curl "http://localhost:8080/api/v1/users?page=1&per_page=10" | jq

# 期待される出力:
# {
#   "code": 200,
#   "message": "success",
#   "data": [
#     {
#       "id": 1,
#       "username": "admin",
#       "email": "admin@example.com",
#       "full_name": "管理者 太郎",
#       "department": "IT部",
#       "role": "admin",
#       "status": "active",
#       "last_login": null,
#       ...
#     }
#   ],
#   "pagination": {
#     "page": 1,
#     "per_page": 10,
#     "total": 5,
#     "total_pages": 1
#   }
# }

# ユーザー詳細
curl http://localhost:8080/api/v1/users/1 | jq

# ユーザー作成
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "full_name": "新規 ユーザー",
    "department": "開発部",
    "password": "password123",
    "role": "user"
  }' | jq

# ユーザー更新
curl -X PUT http://localhost:8080/api/v1/users/6 \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "更新 ユーザー",
    "status": "inactive"
  }' | jq

# ユーザー削除
curl -X DELETE http://localhost:8080/api/v1/users/6
```

#### 9-2. フロントエンドテスト

```bash
# ブラウザで確認
# http://localhost:3000/users

# 期待される動作:
# - ユーザー一覧が表示される
# - ページネーションが動作する
# - 新しいフィールド（full_name, department, status）が表示される
```

**✅ Checkpoint:** 全ての機能が正しく動作する

---

## 完了チェックリスト

Phase 1の実装が完了したことを確認してください:

### バックエンド

- [ ] スキーマが設計仕様に準拠している（full_name, department, password_hash, status, last_login）
- [ ] マイグレーションが正常に実行される
- [ ] シードデータが投入される
- [ ] モデル層が更新されている
- [ ] ユーティリティ関数（エラー、レスポンス、ページネーション）が実装されている
- [ ] リポジトリ層が完全に実装されている
- [ ] サービス層が完全に実装されている（ビジネスロジック、バリデーション、エラーハンドリング）
- [ ] ハンドラー層が完全に実装されている
- [ ] ルーティングが正しく設定されている
- [ ] 単体テストが実装されている（カバレッジ70%以上）
- [ ] `make test` が成功する
- [ ] `make lint` が成功する
- [ ] `make build` が成功する

### API

- [ ] `GET /api/v1/users` が動作する（ページネーション付き）
- [ ] `GET /api/v1/users/:id` が動作する
- [ ] `POST /api/v1/users` が動作する（バリデーション含む）
- [ ] `PUT /api/v1/users/:id` が動作する
- [ ] `DELETE /api/v1/users/:id` が動作する
- [ ] エラーレスポンスが統一されている
- [ ] ページネーション情報が正しく返される

### フロントエンド

- [ ] TypeScript型定義が更新されている
- [ ] API関数が新しいレスポンス形式に対応している
- [ ] カスタムフックが実装されている
- [ ] ユーザー一覧ページが動作する
- [ ] 新しいフィールドが表示される（full_name, department, status）
- [ ] ページネーションが動作する
- [ ] `npm run type-check` が成功する
- [ ] `npm run lint` が成功する
- [ ] `npm run build` が成功する
- [ ] `npm test` が成功する

### ドキュメント

- [ ] CHANGELOG.md に変更内容を記載した
- [ ] コードにコメントを追加した
- [ ] 必要に応じてREADME.mdを更新した

### Git

- [ ] 全ての変更をコミットした
- [ ] 意味のあるコミットメッセージを書いた（Conventional Commits形式）
- [ ] リモートにプッシュした
- [ ] Pull Requestを作成した（チーム開発の場合）

---

## 次のステップ

Phase 1が完了したら、**Phase 2: 認証・認可機能** の実装に進んでください。

詳細は **[IMPLEMENTATION_PHASES.md](IMPLEMENTATION_PHASES.md)** を参照してください。

---

**最終更新**: 2025-01-20
