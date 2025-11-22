package util

import (
	"fmt"
	"net/http"
)

// エラーコード定数
const (
	// 認証エラー (AUTH_xxx)
	ErrCodeUnauthorized           = "AUTH_001"
	ErrCodeInvalidToken           = "AUTH_002"
	ErrCodeTokenExpired           = "AUTH_003"
	ErrCodeInsufficientPermission = "AUTH_004"

	// ユーザーエラー (USER_xxx)
	ErrCodeUserNotFound      = "USER_001"
	ErrCodeUserAlreadyExists = "USER_002"
	ErrCodeInvalidCredentials = "USER_003"

	// バリデーションエラー (VAL_xxx)
	ErrCodeValidationError  = "VAL_001"
	ErrCodeInvalidParameter = "VAL_002"

	// データベースエラー (DB_xxx)
	ErrCodeDatabaseError  = "DB_001"
	ErrCodeRecordNotFound = "DB_002"

	// システムエラー (SYS_xxx)
	ErrCodeInternalError     = "SYS_001"
	ErrCodePasswordHashError = "SYS_002"
)

// AppError はアプリケーションエラーを表します
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

// Error はerrorインターフェースを実装します
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewBadRequestError は400エラーを作成します
func NewBadRequestError(code string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    "Bad request",
		StatusCode: http.StatusBadRequest,
		Err:        err,
	}
}

// NewUnauthorizedError は401エラーを作成します
func NewUnauthorizedError(code string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    "Unauthorized",
		StatusCode: http.StatusUnauthorized,
		Err:        err,
	}
}

// NewForbiddenError は403エラーを作成します
func NewForbiddenError(code string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    "Forbidden",
		StatusCode: http.StatusForbidden,
		Err:        err,
	}
}

// NewNotFoundError は404エラーを作成します
func NewNotFoundError(code string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    "Resource not found",
		StatusCode: http.StatusNotFound,
		Err:        err,
	}
}

// NewConflictError は409エラーを作成します
func NewConflictError(code string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    "Resource conflict",
		StatusCode: http.StatusConflict,
		Err:        err,
	}
}

// NewInternalError は500エラーを作成します
func NewInternalError(code string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
