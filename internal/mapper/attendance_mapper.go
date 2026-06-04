package mapper

import (
        "alhikmah-attendance-api/internal/domain"
        "alhikmah-attendance-api/internal/dto"
)

func ToAttendanceDTO(a *domain.Attendance) *dto.AttendanceResponse {
        if a == nil {
                return nil
        }
        return &dto.AttendanceResponse{
                ID:             a.ID,
                StudentID:      a.StudentID,
                ClassID:        a.ClassID,
                AttendanceDate: a.AttendanceDate,
                Subject:        a.Subject,
                Status:         a.Status,
                RecordedBy:     a.RecordedBy,
                RecordedAt:     a.RecordedAt,
                ScannedAt:      a.ScannedAt,
                Notes:          a.Notes,
                IsManual:       a.IsManual,
        }
}
