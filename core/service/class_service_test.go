package service_test

import (
	"errors"
	"testing"

	"alhikmah-attendance-api/core/domain"
	"alhikmah-attendance-api/core/mocks"
	"alhikmah-attendance-api/core/service"

	"github.com/stretchr/testify/assert"
)

func TestClassService_GetAll(t *testing.T) {
	mockRepo := new(mocks.ClassRepository)
	svc := service.NewClassService(mockRepo)

	classes := []*domain.Class{{ID: "1", ClassName: "Class A"}}
	mockRepo.On("GetAll", "t1", "2025/2026", 0, 10).Return(classes, 1, nil)

	res, count, err := svc.GetAll("t1", "2025/2026", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, classes, res)
	assert.Equal(t, 1, count)
	mockRepo.AssertExpectations(t)
}

func TestClassService_GetByID(t *testing.T) {
	mockRepo := new(mocks.ClassRepository)
	svc := service.NewClassService(mockRepo)

	class := &domain.Class{ID: "1", ClassName: "Class A"}
	mockRepo.On("GetByID", "1").Return(class, nil)

	res, err := svc.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, class, res)
	mockRepo.AssertExpectations(t)
}

func TestClassService_Create(t *testing.T) {
	mockRepo := new(mocks.ClassRepository)
	svc := service.NewClassService(mockRepo)

	class := &domain.Class{ClassName: "Class A", AcademicYear: "2025/2026", TeacherID: "t-1"}
	
	// Valid class
	mockRepo.On("Create", class).Return(nil)
	err := svc.Create(class)
	assert.NoError(t, err)

	// Missing Name
	err = svc.Create(&domain.Class{AcademicYear: "2025/2026", TeacherID: "t-1"})
	assert.Error(t, err)

	// Missing AcademicYear
	err = svc.Create(&domain.Class{ClassName: "Class A", TeacherID: "t-1"})
	assert.Error(t, err)

	// Missing TeacherID
	err = svc.Create(&domain.Class{ClassName: "Class A", AcademicYear: "2025/2026"})
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestClassService_Update(t *testing.T) {
	mockRepo := new(mocks.ClassRepository)
	svc := service.NewClassService(mockRepo)

	class := &domain.Class{ID: "1", ClassName: "Class A", AcademicYear: "2025/2026", TeacherID: "t-1"}
	
	// Valid update
	mockRepo.On("Update", class).Return(nil)
	err := svc.Update(class)
	assert.NoError(t, err)

	// Missing ID
	err = svc.Update(&domain.Class{ClassName: "Class A", AcademicYear: "2025/2026", TeacherID: "t-1"})
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestClassService_SoftDelete(t *testing.T) {
	mockRepo := new(mocks.ClassRepository)
	svc := service.NewClassService(mockRepo)

	mockRepo.On("GetByID", "1").Return(&domain.Class{StudentCount: 0}, nil)
	mockRepo.On("SoftDelete", "1").Return(nil)
	err := svc.SoftDelete("1")
	assert.NoError(t, err)

	mockRepo.On("GetByID", "2").Return((*domain.Class)(nil), errors.New("db error"))
	err = svc.SoftDelete("2")
	assert.Error(t, err)

	mockRepo.On("GetByID", "3").Return(&domain.Class{StudentCount: 5}, nil)
	err = svc.SoftDelete("3")
	assert.Error(t, err)
	assert.Equal(t, "cannot delete class with active students", err.Error())

	mockRepo.AssertExpectations(t)
}
