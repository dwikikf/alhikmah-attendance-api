package service_test

import (
	"errors"
	"testing"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/mocks"
	"alhikmah-attendance-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCreate_MissingRequiredFields(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	err := svc.Create("testuser", "", "", "", "")
	assert.Error(t, err)
	assert.Equal(t, "missing required fields", err.Error())
}

func TestUserCreate_InvalidRole(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	err := svc.Create("testuser", "test@example.com", "password123", "Test User", "student")
	assert.Error(t, err)
	assert.Equal(t, "invalid role", err.Error())
}

func TestUserCreate_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	err := svc.Create("testuser", "test@example.com", "password123", "Test User", "admin")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUpdate_MissingID(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	err := svc.Update("", "Test User Updated", "test@example.com", "")
	assert.Error(t, err)
	assert.Equal(t, "missing user ID", err.Error())
}

func TestUserUpdate_MissingRequiredFields(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	err := svc.Update("user-1", "", "test@example.com", "")
	assert.Error(t, err)
	assert.Equal(t, "missing required fields", err.Error())
}

func TestUserUpdate_SuccessWithPassword(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	err := svc.Update("user-1", "Test User Updated", "test@example.com", "newpassword123")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUpdate_SuccessWithoutPassword(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	err := svc.Update("user-1", "Test User Updated", "test@example.com", "")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAll_PaginationDefaults(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	// Send page=0, limit=200, should default to page=1, limit=10
	// offset = (1-1)*10 = 0
	mockRepo.On("GetAll", "teacher", (*bool)(nil), 0, 10).Return([]*domain.User{}, 0, nil)

	_, _, err := svc.GetAll("teacher", nil, 0, 200)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	mockUser := &domain.User{ID: "user-1"}
	mockRepo.On("GetByID", "user-1").Return(mockUser, nil)

	res, err := svc.GetByID("user-1")
	assert.NoError(t, err)
	assert.Equal(t, "user-1", res.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetByID_Error(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	mockRepo.On("GetByID", "user-1").Return((*domain.User)(nil), errors.New("not found"))

	res, err := svc.GetByID("user-1")
	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestSoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	mockRepo.On("SoftDelete", "user-1").Return(nil)

	err := svc.SoftDelete("user-1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSoftDelete_Error(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := service.NewUserService(mockRepo)

	mockRepo.On("SoftDelete", "user-1").Return(errors.New("db error"))

	err := svc.SoftDelete("user-1")
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
