# API実装ガイド

このドキュメントでは、EffisioプロジェクトでAPIを実装する際の統一的な手法とベストプラクティスを定義します。

## 目次

- [レスポンス形式の統一](#レスポンス形式の統一)
- [エラーハンドリング](#エラーハンドリング)
- [ページネーション実装](#ページネーション実装)
- [バリデーション](#バリデーション)
- [ミドルウェア](#ミドルウェア)
- [実装パターン](#実装パターン)
- [テストの書き方](#テストの書き方)

---

## レスポンス形式の統一

### 標準レスポンス構造

全てのAPIエンドポイントは以下の統一形式を使用します：

```go
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
    Code    string       `json:"code"`
    Message string       `json:"message"`
    Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
```

### ヘルパー関数の実装

`internal/util/response.go` に以下のヘルパー関数を作成：

```go
package util

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// Success は成功レスポンスを返します
func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, gin.H{
        "code":    http.StatusOK,
        "message": "success",
        "data":    data,
    })
}

// Created はリソース作成成功レスポンスを返します
func Created(c *gin.Context, data interface{}) {
    c.JSON(http.StatusCreated, gin.H{
        "code":    http.StatusCreated,
        "message": "created",
        "data":    data,
    })
}

// NoContent は成功（データなし）レスポンスを返します
func NoContent(c *gin.Context) {
    c.JSON(http.StatusNoContent, gin.H{
        "code":    http.StatusNoContent,
        "message": "success",
    })
}

// Error はエラーレスポンスを返します
func Error(c *gin.Context, statusCode int, errorCode string, message string) {
    c.JSON(statusCode, gin.H{
        "code":    statusCode,
        "message": "error",
        "error": gin.H{
            "code":    errorCode,
            "message": message,
        },
    })
}

// ValidationError はバリデーションエラーレスポンスを返します
func ValidationError(c *gin.Context, errors []FieldError) {
    c.JSON(http.StatusUnprocessableEntity, gin.H{
        "code":    http.StatusUnprocessableEntity,
        "message": "validation error",
        "error": gin.H{
            "code":    "VALIDATION_ERROR",
            "message": "リクエストの検証に失敗しました",
            "details": errors,
        },
    })
}

// Paginated はページネーション付きレスポンスを返します
func Paginated(c *gin.Context, data interface{}, pagination *PaginationInfo) {
    c.JSON(http.StatusOK, gin.H{
        "code":       http.StatusOK,
        "message":    "success",
        "data":       data,
        "pagination": pagination,
    })
}
```

### 使用例

```go
// ハンドラーでの使用
func (h *UserHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        util.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
        return
    }

    user, err := h.service.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        util.Error(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
        return
    }

    util.Success(c, gin.H{"user": user})
}
```

---

## エラーハンドリング

### 標準エラーコード

`internal/util/errors.go` でエラーコードを定義：

```go
package util

const (
    // 認証エラー
    ErrCodeAuthFailed          = "AUTH_001"
    ErrCodeTokenInvalid        = "AUTH_002"
    ErrCodeTokenExpired        = "AUTH_003"
    ErrCodePermissionDenied    = "AUTH_004"

    // ユーザーエラー
    ErrCodeUserNotFound        = "USER_001"
    ErrCodeUsernameDuplicate   = "USER_002"
    ErrCodeEmailDuplicate      = "USER_003"

    // バリデーションエラー
    ErrCodeValidation          = "VALIDATION_001"
    ErrCodeInvalidFormat       = "VALIDATION_002"

    // サーバーエラー
    ErrCodeInternal            = "SERVER_001"
    ErrCodeDatabaseError       = "SERVER_002"
)

var ErrorMessages = map[string]string{
    ErrCodeAuthFailed:        "認証に失敗しました",
    ErrCodeTokenInvalid:      "トークンが無効です",
    ErrCodeTokenExpired:      "トークンの有効期限が切れています",
    ErrCodePermissionDenied:  "権限がありません",
    ErrCodeUserNotFound:      "ユーザーが見つかりません",
    ErrCodeUsernameDuplicate: "ユーザー名は既に使用されています",
    ErrCodeEmailDuplicate:    "メールアドレスは既に使用されています",
    ErrCodeValidation:        "入力値に誤りがあります",
    ErrCodeInvalidFormat:     "データ形式が正しくありません",
    ErrCodeInternal:          "サーバー内部エラーが発生しました",
    ErrCodeDatabaseError:     "データベースエラーが発生しました",
}

func GetErrorMessage(code string) string {
    if msg, ok := ErrorMessages[code]; ok {
        return msg
    }
    return "不明なエラーが発生しました"
}
```

### カスタムエラー型

```go
package util

import "fmt"

type AppError struct {
    Code       string
    Message    string
    StatusCode int
    Err        error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func NewAppError(statusCode int, code string, message string, err error) *AppError {
    return &AppError{
        Code:       code,
        Message:    message,
        StatusCode: statusCode,
        Err:        err,
    }
}

// 便利な関数
func NewNotFoundError(code string, message string) *AppError {
    return NewAppError(http.StatusNotFound, code, message, nil)
}

func NewBadRequestError(code string, message string) *AppError {
    return NewAppError(http.StatusBadRequest, code, message, nil)
}

func NewInternalError(code string, err error) *AppError {
    return NewAppError(http.StatusInternalServerError, code, GetErrorMessage(code), err)
}
```

### エラーハンドリングの統一

サービス層でのエラーハンドリング：

```go
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.UserResponse, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, util.NewNotFoundError(util.ErrCodeUserNotFound, "User not found")
        }
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    return user.ToResponse(), nil
}
```

ハンドラー層でのエラー処理：

```go
func (h *UserHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        util.Error(c, http.StatusBadRequest, util.ErrCodeInvalidFormat, "Invalid user ID")
        return
    }

    user, err := h.service.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        if appErr, ok := err.(*util.AppError); ok {
            util.Error(c, appErr.StatusCode, appErr.Code, appErr.Message)
            return
        }
        // 予期しないエラー
        h.logger.Error("Unexpected error", zap.Error(err))
        util.Error(c, http.StatusInternalServerError, util.ErrCodeInternal, "Internal server error")
        return
    }

    util.Success(c, gin.H{"user": user})
}
```

---

## ページネーション実装

### ページネーション構造体

`internal/util/pagination.go`:

```go
package util

import (
    "math"
    "strconv"
    "github.com/gin-gonic/gin"
)

type PaginationInfo struct {
    Page       int  `json:"page"`
    PerPage    int  `json:"per_page"`
    Total      int64 `json:"total"`
    TotalPages int  `json:"total_pages"`
    HasNext    bool `json:"has_next"`
    HasPrev    bool `json:"has_prev"`
}

type PaginationParams struct {
    Page    int
    PerPage int
    Offset  int
}

// GetPaginationParams はクエリパラメータからページネーション情報を取得
func GetPaginationParams(c *gin.Context) *PaginationParams {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

    // バリデーション
    if page < 1 {
        page = 1
    }
    if perPage < 1 {
        perPage = 20
    }
    if perPage > 100 {
        perPage = 100 // 最大100件まで
    }

    offset := (page - 1) * perPage

    return &PaginationParams{
        Page:    page,
        PerPage: perPage,
        Offset:  offset,
    }
}

// NewPaginationInfo はページネーション情報を生成
func NewPaginationInfo(page, perPage int, total int64) *PaginationInfo {
    totalPages := int(math.Ceil(float64(total) / float64(perPage)))

    return &PaginationInfo{
        Page:       page,
        PerPage:    perPage,
        Total:      total,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
        HasPrev:    page > 1,
    }
}
```

### リポジトリでの実装

```go
func (r *UserRepository) FindAll(ctx context.Context, params *util.PaginationParams) ([]*model.User, int64, error) {
    var users []*model.User
    var total int64

    // 総件数を取得
    if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // ページネーション付きでデータ取得
    if err := r.db.WithContext(ctx).
        Offset(params.Offset).
        Limit(params.PerPage).
        Order("created_at DESC").
        Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
}
```

### ハンドラーでの使用

```go
func (h *UserHandler) List(c *gin.Context) {
    // ページネーションパラメータ取得
    params := util.GetPaginationParams(c)

    // データ取得
    users, total, err := h.repo.FindAll(c.Request.Context(), params)
    if err != nil {
        util.Error(c, http.StatusInternalServerError, util.ErrCodeDatabaseError, "Failed to fetch users")
        return
    }

    // レスポンス変換
    responses := make([]*model.UserResponse, len(users))
    for i, user := range users {
        responses[i] = user.ToResponse()
    }

    // ページネーション情報生成
    pagination := util.NewPaginationInfo(params.Page, params.PerPage, total)

    // レスポンス返却
    util.Paginated(c, responses, pagination)
}
```

---

## バリデーション

### カスタムバリデータの登録

`internal/util/validator.go`:

```go
package util

import (
    "regexp"
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()

    // カスタムバリデータ登録
    validate.RegisterValidation("password", validatePassword)
    validate.RegisterValidation("username", validateUsername)
}

// validatePassword はパスワード強度をチェック
func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    if len(password) < 8 {
        return false
    }

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>_\-+=\[\]]`).MatchString(password)

    return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateUsername はユーザー名形式をチェック
func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
    return matched
}

// GetValidator はバリデータインスタンスを返す
func GetValidator() *validator.Validate {
    return validate
}

// ValidateStruct は構造体をバリデーション
func ValidateStruct(s interface{}) []FieldError {
    err := validate.Struct(s)
    if err == nil {
        return nil
    }

    var fieldErrors []FieldError
    for _, err := range err.(validator.ValidationErrors) {
        fieldErrors = append(fieldErrors, FieldError{
            Field:   err.Field(),
            Message: getValidationMessage(err),
        })
    }

    return fieldErrors
}

// getValidationMessage はバリデーションエラーメッセージを取得
func getValidationMessage(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return err.Field() + "は必須です"
    case "email":
        return "有効なメールアドレスを入力してください"
    case "min":
        return err.Field() + "は" + err.Param() + "文字以上である必要があります"
    case "max":
        return err.Field() + "は" + err.Param() + "文字以下である必要があります"
    case "password":
        return "パスワードは8文字以上で、英大小文字、数字、記号を含めてください"
    case "username":
        return "ユーザー名は英数字とアンダースコアのみ使用できます"
    default:
        return err.Field() + "の形式が正しくありません"
    }
}
```

### リクエスト構造体でのバリデーション使用

```go
type CreateUserRequest struct {
    Username   string `json:"username" binding:"required,min=3,max=50,username"`
    Email      string `json:"email" binding:"required,email,max=255"`
    Password   string `json:"password" binding:"required,password"`
    FullName   string `json:"full_name" binding:"omitempty,max=255"`
    Department string `json:"department" binding:"omitempty,max=100"`
    Role       string `json:"role" binding:"omitempty,oneof=admin manager user viewer"`
}
```

### ハンドラーでのバリデーション

```go
func (h *UserHandler) Create(c *gin.Context) {
    var req model.CreateUserRequest

    // JSONバインドとバリデーション
    if err := c.ShouldBindJSON(&req); err != nil {
        // Ginのバリデーションエラーを変換
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            fieldErrors := make([]util.FieldError, len(validationErrors))
            for i, err := range validationErrors {
                fieldErrors[i] = util.FieldError{
                    Field:   err.Field(),
                    Message: util.getValidationMessage(err),
                }
            }
            util.ValidationError(c, fieldErrors)
            return
        }
        util.Error(c, http.StatusBadRequest, util.ErrCodeInvalidFormat, "Invalid request format")
        return
    }

    // サービス層呼び出し
    user, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        // エラーハンドリング...
        return
    }

    util.Created(c, gin.H{"user": user})
}
```

---

## ミドルウェア

### 認証ミドルウェア（JWT）

`internal/middleware/auth.go`:

```go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "go.uber.org/zap"

    "github.com/varubogu/effisio/backend/internal/util"
)

type AuthMiddleware struct {
    jwtSecret []byte
    logger    *zap.Logger
}

func NewAuthMiddleware(jwtSecret string, logger *zap.Logger) *AuthMiddleware {
    return &AuthMiddleware{
        jwtSecret: []byte(jwtSecret),
        logger:    logger,
    }
}

// RequireAuth は認証必須のミドルウェア
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            util.Error(c, http.StatusUnauthorized, util.ErrCodeAuthFailed, "Authorization header required")
            c.Abort()
            return
        }

        // "Bearer {token}" 形式をパース
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            util.Error(c, http.StatusUnauthorized, util.ErrCodeTokenInvalid, "Invalid authorization format")
            c.Abort()
            return
        }

        tokenString := parts[1]

        // トークン検証
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return m.jwtSecret, nil
        })

        if err != nil || !token.Valid {
            util.Error(c, http.StatusUnauthorized, util.ErrCodeTokenInvalid, "Invalid token")
            c.Abort()
            return
        }

        // クレーム取得
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            util.Error(c, http.StatusUnauthorized, util.ErrCodeTokenInvalid, "Invalid token claims")
            c.Abort()
            return
        }

        // コンテキストにユーザー情報を設定
        c.Set("user_id", uint(claims["user_id"].(float64)))
        c.Set("username", claims["username"].(string))
        c.Set("role", claims["role"].(string))

        c.Next()
    }
}

