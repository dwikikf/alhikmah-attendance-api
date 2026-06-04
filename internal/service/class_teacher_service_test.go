package service_test

import (
	"testing"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/mocks"
	"alhikmah-attendance-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClassTeacherService_Assign(t *testing.T) {
	mockRepo := new(mocks.ClassTeacherRepository)
	mockUserRepo := new(mocks.UserRepository)
	mockClassRepo := new(mocks.ClassRepository)
	svc := service.NewClassTeacherService(mockRepo, mockUserRepo, mockClassRepo)

	// Fail: empty subject
	err := svc.Assign("u1", "c1", "")
	assert.Error(t, err)

	// Fail: user not teacher
	mockUserRepo.On("GetByID", "u1").Return(&domain.User{ID: "u1", Role: "admin"}, nil).Once()
	err = svc.Assign("u1", "c1", "PJOK")
	assert.Error(t, err)

	// Success
	mockUserRepo.On("GetByID", "u2").Return(&domain.User{ID: "u2", Role: "teacher"}, nil).Once()
	mockClassRepo.On("GetByID", "c1").Return(&domain.Class{ID: "c1", AcademicYear: "2026/2027"}, nil).Once()
	mockRepo.On("Assign", mock.AnythingOfType("*domain.ClassTeacher")).Return(nil).Once()

	err = svc.Assign("u2", "c1", "PJOK")
	assert.NoError(t, err)
}
