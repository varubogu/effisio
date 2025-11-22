package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
	"go.uber.org/zap"
)

// MockUserRepository mocks the UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(ctx context.Context, params *util.PaginationParams) ([]*model.User, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) CountByRole(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockUserRepository) CountByDepartment(ctx context.Context) ([]struct {
	Department string
	Count      int64
}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]struct {
		Department string
		Count      int64
	}), args.Error(1)
}

func (m *MockUserRepository) GetLastLoginStats(ctx context.Context, days int) ([]struct {
	Date  string
	Count int64
}, error) {
	args := m.Called(ctx, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]struct {
		Date  string
		Count int64
	}), args.Error(1)
}

// MockRefreshTokenRepository mocks the RefreshTokenRepository
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	return m.Called(ctx, token).Error(0)
}

func (m *MockRefreshTokenRepository) FindByTokenID(ctx context.Context, tokenID string) (*model.RefreshToken, error) {
	args := m.Called(ctx, tokenID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) FindByUserID(ctx context.Context, userID uint) ([]*model.RefreshToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, tokenID string) error {
	return m.Called(ctx, tokenID).Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func getHashedPassword(t *testing.T, password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

func getLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func TestAuthServiceLogin_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()
	hashedPassword := getHashedPassword(t, "password123")

	user := &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         "admin",
		Status:       model.UserStatusActive,
		FullName:     "Test User",
	}

	mockUserRepo.On("FindByUsername", ctx, "testuser").Return(user, nil)
	mockUserRepo.On("Update", ctx, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == 1 && u.LastLogin != nil
	})).Return(nil)
	mockTokenRepo.On("Create", ctx, mock.MatchedBy(func(t *model.RefreshToken) bool {
		return t.UserID == 1 && !t.Revoked
	})).Return(nil)

	req := &LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := authService.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.NotNil(t, resp.User)
	assert.Equal(t, "testuser", resp.User.Username)

	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthServiceLogin_InvalidUsername(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	mockUserRepo.On("FindByUsername", ctx, "nonexistent").Return(nil, gorm.ErrRecordNotFound)

	req := &LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	resp, err := authService.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthServiceLogin_InvalidPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()
	hashedPassword := getHashedPassword(t, "correctpassword")

	user := &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         "admin",
		Status:       model.UserStatusActive,
	}

	mockUserRepo.On("FindByUsername", ctx, "testuser").Return(user, nil)

	req := &LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	resp, err := authService.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthServiceLogin_InactiveUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()
	hashedPassword := getHashedPassword(t, "password123")

	user := &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         "admin",
		Status:       model.UserStatusInactive, // Inactive status
	}

	mockUserRepo.On("FindByUsername", ctx, "testuser").Return(user, nil)

	req := &LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := authService.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthServiceRefreshToken_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	// Generate valid refresh token
	refreshToken, err := jwtService.GenerateRefreshToken(1, "token-123")
	require.NoError(t, err)

	user := &model.User{
		ID:       1,
		Username: "testuser",
		Role:     "admin",
		Status:   model.UserStatusActive,
	}

	dbToken := &model.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenID:   "token-123",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	mockTokenRepo.On("FindByTokenID", ctx, "token-123").Return(dbToken, nil)
	mockUserRepo.On("FindByID", ctx, uint(1)).Return(user, nil)
	mockTokenRepo.On("Revoke", ctx, "token-123").Return(nil)
	mockTokenRepo.On("Create", ctx, mock.MatchedBy(func(t *model.RefreshToken) bool {
		return t.UserID == 1 && !t.Revoked
	})).Return(nil)

	req := &RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	resp, err := authService.RefreshToken(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)

	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthServiceRefreshToken_InvalidToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	req := &RefreshTokenRequest{
		RefreshToken: "invalid-token",
	}

	resp, err := authService.RefreshToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAuthServiceRefreshToken_RevokedToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	// Generate valid token
	refreshToken, err := jwtService.GenerateRefreshToken(1, "token-123")
	require.NoError(t, err)

	// But in DB it's revoked
	dbToken := &model.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenID:   "token-123",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   true, // Revoked
	}

	mockTokenRepo.On("FindByTokenID", ctx, "token-123").Return(dbToken, nil)

	req := &RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	resp, err := authService.RefreshToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthServiceLogout_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	// Generate valid token
	refreshToken, err := jwtService.GenerateRefreshToken(1, "token-123")
	require.NoError(t, err)

	mockTokenRepo.On("Revoke", ctx, "token-123").Return(nil)

	err = authService.Logout(ctx, refreshToken)

	assert.NoError(t, err)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthServiceLogout_InvalidToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	// Logout with invalid token should still succeed
	err := authService.Logout(ctx, "invalid-token")

	assert.NoError(t, err)
}

func TestAuthServiceLogoutAll_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	mockTokenRepo.On("RevokeAllByUserID", ctx, uint(1)).Return(nil)

	err := authService.LogoutAll(ctx, 1)

	assert.NoError(t, err)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthServiceLogoutAll_DatabaseError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	mockTokenRepo.On("RevokeAllByUserID", ctx, uint(1)).Return(errors.New("database error"))

	err := authService.LogoutAll(ctx, 1)

	assert.Error(t, err)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthServiceRefreshToken_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	refreshToken, err := jwtService.GenerateRefreshToken(1, "token-123")
	require.NoError(t, err)

	dbToken := &model.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenID:   "token-123",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	mockTokenRepo.On("FindByTokenID", ctx, "token-123").Return(dbToken, nil)
	mockUserRepo.On("FindByID", ctx, uint(1)).Return(nil, gorm.ErrRecordNotFound)

	req := &RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	resp, err := authService.RefreshToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthServiceRefreshToken_InactiveUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockRefreshTokenRepository)
	jwtService := util.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	authService := NewAuthService(mockUserRepo, mockTokenRepo, jwtService, getLogger())

	ctx := context.Background()

	refreshToken, err := jwtService.GenerateRefreshToken(1, "token-123")
	require.NoError(t, err)

	user := &model.User{
		ID:       1,
		Username: "testuser",
		Role:     "admin",
		Status:   model.UserStatusInactive, // Inactive
	}

	dbToken := &model.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenID:   "token-123",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	mockTokenRepo.On("FindByTokenID", ctx, "token-123").Return(dbToken, nil)
	mockUserRepo.On("FindByID", ctx, uint(1)).Return(user, nil)

	req := &RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	resp, err := authService.RefreshToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockUserRepo.AssertExpectations(t)
}