// RequireRole はロール必須のミドルウェア
func (m *AuthMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            util.Error(c, http.StatusForbidden, util.ErrCodePermissionDenied, "Role not found")
            c.Abort()
            return
        }

        userRole := role.(string)
        for _, allowedRole := range allowedRoles {
            if userRole == allowedRole {
                c.Next()
                return
            }
        }

        util.Error(c, http.StatusForbidden, util.ErrCodePermissionDenied, "Insufficient permissions")
        c.Abort()
    }
}
```

### 使用例

```go
// ルーター設定で使用
api := router.Group("/api/v1")
{
    // 認証不要
    api.POST("/auth/login", authHandler.Login)

    // 認証必須
    authenticated := api.Group("")
    authenticated.Use(authMiddleware.RequireAuth())
    {
        authenticated.GET("/users", userHandler.List)
        authenticated.GET("/users/:id", userHandler.GetByID)

        // 管理者のみ
        admin := authenticated.Group("")
        admin.Use(authMiddleware.RequireRole("admin"))
        {
            admin.POST("/users", userHandler.Create)
            admin.DELETE("/users/:id", userHandler.Delete)
        }
    }
}
```

---

## 実装パターン

### 標準的なCRUD実装テンプレート

#### 1. モデル定義 (`internal/model/resource.go`)

```go
package model

