package util

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response は統一レスポンス形式です
type Response struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail はエラー詳細です
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success は成功レスポンスを返します (200 OK)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// Created は作成成功レスポンスを返します (201 Created)
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "created",
		Data:    data,
	})
}

// NoContent は内容なしレスポンスを返します (204 No Content)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error はエラーレスポンスを返します
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

// ValidationError はバリデーションエラーレスポンスを返します
func ValidationError(c *gin.Context, errors map[string]string) {
	Error(c, http.StatusBadRequest, ErrCodeValidationError, "Validation failed", errors)
}

// HandleError はAppErrorからレスポンスを生成します
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		Error(c, appErr.StatusCode, appErr.Code, appErr.Message, nil)
	} else {
		Error(c, http.StatusInternalServerError, ErrCodeInternalError, "Internal server error", nil)
	}
}

// ParseValidationErrors はバリデーションエラーをパースします
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
