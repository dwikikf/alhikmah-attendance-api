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

type ClassTeacherHandler struct {
	service domain.ClassTeacherService
}

func NewClassTeacherHandler(service domain.ClassTeacherService) *ClassTeacherHandler {
	return &ClassTeacherHandler{service: service}
}

func (h *ClassTeacherHandler) Assign(c *gin.Context) {
	classID := c.Param("class_id")

	var req dto.AssignTeacherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	if err := h.service.Assign(req.TeacherID, classID, req.Subject); err != nil {
		HandleAppError(c, err, "")
		return
	}

	c.JSON(http.StatusCreated, response.Success("Teacher assigned successfully", nil))
}

func (h *ClassTeacherHandler) Unassign(c *gin.Context) {
	classID := c.Param("class_id")
	teacherID := c.Param("teacher_id")
	subject := c.Query("subject")

	if err := h.service.Unassign(teacherID, classID, subject); err != nil {
		HandleAppError(c, err, "")
		return
	}

	c.JSON(http.StatusOK, response.Success("Teacher unassigned successfully", nil))
}

func (h *ClassTeacherHandler) GetByClassID(c *gin.Context) {
	classID := c.Param("class_id")

	list, err := h.service.GetByClassID(classID)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch teachers data")
		return
	}

	var dtoList []*dto.ClassTeacherResponse
	for _, t := range list {
		dtoList = append(dtoList, mapper.ToClassTeacherDTO(t))
	}

	c.JSON(http.StatusOK, response.Success("Teachers fetched successfully", dtoList))
}

func (h *ClassTeacherHandler) GetSubjectAssignments(c *gin.Context) {
	teacherID := c.Param("teacher_id")

	list, err := h.service.GetSubjectAssignments(teacherID)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch assignments data")
		return
	}

	var dtoList []*dto.ClassTeacherResponse
	for _, t := range list {
		dtoList = append(dtoList, mapper.ToClassTeacherDTO(t))
	}

	c.JSON(http.StatusOK, response.Success("Assignments fetched successfully", dtoList))
}