import (
    "time"
    "gorm.io/gorm"
)

type Resource struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"not null;size:255" json:"name"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateResourceRequest struct {
    Name string `json:"name" binding:"required,max=255"`
}

type UpdateResourceRequest struct {
    Name string `json:"name" binding:"omitempty,max=255"`
}

type ResourceResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

func (r *Resource) ToResponse() *ResourceResponse {
    return &ResourceResponse{
        ID:        r.ID,
        Name:      r.Name,
        CreatedAt: r.CreatedAt,
    }
}
```

#### 2. リポジトリ (`internal/repository/resource.go`)

```go
package repository

import (
    "context"
    "gorm.io/gorm"
    "github.com/varubogu/effisio/backend/internal/model"
    "github.com/varubogu/effisio/backend/internal/util"
)

type ResourceRepository struct {
    db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) *ResourceRepository {
    return &ResourceRepository{db: db}
}

func (r *ResourceRepository) FindAll(ctx context.Context, params *util.PaginationParams) ([]*model.Resource, int64, error) {
    var resources []*model.Resource
    var total int64

    if err := r.db.WithContext(ctx).Model(&model.Resource{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := r.db.WithContext(ctx).
        Offset(params.Offset).
        Limit(params.PerPage).
        Order("created_at DESC").
        Find(&resources).Error; err != nil {
        return nil, 0, err
    }

    return resources, total, nil
}

func (r *ResourceRepository) FindByID(ctx context.Context, id uint) (*model.Resource, error) {
    var resource model.Resource
    if err := r.db.WithContext(ctx).First(&resource, id).Error; err != nil {
        return nil, err
    }
    return &resource, nil
}

func (r *ResourceRepository) Create(ctx context.Context, resource *model.Resource) error {
    return r.db.WithContext(ctx).Create(resource).Error
}

func (r *ResourceRepository) Update(ctx context.Context, resource *model.Resource) error {
    return r.db.WithContext(ctx).Save(resource).Error
}

func (r *ResourceRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&model.Resource{}, id).Error
}
```

#### 3. サービス (`internal/service/resource.go`)

```go
package service

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "go.uber.org/zap"

    "github.com/varubogu/effisio/backend/internal/model"
    "github.com/varubogu/effisio/backend/internal/repository"
    "github.com/varubogu/effisio/backend/internal/util"
)

