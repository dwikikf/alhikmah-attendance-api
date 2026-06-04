package mapper

import (
        "alhikmah-attendance-api/internal/domain"
        "alhikmah-attendance-api/internal/dto"
)

func ToClassDTO(class *domain.Class) *dto.ClassResponse {
        if class == nil {
                return nil
        }
        return &dto.ClassResponse{
                ID:           class.ID,
                RoomName:     class.RoomName,
                Grade:        class.Grade,
                Section:      class.Section,
                DisplayName:  class.GetDisplayName(),
                TeacherID:    class.TeacherID,
                TeacherName:  class.TeacherName,
                AcademicYear: class.AcademicYear,
                Capacity:     class.Capacity,
                StudentCount: class.StudentCount,
                Description:  class.Description,
                CreatedAt:    class.CreatedAt,
                UpdatedAt:    class.UpdatedAt,
        }
}
