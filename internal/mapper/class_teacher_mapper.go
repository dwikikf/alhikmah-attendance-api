package mapper

import (
        "alhikmah-attendance-api/internal/domain"
        "alhikmah-attendance-api/internal/dto"
)

func ToClassTeacherDTO(ct *domain.ClassTeacher) *dto.ClassTeacherResponse {
        if ct == nil {
                return nil
        }
        return &dto.ClassTeacherResponse{
                ID:           ct.ID,
                TeacherID:    ct.TeacherID,
                TeacherName:  ct.TeacherName,
                ClassID:      ct.ClassID,
                ClassDisplay: ct.ClassDisplay,
                AcademicYear: ct.AcademicYear,
                Subject:      ct.Subject,
                Role:         ct.Role,
                CreatedAt:    ct.CreatedAt,
        }
}