type ResourceService struct {
    repo   *repository.ResourceRepository
    logger *zap.Logger
}

func NewResourceService(repo *repository.ResourceRepository, logger *zap.Logger) *ResourceService {
    return &ResourceService{
        repo:   repo,
        logger: logger,
    }
}

func (s *ResourceService) List(ctx context.Context, params *util.PaginationParams) ([]*model.ResourceResponse, *util.PaginationInfo, error) {
    resources, total, err := s.repo.FindAll(ctx, params)
    if err != nil {
        return nil, nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    responses := make([]*model.ResourceResponse, len(resources))
    for i, resource := range resources {
        responses[i] = resource.ToResponse()
    }

    pagination := util.NewPaginationInfo(params.Page, params.PerPage, total)

    return responses, pagination, nil
}

func (s *ResourceService) GetByID(ctx context.Context, id uint) (*model.ResourceResponse, error) {
    resource, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, util.NewNotFoundError("RESOURCE_NOT_FOUND", "Resource not found")
        }
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    return resource.ToResponse(), nil
}

func (s *ResourceService) Create(ctx context.Context, req *model.CreateResourceRequest) (*model.ResourceResponse, error) {
    resource := &model.Resource{
        Name: req.Name,
    }

    if err := s.repo.Create(ctx, resource); err != nil {
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    s.logger.Info("Resource created", zap.Uint("id", resource.ID))

    return resource.ToResponse(), nil
}

func (s *ResourceService) Update(ctx context.Context, id uint, req *model.UpdateResourceRequest) (*model.ResourceResponse, error) {
    resource, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, util.NewNotFoundError("RESOURCE_NOT_FOUND", "Resource not found")
        }
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    if req.Name != "" {
        resource.Name = req.Name
    }

    if err := s.repo.Update(ctx, resource); err != nil {
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    s.logger.Info("Resource updated", zap.Uint("id", resource.ID))

    return resource.ToResponse(), nil
}

func (s *ResourceService) Delete(ctx context.Context, id uint) error {
    if err := s.repo.Delete(ctx, id); err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return util.NewNotFoundError("RESOURCE_NOT_FOUND", "Resource not found")
        }
        return util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    s.logger.Info("Resource deleted", zap.Uint("id", id))

    return nil
}
```

#### 4. ハンドラー (`internal/handler/resource.go`)

```go
package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/varubogu/effisio/backend/internal/model"
    "github.com/varubogu/effisio/backend/internal/service"
    "github.com/varubogu/effisio/backend/internal/util"
)

