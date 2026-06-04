package service_test

import (
	"testing"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/mocks"
	"alhikmah-attendance-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEnrollmentService_Enroll(t *testing.T) {
	mockRepo := new(mocks.EnrollmentRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	svc := service.NewEnrollmentService(mockRepo, mockStudentRepo, mockClassRepo)

	// Success
	mockStudentRepo.On("GetByID", "s1").Return(&domain.Student{ID: "s1"}, nil).Once()
	mockClassRepo.On("GetByID", "c1").Return(&domain.Class{ID: "c1"}, nil).Once()
	mockRepo.On("GetActiveByStudentID", "s1").Return((*domain.StudentEnrollment)(nil), nil).Once()
	mockRepo.On("Enroll", mock.AnythingOfType("*domain.StudentEnrollment")).Return(nil).Once()

	err := svc.Enroll("s1", "c1", "2025/2026")
	assert.NoError(t, err)

	// Fail: already active in same TA
	mockStudentRepo.On("GetByID", "s2").Return(&domain.Student{ID: "s2"}, nil).Once()
	mockClassRepo.On("GetByID", "c2").Return(&domain.Class{ID: "c2"}, nil).Once()
	mockRepo.On("GetActiveByStudentID", "s2").Return(&domain.StudentEnrollment{ClassDisplay: "Kelas X", AcademicYear: "2025/2026"}, nil).Once()
	
	err = svc.Enroll("s2", "c2", "2025/2026")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already has active enrollment")

	mockRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
	mockClassRepo.AssertExpectations(t)
}

func TestEnrollmentService_PromoteClass(t *testing.T) {
	mockRepo := new(mocks.EnrollmentRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	svc := service.NewEnrollmentService(mockRepo, mockStudentRepo, mockClassRepo)

	items := []domain.PromoteItem{
		{StudentID: "s1", TargetClassID: "c2"},
		{StudentID: "s2", TargetClassID: "c2"},
	}

	domainItems := []domain.PromoteItem{
		{StudentID: "s1", TargetClassID: "c2"},
		{StudentID: "s2", TargetClassID: "c2"},
	}

	mockRepo.On("BulkEnroll", domainItems, "2026/2027").Return(2, nil).Once()

	count, err := svc.PromoteClass(items, "2026/2027")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	mockRepo.AssertExpectations(t)
}

func TestEnrollmentService_TransferStudent(t *testing.T) {
	mockRepo := new(mocks.EnrollmentRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	svc := service.NewEnrollmentService(mockRepo, mockStudentRepo, mockClassRepo)

	// Success
	mockRepo.On("GetActiveByStudentID", "s1").Return(&domain.StudentEnrollment{ClassID: "c1"}, nil).Once()
	mockRepo.On("EndEnrollment", "s1", "c1", "transferred").Return(nil).Once()
	
	// Enroll calls GetActiveByStudentID again
	mockStudentRepo.On("GetByID", "s1").Return(&domain.Student{ID: "s1"}, nil).Once()
	mockClassRepo.On("GetByID", "c2").Return(&domain.Class{ID: "c2"}, nil).Once()
	mockRepo.On("GetActiveByStudentID", "s1").Return((*domain.StudentEnrollment)(nil), nil).Once()
	mockRepo.On("Enroll", mock.AnythingOfType("*domain.StudentEnrollment")).Return(nil).Once()

	err := svc.TransferStudent("s1", "c1", "c2", "2025/2026")
	assert.NoError(t, err)

	// Fail: not in fromClassID
	mockRepo.On("GetActiveByStudentID", "s2").Return(&domain.StudentEnrollment{ClassID: "c99"}, nil).Once()
	err = svc.TransferStudent("s2", "c1", "c2", "2025/2026")
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
