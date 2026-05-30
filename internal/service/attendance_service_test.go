package service_test

import (
	"errors"
	"testing"
	"time"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/mocks"
	"alhikmah-attendance-api/internal/service"
	"alhikmah-attendance-api/pkg/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessQRScan_DuplicateCachePrevention(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()

	// Pre-populate cache
	nisn := "12345"
	c.Set("qr_scan_"+nisn, true, 1*time.Minute)

	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	err := svc.ProcessQRScan(nisn, "teacher1", "teacher")

	assert.Error(t, err)
	assert.Equal(t, "Siswa baru saja melakukan scan QR", err.Error())
}

func TestProcessQRScan_StudentNotFound(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()

	nisn := "12345"
	mockStudentRepo.On("GetByNISN", nisn).Return(nil, errors.New("not found"))

	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	err := svc.ProcessQRScan(nisn, "teacher1", "teacher")

	assert.Error(t, err)
	assert.Equal(t, "Siswa tidak ditemukan", err.Error())
}

func TestProcessQRScan_UnauthorizedTeacher(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()

	nisn := "12345"
	teacherID := "teacher1"
	student := &domain.Student{ID: "student1", ClassID: "class1"}
	
	mockStudentRepo.On("GetByNISN", nisn).Return(student, nil)
	mockClassRepo.On("IsTeacherResponsibleForStudent", student.ID, teacherID).Return(false, nil)

	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	err := svc.ProcessQRScan(nisn, teacherID, "teacher")

	assert.Error(t, err)
	assert.Equal(t, "Guru tidak memiliki akses untuk kelas siswa ini", err.Error())
}

func TestProcessQRScan_AlreadyScannedToday(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()

	nisn := "12345"
	teacherID := "teacher1"
	student := &domain.Student{ID: "student1", ClassID: "class1"}
	today := time.Now().Truncate(24 * time.Hour)
	attendances := []*domain.Attendance{
		{StudentID: student.ID, ClassID: student.ClassID},
	}
	
	mockStudentRepo.On("GetByNISN", nisn).Return(student, nil)
	mockClassRepo.On("IsTeacherResponsibleForStudent", student.ID, teacherID).Return(true, nil)
	mockRepo.On("GetByClassAndDate", student.ClassID, today).Return(attendances, nil)

	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	err := svc.ProcessQRScan(nisn, teacherID, "teacher")

	assert.Error(t, err)
	assert.Equal(t, "Siswa sudah melakukan absensi hari ini", err.Error())
}

func TestProcessQRScan_Success(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()

	nisn := "12345"
	teacherID := "teacher1"
	student := &domain.Student{ID: "student1", ClassID: "class1"}
	today := time.Now().Truncate(24 * time.Hour)
	
	mockStudentRepo.On("GetByNISN", nisn).Return(student, nil)
	mockClassRepo.On("IsTeacherResponsibleForStudent", student.ID, teacherID).Return(true, nil)
	mockRepo.On("GetByClassAndDate", student.ClassID, today).Return([]*domain.Attendance{}, nil)
	mockRepo.On("MarkAttendance", mock.AnythingOfType("*domain.Attendance")).Return(nil)

	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	err := svc.ProcessQRScan(nisn, teacherID, "teacher")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProcessManualAttendance_InvalidStatus(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()
	
	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	err := svc.ProcessManualAttendance("student1", "invalid_status", nil, "teacher1", "teacher")
	
	assert.Error(t, err)
	assert.Equal(t, "Status absensi tidak valid", err.Error())
}

func TestProcessManualAttendance_Success(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()
	
	studentID := "student1"
	teacherID := "teacher1"
	student := &domain.Student{ID: studentID, ClassID: "class1"}
	today := time.Now().Truncate(24 * time.Hour)
	
	mockStudentRepo.On("GetByID", studentID).Return(student, nil)
	mockClassRepo.On("IsTeacherResponsibleForStudent", studentID, teacherID).Return(true, nil)
	mockRepo.On("GetByClassAndDate", student.ClassID, today).Return([]*domain.Attendance{}, nil)
	mockRepo.On("MarkAttendance", mock.AnythingOfType("*domain.Attendance")).Return(nil)
	
	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	err := svc.ProcessManualAttendance(studentID, "izin", nil, teacherID, "teacher")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateManual_Success(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()

	attendanceID := "att-1"
	status := "izin"
	reason := "sakit"
	changedBy := "admin-1"
	oldStatus := "hadir"

	mockRepo.On("GetByID", attendanceID).Return(&domain.Attendance{ID: attendanceID, Status: oldStatus}, nil)
	mockRepo.On("UpdateAttendance", mock.AnythingOfType("*domain.Attendance"), mock.AnythingOfType("*domain.AttendanceAudit")).Return(nil)

	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	err := svc.UpdateManual(attendanceID, status, reason, changedBy)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateManual_NotFound(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()

	attendanceID := "att-1"
	mockRepo.On("GetByID", attendanceID).Return((*domain.Attendance)(nil), errors.New("not found"))

	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	err := svc.UpdateManual(attendanceID, "izin", "", "admin-1")

	assert.Error(t, err)
	assert.Equal(t, "Data absensi tidak ditemukan", err.Error())
}

func TestGetClassAttendanceForToday(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()

	classID := "class1"
	today := time.Now().Truncate(24 * time.Hour)
	attendances := []*domain.Attendance{{ID: "att1"}}

	mockRepo.On("GetByClassAndDate", classID, today).Return(attendances, nil)

	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	res, err := svc.GetClassAttendanceForToday(classID)

	assert.NoError(t, err)
	assert.Equal(t, attendances, res)
	mockRepo.AssertExpectations(t)
}

func TestGetByClassAndDate(t *testing.T) {
	mockRepo := new(mocks.AttendanceRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockClassRepo := new(mocks.ClassRepository)
	c := cache.NewCache()

	classID := "class1"
	date := time.Now()
	attendances := []*domain.Attendance{{ID: "att1"}}

	mockRepo.On("GetByClassAndDate", classID, date).Return(attendances, nil)

	svc := service.NewAttendanceService(mockRepo, mockStudentRepo, mockClassRepo, c)
	res, err := svc.GetByClassAndDate(classID, date)

	assert.NoError(t, err)
	assert.Equal(t, attendances, res)
	mockRepo.AssertExpectations(t)
}
