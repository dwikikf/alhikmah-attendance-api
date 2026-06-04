package service_test

import (
	"testing"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/mocks"
	"alhikmah-attendance-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStudentService_GenerateQRCodeData(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	res := svc.GenerateQRCodeData("123", "John", "10A")
	assert.Equal(t, "123|John|10A", res)
}

func TestStudentService_Create(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*domain.Student")).Return(nil)

	err := svc.Create("123", "John", "c1", nil, nil)
	assert.NoError(t, err)

	// Missing fields
	err = svc.Create("", "John", "c1", nil, nil)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestStudentService_CreateBulk(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	mockRepo.On("CreateBulk", mock.AnythingOfType("[]*domain.Student")).Return(nil)

	err := svc.CreateBulkFromRaw([]domain.StudentRawEntry{
		{NISN: "123", FullName: "John", ClassID: "c1", ClassName: "10A"},
		{NISN: "124", FullName: "Jane", ClassID: "c1", ClassName: "10A"},
	})
	assert.NoError(t, err)

	// Missing fields
	err = svc.CreateBulkFromRaw([]domain.StudentRawEntry{{FullName: "Fail", ClassID: "c1"}})
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestStudentService_GetByID(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	student := &domain.Student{ID: "s1"}
	mockRepo.On("GetByID", "s1").Return(student, nil)

	res, err := svc.GetByID("s1")
	assert.NoError(t, err)
	assert.Equal(t, "s1", res.ID)
}

func TestStudentService_GetByClassID(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	students := []*domain.Student{{ID: "s1"}}
	mockRepo.On("GetByClassID", "c1").Return(students, nil)

	res, err := svc.GetByClassID("c1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "s1", res[0].ID)
}

func TestStudentService_GetAll(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	students := []*domain.Student{{ID: "s1"}}
	mockRepo.On("GetAll", "teacher1", (*bool)(nil), "", "", 1, 10).Return(students, 1, nil)

	res, count, err := svc.GetAll("teacher1", nil, "", "", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "s1", res[0].ID)
	assert.Equal(t, 1, count)
}

func TestStudentService_Update(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	mockRepo.On("GetByID", "s1").Return(&domain.Student{ID: "s1"}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Student")).Return(nil)

	err := svc.Update("s1", "123", "John", "c1", "10A", nil, nil, true)
	assert.NoError(t, err)
}

func TestStudentService_SoftDelete(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	mockRepo.On("SoftDelete", "s1").Return(nil)

	err := svc.SoftDelete("s1")
	assert.NoError(t, err)
}
