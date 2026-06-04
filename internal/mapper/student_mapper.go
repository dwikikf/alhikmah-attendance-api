package mapper

import (
        "alhikmah-attendance-api/internal/domain"
        "alhikmah-attendance-api/internal/dto"
)

func ToStudentDTO(student *domain.Student) *dto.StudentResponse {
        if student == nil {
                return nil
        }
        return &dto.StudentResponse{
                ID:              student.ID,
                NISN:            student.NISN,
                FullName:        student.FullName,
                ClassID:         student.ClassID,
                ClassName:       student.ClassName,
                DOB:             student.DOB,
                Gender:          student.Gender,
                PhotoURL:        student.PhotoURL,
                QRCodeData:      student.QRCodeData,
                IsActive:        student.IsActive,
                CreatedAt:       student.CreatedAt,
                UpdatedAt:       student.UpdatedAt,
                AttendanceToday: student.AttendanceToday,
                ScannedAt:       student.ScannedAt,
        }
}
