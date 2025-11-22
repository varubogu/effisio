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
	"github.com/varubogu/effisio/backend/pkg/util"
	"go.uber.org/zap"
)

func TestUserService_List_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()
	params := &util.PaginationParams{Page: 1, PerPage: 10, Offset: 0}

	users := []*model.User{
		{
			ID:           1,
			Username:     "user1",
			Email:        "user1@example.com",
			PasswordHash: "hash1",
			Role:         "admin",
			Status:       model.UserStatusActive,
			FullName:     "User One",
		},
		{
			ID:           2,
			Username:     "user2",
			Email:        "user2@example.com",
			PasswordHash: "hash2",
			Role:         "user",
			Status:       model.UserStatusActive,
			FullName:     "User Two",
		},
	}

	mockRepo.On("FindAll", ctx, params).Return(users, int64(2), nil)

	resp, err := userService.List(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, int64(2), resp.Total)

	mockRepo.AssertExpectations(t)
}

func TestUserService_List_EmptyResult(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()
	params := &util.PaginationParams{Page: 1, PerPage: 10, Offset: 0}

	mockRepo.On("FindAll", ctx, params).Return([]*model.User{}, int64(0), nil)

	resp, err := userService.List(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Data, 0)
	assert.Equal(t, int64(0), resp.Total)

	mockRepo.AssertExpectations(t)
}

func TestUserService_List_DatabaseError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()
	params := &util.PaginationParams{Page: 1, PerPage: 10, Offset: 0}

	mockRepo.On("FindAll", ctx, params).Return(nil, int64(0), errors.New("database error"))

	resp, err := userService.List(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()
	user := &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         "admin",
		Status:       model.UserStatusActive,
		FullName:     "Test User",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(user, nil)

	resp, err := userService.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	resp, err := userService.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()

	req := &model.CreateUserRequest{
		Username:   "newuser",
		Email:      "newuser@example.com",
		Password:   "SecurePassword123!",
		FullName:   "New User",
		Department: "Engineering",
		Role:       "user",
	}

	mockRepo.On("ExistsByUsername", ctx, "newuser").Return(false, nil)
	mockRepo.On("ExistsByEmail", ctx, "newuser@example.com").Return(false, nil)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(u *model.User) bool {
		return u.Username == "newuser" &&
			u.Email == "newuser@example.com" &&
			u.FullName == "New User" &&
			u.Department == "Engineering" &&
			u.Role == "user" &&
			u.Status == model.UserStatusActive &&
			len(u.PasswordHash) > 0
	})).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*model.User)
		u.ID = 1
	})

	resp, err := userService.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "newuser", resp.Username)
	assert.Equal(t, "newuser@example.com", resp.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_UsernameDuplicate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()

	req := &model.CreateUserRequest{
		Username:   "existinguser",
		Email:      "newuser@example.com",
		Password:   "SecurePassword123!",
		FullName:   "New User",
		Department: "Engineering",
		Role:       "user",
	}

	mockRepo.On("ExistsByUsername", ctx, "existinguser").Return(true, nil)

	resp, err := userService.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_EmailDuplicate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()

	req := &model.CreateUserRequest{
		Username:   "newuser",
		Email:      "existing@example.com",
		Password:   "SecurePassword123!",
		FullName:   "New User",
		Department: "Engineering",
		Role:       "user",
	}

	mockRepo.On("ExistsByUsername", ctx, "newuser").Return(false, nil)
	mockRepo.On("ExistsByEmail", ctx, "existing@example.com").Return(true, nil)

	resp, err := userService.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()
	newEmail := "updated@example.com"
	newFullName := "Updated Name"
	newRole := "manager"

	user := &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "old@example.com",
		PasswordHash: "hash",
		Role:         "user",
		Status:       model.UserStatusActive,
		FullName:     "Old Name",
	}

	req := &model.UpdateUserRequest{
		Email:    &newEmail,
		FullName: &newFullName,
		Role:     &newRole,
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(user, nil)
	mockRepo.On("FindByEmail", ctx, newEmail).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == 1 &&
			u.Email == newEmail &&
			u.FullName == newFullName &&
			u.Role == newRole
	})).Return(nil)

	resp, err := userService.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, newEmail, resp.Email)
	assert.Equal(t, newFullName, resp.FullName)
	assert.Equal(t, newRole, resp.Role)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()
	newEmail := "updated@example.com"

	req := &model.UpdateUserRequest{
		Email: &newEmail,
	}

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	resp, err := userService.Update(ctx, 999, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_EmailConflict(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()
	newEmail := "taken@example.com"

	user := &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "old@example.com",
		PasswordHash: "hash",
		Role:         "user",
		Status:       model.UserStatusActive,
	}

	otherUser := &model.User{
		ID:    2,
		Email: "taken@example.com",
	}

	req := &model.UpdateUserRequest{
		Email: &newEmail,
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(user, nil)
	mockRepo.On("FindByEmail", ctx, newEmail).Return(otherUser, nil)

	resp, err := userService.Update(ctx, 1, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()
	user := &model.User{
		ID:       1,
		Username: "testuser",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(user, nil)
	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := userService.Delete(ctx, 1)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	err := userService.Delete(ctx, 999)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_DatabaseError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()
	user := &model.User{
		ID:       1,
		Username: "testuser",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(user, nil)
	mockRepo.On("Delete", ctx, uint(1)).Return(errors.New("database error"))

	err := userService.Delete(ctx, 1)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_PartialUpdate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo, getLogger())

	ctx := context.Background()

	user := &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "old@example.com",
		PasswordHash: "hash",
		Role:         "user",
		Status:       model.UserStatusActive,
		FullName:     "Old Name",
		Department:   "Old Dept",
	}

	// Only update status
	newStatus := model.UserStatusInactive
	req := &model.UpdateUserRequest{
		Status: &newStatus,
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(user, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == 1 &&
			u.Status == model.UserStatusInactive &&
			u.Email == "old@example.com" && // Unchanged
			u.FullName == "Old Name"        // Unchanged
	})).Return(nil)

	resp, err := userService.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, model.UserStatusInactive, resp.Status)

	mockRepo.AssertExpectations(t)
}
