package mapper

import (
        "alhikmah-attendance-api/internal/domain"
        "alhikmah-attendance-api/internal/dto"
)

func ToEnrollmentDTO(e *domain.StudentEnrollment) *dto.EnrollmentResponse {
        if e == nil {
                return nil
        }
        return &dto.EnrollmentResponse{
                ID:           e.ID,
                StudentID:    e.StudentID,
                StudentName:  e.StudentName,
                ClassID:      e.ClassID,
                ClassDisplay: e.ClassDisplay,
                AcademicYear: e.AcademicYear,
                EnrolledAt:   e.EnrolledAt,
        }
}

func ToEnrollmentHistoryDTO(e *domain.StudentEnrollment) *dto.EnrollmentHistoryResponse {
	if e == nil {
		return nil
	}
	return &dto.EnrollmentHistoryResponse{
		ID:           e.ID,
		ClassID:      e.ClassID,
		ClassDisplay: e.ClassDisplay,
		AcademicYear: e.AcademicYear,
		EnrolledAt:   e.EnrolledAt,
		EndedAt:      e.EndedAt,
		EndReason:    e.EndReason,
	}
}