type ResourceHandler struct {
    service *service.ResourceService
    logger  *zap.Logger
}

func NewResourceHandler(service *service.ResourceService, logger *zap.Logger) *ResourceHandler {
    return &ResourceHandler{
        service: service,
        logger:  logger,
    }
}

func (h *ResourceHandler) List(c *gin.Context) {
    params := util.GetPaginationParams(c)

    resources, pagination, err := h.service.List(c.Request.Context(), params)
    if err != nil {
        h.handleError(c, err)
        return
    }

    util.Paginated(c, resources, pagination)
}

func (h *ResourceHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        util.Error(c, http.StatusBadRequest, util.ErrCodeInvalidFormat, "Invalid resource ID")
        return
    }

    resource, err := h.service.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        h.handleError(c, err)
        return
    }

    util.Success(c, gin.H{"resource": resource})
}

func (h *ResourceHandler) Create(c *gin.Context) {
    var req model.CreateResourceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        util.ValidationError(c, util.ParseValidationErrors(err))
        return
    }

    resource, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    util.Created(c, gin.H{"resource": resource})
}

func (h *ResourceHandler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        util.Error(c, http.StatusBadRequest, util.ErrCodeInvalidFormat, "Invalid resource ID")
        return
    }

    var req model.UpdateResourceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        util.ValidationError(c, util.ParseValidationErrors(err))
        return
    }

    resource, err := h.service.Update(c.Request.Context(), uint(id), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    util.Success(c, gin.H{"resource": resource})
}

