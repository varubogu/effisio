package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// MockAuditLogService mocks the AuditLogService
type MockAuditLogService struct {
	mock.Mock
}

func (m *MockAuditLogService) LogAction(ctx interface{}, req *model.CreateAuditLogRequest) (*model.AuditLogResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuditLogResponse), args.Error(1)
}

func (m *MockAuditLogService) GetByID(ctx interface{}, id uint) (*model.AuditLogResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuditLogResponse), args.Error(1)
}

func (m *MockAuditLogService) List(ctx interface{}, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*util.PaginatedResponse), args.Error(1)
}

func (m *MockAuditLogService) ListByUserID(ctx interface{}, userID uint, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	args := m.Called(ctx, userID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*util.PaginatedResponse), args.Error(1)
}

func (m *MockAuditLogService) ListByResource(ctx interface{}, resourceType, resourceID string, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	args := m.Called(ctx, resourceType, resourceID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*util.PaginatedResponse), args.Error(1)
}

func (m *MockAuditLogService) ListByAction(ctx interface{}, action string, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	args := m.Called(ctx, action, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*util.PaginatedResponse), args.Error(1)
}

func (m *MockAuditLogService) ListByDateRange(ctx interface{}, startDate, endDate time.Time, params *util.PaginationParams) (*util.PaginatedResponse, error) {
	args := m.Called(ctx, startDate, endDate, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*util.PaginatedResponse), args.Error(1)
}

func (m *MockAuditLogService) GetStatistics(ctx interface{}) (interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockAuditLogService) DeleteOldLogs(ctx interface{}, days int) error {
	return m.Called(ctx, days).Error(0)
}

func getHandlerLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func TestAuditLogHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuditLogService)

	logs := []*model.AuditLogResponse{
		{
			ID:           1,
			UserID:       1,
			Action:       model.ActionCreate,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   "user-1",
			Status:       model.AuditStatusSuccess,
			CreatedAt:    time.Now(),
		},
	}

	paginatedResp := &util.PaginatedResponse{
		Status:  200,
		Message: "OK",
		Data:    logs,
		Total:   1,
		Page:    1,
		PerPage: 10,
		Pages:   1,
	}

	mockService.On("List", mock.Anything, mock.Anything).Return(paginatedResp, nil)

	handler := NewAuditLogHandler(mockService, getHandlerLogger())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/audit-logs", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuditLogHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuditLogService)

	auditLog := &model.AuditLogResponse{
		ID:           1,
		UserID:       1,
		Action:       model.ActionCreate,
		ResourceType: model.ResourceTypeUser,
		ResourceID:   "user-1",
		Status:       model.AuditStatusSuccess,
		CreatedAt:    time.Now(),
	}

	mockService.On("GetByID", mock.Anything, uint(1)).Return(auditLog, nil)

	handler := NewAuditLogHandler(mockService, getHandlerLogger())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/audit-logs/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.GetByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuditLogHandler_GetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuditLogService)
	handler := NewAuditLogHandler(mockService, getHandlerLogger())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/audit-logs/invalid", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "invalid"})

	handler.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuditLogHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuditLogService)

	createReq := model.CreateAuditLogRequest{
		UserID:       1,
		Action:       model.ActionCreate,
		ResourceType: model.ResourceTypeUser,
		ResourceID:   "user-1",
		Changes: model.AuditLogChanges{
			Before: map[string]interface{}{},
			After: map[string]interface{}{
				"username": "test",
			},
		},
		Status: model.AuditStatusSuccess,
	}

	auditLog := &model.AuditLogResponse{
		ID:           1,
		UserID:       1,
		Action:       model.ActionCreate,
		ResourceType: model.ResourceTypeUser,
		ResourceID:   "user-1",
		Status:       model.AuditStatusSuccess,
		CreatedAt:    time.Now(),
	}

	mockService.On("LogAction", mock.Anything, mock.MatchedBy(func(req *model.CreateAuditLogRequest) bool {
		return req.UserID == 1 &&
			req.Action == model.ActionCreate &&
			req.ResourceType == model.ResourceTypeUser
	})).Return(auditLog, nil)

	handler := NewAuditLogHandler(mockService, getHandlerLogger())

	body, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/audit-logs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.RemoteIP = func() string { return "192.168.1.1" }

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuditLogHandler_GetStatistics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuditLogService)

	stats := &map[string]interface{}{
		"total_logs": int64(100),
		"by_action": map[string]int64{
			"create": 50,
			"update": 30,
		},
		"by_status": map[string]int64{
			"success": 95,
			"failed":  5,
		},
		"success_rate": 0.95,
	}

	mockService.On("GetStatistics", mock.Anything).Return(stats, nil)

	handler := NewAuditLogHandler(mockService, getHandlerLogger())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/audit-logs/statistics", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetStatistics(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuditLogHandler_DeleteOldLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuditLogService)
	mockService.On("DeleteOldLogs", mock.Anything, 90).Return(nil)

	handler := NewAuditLogHandler(mockService, getHandlerLogger())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/api/v1/audit-logs/delete-old?days=90", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.DeleteOldLogs(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuditLogHandler_ListByResource(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuditLogService)

	paginatedResp := &util.PaginatedResponse{
		Status:  200,
		Message: "OK",
		Data:    []*model.AuditLogResponse{},
		Total:   0,
	}

	mockService.On("ListByResource", mock.Anything, model.ResourceTypeUser, "user-123", mock.Anything).
		Return(paginatedResp, nil)

	handler := NewAuditLogHandler(mockService, getHandlerLogger())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/audit-logs/resource?resource_type=user&resource_id=user-123", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListByResource(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuditLogHandler_ListByAction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuditLogService)

	paginatedResp := &util.PaginatedResponse{
		Status:  200,
		Message: "OK",
		Data:    []*model.AuditLogResponse{},
		Total:   0,
	}

	mockService.On("ListByAction", mock.Anything, model.ActionCreate, mock.Anything).
		Return(paginatedResp, nil)

	handler := NewAuditLogHandler(mockService, getHandlerLogger())

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/audit-logs/action?action=create", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListByAction(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
