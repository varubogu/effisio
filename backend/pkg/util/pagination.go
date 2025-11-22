package util

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams はページネーションパラメータです
type PaginationParams struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Offset  int `json:"-"`
}

// PaginationInfo はページネーション情報です
type PaginationInfo struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse はページネーション付きレスポンスです
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// GetPaginationParams はリクエストからページネーションパラメータを取得します
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

// NewPaginationInfo はページネーション情報を生成します
func NewPaginationInfo(total int64, params *PaginationParams) *PaginationInfo {
	totalPages := int(math.Ceil(float64(total) / float64(params.PerPage)))
	return &PaginationInfo{
		Page:       params.Page,
		PerPage:    params.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

// NewPaginatedResponse はページネーション付きレスポンスを生成します
func NewPaginatedResponse(data interface{}, total int64, params *PaginationParams) *PaginatedResponse {
	return &PaginatedResponse{
		Data:       data,
		Pagination: *NewPaginationInfo(total, params),
	}
}

// Paginated はページネーション付きレスポンスを返します
func Paginated(c *gin.Context, response *PaginatedResponse) {
	c.JSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"message":    "success",
		"data":       response.Data,
		"pagination": response.Pagination,
	})
}