func (h *ResourceHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        util.Error(c, http.StatusBadRequest, util.ErrCodeInvalidFormat, "Invalid resource ID")
        return
    }

    if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
        h.handleError(c, err)
        return
    }

    util.Success(c, gin.H{"message": "Resource deleted successfully"})
}

// handleError は共通のエラーハンドリング
func (h *ResourceHandler) handleError(c *gin.Context, err error) {
    if appErr, ok := err.(*util.AppError); ok {
        util.Error(c, appErr.StatusCode, appErr.Code, appErr.Message)
        return
    }

    h.logger.Error("Unexpected error", zap.Error(err))
    util.Error(c, http.StatusInternalServerError, util.ErrCodeInternal, "Internal server error")
}
```

---

## テストの書き方

### ユニットテストの基本

```go
package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go.uber.org/zap"

    "github.com/varubogu/effisio/backend/internal/model"
    "github.com/varubogu/effisio/backend/internal/service"
)

// モックリポジトリ
type MockResourceRepository struct {
    mock.Mock
}

func (m *MockResourceRepository) FindByID(ctx context.Context, id uint) (*model.Resource, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.Resource), args.Error(1)
}

func TestResourceService_GetByID(t *testing.T) {
    logger, _ := zap.NewDevelopment()
    mockRepo := new(MockResourceRepository)
    service := service.NewResourceService(mockRepo, logger)

    t.Run("成功ケース", func(t *testing.T) {
        // モックの設定
        expectedResource := &model.Resource{
            ID:   1,
            Name: "Test Resource",
        }
        mockRepo.On("FindByID", mock.Anything, uint(1)).Return(expectedResource, nil)

        // テスト実行
        result, err := service.GetByID(context.Background(), 1)

        // 検証
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Equal(t, "Test Resource", result.Name)
        mockRepo.AssertExpectations(t)
    })

    t.Run("リソースが見つからない", func(t *testing.T) {
        mockRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

        result, err := service.GetByID(context.Background(), 999)

        assert.Error(t, err)
        assert.Nil(t, result)
        mockRepo.AssertExpectations(t)
    })
}
```

---

## まとめ

### 実装チェックリスト

新しいAPIエンドポイントを実装する際は、以下を確認してください：

- [ ] 統一されたレスポンス形式を使用している
- [ ] エラーコードを定義している
- [ ] バリデーションを実装している
- [ ] ページネーションが必要な場合は実装している
- [ ] 適切なHTTPステータスコードを返している
- [ ] ログを適切に記録している
- [ ] ユニットテストを書いている
- [ ] API仕様書を更新している
- [ ] エラーハンドリングを統一している
- [ ] トランザクション処理が必要な場合は実装している
