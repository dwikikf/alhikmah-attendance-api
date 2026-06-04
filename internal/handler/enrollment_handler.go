package handler

import (
	"net/http"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/dto"
	"alhikmah-attendance-api/internal/mapper"
	"alhikmah-attendance-api/pkg/response"
	"alhikmah-attendance-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type EnrollmentHandler struct {
	service domain.EnrollmentService
}

func NewEnrollmentHandler(service domain.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{service: service}
}

func (h *EnrollmentHandler) Enroll(c *gin.Context) {
	classID := c.Param("class_id")

	var req dto.EnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	if err := h.service.Enroll(req.StudentID, classID, req.AcademicYear); err != nil {
		HandleAppError(c, err, "")
		return
	}

	c.JSON(http.StatusCreated, response.Success("Student enrolled successfully", nil))
}

func (h *EnrollmentHandler) PromoteClass(c *gin.Context) {
	var req dto.PromoteClassRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	var items []domain.PromoteItem
	for _, item := range req.Items {
		items = append(items, domain.PromoteItem{
			StudentID:     item.StudentID,
			TargetClassID: item.TargetClassID,
		})
	}

	count, err := h.service.PromoteClass(items, req.AcademicYear)
	if err != nil {
		HandleAppError(c, err, "")
		return
	}

	c.JSON(http.StatusOK, response.Success("Students promoted successfully", gin.H{"promoted_count": count}))
}

func (h *EnrollmentHandler) TransferStudent(c *gin.Context) {
	studentID := c.Param("student_id")

	var req dto.TransferStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	if err := h.service.TransferStudent(studentID, req.FromClassID, req.ToClassID, req.AcademicYear); err != nil {
		HandleAppError(c, err, "")
		return
	}

	c.JSON(http.StatusOK, response.Success("Student transferred successfully", nil))
}

func (h *EnrollmentHandler) GetHistoryByStudentID(c *gin.Context) {
	studentID := c.Param("student_id")

	history, err := h.service.GetHistoryByStudentID(studentID)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch student history")
		return
	}

	var dtoHistory []*dto.EnrollmentHistoryResponse
	for _, e := range history {
		dtoHistory = append(dtoHistory, mapper.ToEnrollmentHistoryDTO(e))
	}

	c.JSON(http.StatusOK, response.Success("Enrollment history fetched successfully", dtoHistory))
}

func (h *EnrollmentHandler) GetActiveByClassID(c *gin.Context) {
	classID := c.Param("class_id")

	enrollments, err := h.service.GetActiveByClassID(classID)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch active data")
		return
	}

	var dtoEnrollments []*dto.EnrollmentResponse
	for _, e := range enrollments {
		dtoEnrollments = append(dtoEnrollments, mapper.ToEnrollmentDTO(e))
	}

	c.JSON(http.StatusOK, response.Success("Active enrollments fetched successfully", dtoEnrollments))
}
