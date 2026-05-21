package handler

import (
	"net/http"

	"alhikmah-attendance-api/internal/domain"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	service domain.StudentService
}

func NewStudentHandler(service domain.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req struct {
		NISN     string `json:"nisn" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		ClassID  string `json:"class_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student := &domain.Student{
		NISN:     req.NISN,
		FullName: req.FullName,
		ClassID:  req.ClassID,
	}

	if err := h.service.Create(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    student,
		"message": "Student created successfully",
	})
}

func (h *StudentHandler) GetByClass(c *gin.Context) {
	classID := c.Param("class_id")

	students, err := h.service.GetByClassID(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    students,
	})
}
