package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/pkg/util"
	"go.uber.org/zap"
)

// MockAuditLogRepository mocks the AuditLogRepository
type MockAuditLogRepository struct {
	mock.Mock
}

func (m *MockAuditLogRepository) Create(ctx context.Context, auditLog *model.AuditLog) error {
	return m.Called(ctx, auditLog).Error(0)
}

func (m *MockAuditLogRepository) FindByID(ctx context.Context, id uint) (*model.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) FindAll(ctx context.Context, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogRepository) FindByUserID(ctx context.Context, userID uint, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	args := m.Called(ctx, userID, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogRepository) FindByResourceID(ctx context.Context, resourceType, resourceID string, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	args := m.Called(ctx, resourceType, resourceID, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogRepository) FindByAction(ctx context.Context, action string, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	args := m.Called(ctx, action, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time, params *util.PaginationParams) ([]*model.AuditLog, int64, error) {
	args := m.Called(ctx, startDate, endDate, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogRepository) DeleteOldLogs(ctx context.Context, days int) error {
	return m.Called(ctx, days).Error(0)
}

func (m *MockAuditLogRepository) CountByAction(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockAuditLogRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func getAuditLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func TestAuditLogService_LogAction_Success(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()

	req := &model.CreateAuditLogRequest{
		UserID:       1,
		Action:       model.ActionCreate,
		ResourceType: model.ResourceTypeUser,
		ResourceID:   "user-123",
		Changes: model.AuditLogChanges{
			Before: map[string]interface{}{},
			After: map[string]interface{}{
				"username": "newuser",
				"email":    "new@example.com",
			},
		},
		IPAddress:    "192.168.1.1",
		UserAgent:    "Mozilla/5.0",
		Status:       model.AuditStatusSuccess,
		ErrorMessage: "",
	}

	mockRepo.On("Create", ctx, mock.MatchedBy(func(log *model.AuditLog) bool {
		return log.UserID == 1 &&
			log.Action == model.ActionCreate &&
			log.ResourceType == model.ResourceTypeUser &&
			log.Status == model.AuditStatusSuccess
	})).Run(func(args mock.Arguments) {
		log := args.Get(1).(*model.AuditLog)
		log.ID = 1
	}).Return(nil)

	resp, err := service.LogAction(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.UserID)
	assert.Equal(t, model.ActionCreate, resp.Action)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_LogAction_ValidationError(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()

	tests := []struct {
		name      string
		request   *model.CreateAuditLogRequest
		expectErr bool
	}{
		{
			name: "Missing UserID",
			request: &model.CreateAuditLogRequest{
				UserID:       0,
				Action:       model.ActionCreate,
				ResourceType: model.ResourceTypeUser,
				ResourceID:   "123",
			},
			expectErr: true,
		},
		{
			name: "Missing Action",
			request: &model.CreateAuditLogRequest{
				UserID:       1,
				Action:       "",
				ResourceType: model.ResourceTypeUser,
				ResourceID:   "123",
			},
			expectErr: true,
		},
		{
			name: "Invalid Action",
			request: &model.CreateAuditLogRequest{
				UserID:       1,
				Action:       "invalid_action",
				ResourceType: model.ResourceTypeUser,
				ResourceID:   "123",
				Status:       model.AuditStatusSuccess,
			},
			expectErr: true,
		},
		{
			name: "Missing ResourceType",
			request: &model.CreateAuditLogRequest{
				UserID:       1,
				Action:       model.ActionCreate,
				ResourceType: "",
				ResourceID:   "123",
			},
			expectErr: true,
		},
		{
			name: "Missing ResourceID",
			request: &model.CreateAuditLogRequest{
				UserID:       1,
				Action:       model.ActionCreate,
				ResourceType: model.ResourceTypeUser,
				ResourceID:   "",
			},
			expectErr: true,
		},
		{
			name: "Invalid Status",
			request: &model.CreateAuditLogRequest{
				UserID:       1,
				Action:       model.ActionCreate,
				ResourceType: model.ResourceTypeUser,
				ResourceID:   "123",
				Status:       "invalid_status",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.LogAction(ctx, tt.request)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	}
}

func TestAuditLogService_GetByID(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()

	auditLog := &model.AuditLog{
		ID:           1,
		UserID:       1,
		Action:       model.ActionCreate,
		ResourceType: model.ResourceTypeUser,
		ResourceID:   "user-123",
		Changes:      []byte(`{"before":{},"after":{"username":"test"}}`),
		Status:       model.AuditStatusSuccess,
		CreatedAt:    time.Now(),
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(auditLog, nil)

	resp, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	resp, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_List(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()
	params := &util.PaginationParams{Page: 1, PerPage: 10, Offset: 0}

	logs := []*model.AuditLog{
		{
			ID:           1,
			UserID:       1,
			Action:       model.ActionCreate,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   "user-1",
			Changes:      []byte(`{"before":{},"after":{}}`),
			Status:       model.AuditStatusSuccess,
			CreatedAt:    time.Now(),
		},
		{
			ID:           2,
			UserID:       2,
			Action:       model.ActionUpdate,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   "user-2",
			Changes:      []byte(`{"before":{},"after":{}}`),
			Status:       model.AuditStatusSuccess,
			CreatedAt:    time.Now(),
		},
	}

	mockRepo.On("FindAll", ctx, params).Return(logs, int64(2), nil)

	resp, err := service.List(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, int64(2), resp.Total)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_ListByUserID(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()
	params := &util.PaginationParams{Page: 1, PerPage: 10, Offset: 0}

	logs := []*model.AuditLog{
		{
			ID:           1,
			UserID:       1,
			Action:       model.ActionCreate,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   "user-1",
			Changes:      []byte(`{"before":{},"after":{}}`),
			Status:       model.AuditStatusSuccess,
		},
	}

	mockRepo.On("FindByUserID", ctx, uint(1), params).Return(logs, int64(1), nil)

	resp, err := service.ListByUserID(ctx, 1, params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Data, 1)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_ListByResource(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()
	params := &util.PaginationParams{Page: 1, PerPage: 10, Offset: 0}

	logs := []*model.AuditLog{
		{
			ID:           1,
			UserID:       1,
			Action:       model.ActionCreate,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   "user-123",
			Changes:      []byte(`{"before":{},"after":{}}`),
			Status:       model.AuditStatusSuccess,
		},
	}

	mockRepo.On("FindByResourceID", ctx, model.ResourceTypeUser, "user-123", params).
		Return(logs, int64(1), nil)

	resp, err := service.ListByResource(ctx, model.ResourceTypeUser, "user-123", params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Data, 1)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_ListByResource_MissingParams(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()
	params := &util.PaginationParams{}

	resp, err := service.ListByResource(ctx, "", "user-123", params)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAuditLogService_GetStatistics(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()

	byAction := map[string]int64{
		model.ActionCreate: 10,
		model.ActionUpdate: 5,
		model.ActionDelete: 2,
	}

	byStatus := map[string]int64{
		model.AuditStatusSuccess: 16,
		model.AuditStatusFailed:  1,
	}

	mockRepo.On("CountByAction", ctx).Return(byAction, nil)
	mockRepo.On("CountByStatus", ctx).Return(byStatus, nil)

	stats, err := service.GetStatistics(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(17), stats.TotalLogs)
	assert.Equal(t, byAction, stats.ByAction)
	assert.Equal(t, byStatus, stats.ByStatus)
	assert.InDelta(t, float64(16)/float64(17), stats.SuccessRate, 0.01)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_DeleteOldLogs(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()

	mockRepo.On("DeleteOldLogs", ctx, 90).Return(nil)

	err := service.DeleteOldLogs(ctx, 90)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_DeleteOldLogs_InvalidDays(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()

	err := service.DeleteOldLogs(ctx, 0)

	assert.Error(t, err)
}

func TestAuditLogService_ListByAction(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()
	params := &util.PaginationParams{Page: 1, PerPage: 10, Offset: 0}

	logs := []*model.AuditLog{
		{
			ID:           1,
			UserID:       1,
			Action:       model.ActionCreate,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   "user-1",
			Changes:      []byte(`{"before":{},"after":{}}`),
			Status:       model.AuditStatusSuccess,
		},
	}

	mockRepo.On("FindByAction", ctx, model.ActionCreate, params).Return(logs, int64(1), nil)

	resp, err := service.ListByAction(ctx, model.ActionCreate, params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Data, 1)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_ListByDateRange(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()
	params := &util.PaginationParams{Page: 1, PerPage: 10, Offset: 0}
	startDate := time.Now().AddDate(0, 0, -7)
	endDate := time.Now()

	logs := []*model.AuditLog{
		{
			ID:           1,
			UserID:       1,
			Action:       model.ActionCreate,
			ResourceType: model.ResourceTypeUser,
			ResourceID:   "user-1",
			Changes:      []byte(`{"before":{},"after":{}}`),
			Status:       model.AuditStatusSuccess,
			CreatedAt:    time.Now(),
		},
	}

	mockRepo.On("FindByDateRange", ctx, startDate, endDate, params).Return(logs, int64(1), nil)

	resp, err := service.ListByDateRange(ctx, startDate, endDate, params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Data, 1)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_ListByDateRange_InvalidRange(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	service := NewAuditLogService(mockRepo, getAuditLogger())

	ctx := context.Background()
	params := &util.PaginationParams{}
	startDate := time.Now()
	endDate := time.Now().AddDate(0, 0, -7) // endDate before startDate

	resp, err := service.ListByDateRange(ctx, startDate, endDate, params)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
