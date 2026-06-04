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

func TestClassService_GetAll(t *testing.T) {
	mockRepo := new(mocks.ClassRepository)
	svc := service.NewClassService(mockRepo)

	classes := []*domain.Class{{ID: "1", RoomName: "Aqoba", Grade: 1}}
	mockRepo.On("GetAll", "t1", "2025/2026", 0, 10).Return(classes, 1, nil)

	res, count, err := svc.GetAll("t1", "2025/2026", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "1", res[0].ID)
	assert.Equal(t, 1, count)
	mockRepo.AssertExpectations(t)
}

func TestClassService_GetByID(t *testing.T) {
	mockRepo := new(mocks.ClassRepository)
	svc := service.NewClassService(mockRepo)

	class := &domain.Class{ID: "1", RoomName: "Aqoba", Grade: 1}
	mockRepo.On("GetByID", "1").Return(class, nil)

	res, err := svc.GetByID("1")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "1", res.ID)
	mockRepo.AssertExpectations(t)
}

func TestClassService_Create(t *testing.T) {
        mockRepo := new(mocks.ClassRepository)
        svc := service.NewClassService(mockRepo)

        // Valid class
        mockRepo.On("Create", mock.AnythingOfType("*domain.Class")).Return(nil)
        err := svc.Create("Aqoba", 1, nil, "t-1", "2025/2026", 30, "")
        assert.NoError(t, err)

        // Missing Name
        err = svc.Create("", 1, nil, "t-1", "2025/2026", 30, "")
        assert.Error(t, err)

        // Invalid Grade
        err = svc.Create("Aqoba", 13, nil, "t-1", "2025/2026", 30, "")
        assert.Error(t, err)

        // Missing AcademicYear
        err = svc.Create("Aqoba", 1, nil, "t-1", "", 30, "")
        assert.Error(t, err)

        // Missing TeacherID
        err = svc.Create("Aqoba", 1, nil, "", "2025/2026", 30, "")
        assert.Error(t, err)

        mockRepo.AssertExpectations(t)
}

func TestClassService_Update(t *testing.T) {
        mockRepo := new(mocks.ClassRepository)
        svc := service.NewClassService(mockRepo)

        // Valid update
        mockRepo.On("Update", mock.AnythingOfType("*domain.Class")).Return(nil)
        err := svc.Update("1", "Aqoba", 1, nil, "t-1", 30, "")
        assert.NoError(t, err)

        // Missing ID
        err = svc.Update("", "Aqoba", 1, nil, "t-1", 30, "")
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
