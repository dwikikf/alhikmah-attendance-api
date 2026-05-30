package service_test

import (
	"testing"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/mocks"
	"alhikmah-attendance-api/internal/service"

	"github.com/stretchr/testify/assert"
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

	student := &domain.Student{NISN: "123", FullName: "John", ClassID: "c1", ClassName: "10A"}
	mockRepo.On("Create", student).Return(nil)

	err := svc.Create(student)
	assert.NoError(t, err)
	assert.Equal(t, "123|John|10A", student.QRCodeData)

	// Missing fields
	err = svc.Create(&domain.Student{FullName: "John", ClassID: "c1"})
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestStudentService_CreateBulk(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	students := []*domain.Student{
		{NISN: "123", FullName: "John", ClassID: "c1", ClassName: "10A"},
		{NISN: "124", FullName: "Jane", ClassID: "c1", ClassName: "10A"},
	}

	mockRepo.On("CreateBulk", students).Return(nil)

	err := svc.CreateBulk(students)
	assert.NoError(t, err)
	assert.Equal(t, "123|John|10A", students[0].QRCodeData)

	// Missing fields
	err = svc.CreateBulk([]*domain.Student{{FullName: "Fail", ClassID: "c1"}})
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
	assert.Equal(t, student, res)
}

func TestStudentService_GetByClassID(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	students := []*domain.Student{{ID: "s1"}}
	mockRepo.On("GetByClassID", "c1").Return(students, nil)

	res, err := svc.GetByClassID("c1")
	assert.NoError(t, err)
	assert.Equal(t, students, res)
}

func TestStudentService_GetAll(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	students := []*domain.Student{{ID: "s1"}}
	mockRepo.On("GetAll", "teacher1", (*bool)(nil), "", "", 1, 10).Return(students, 1, nil)

	res, count, err := svc.GetAll("teacher1", nil, "", "", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, students, res)
	assert.Equal(t, 1, count)
}

func TestStudentService_Update(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	student := &domain.Student{ID: "s1", NISN: "123", FullName: "John", ClassName: "10A"}
	mockRepo.On("Update", student).Return(nil)

	err := svc.Update(student)
	assert.NoError(t, err)
	assert.Equal(t, "123|John|10A", student.QRCodeData)
}

func TestStudentService_SoftDelete(t *testing.T) {
	mockRepo := new(mocks.StudentRepository)
	svc := service.NewStudentService(mockRepo)

	mockRepo.On("SoftDelete", "s1").Return(nil)

	err := svc.SoftDelete("s1")
	assert.NoError(t, err)
}
