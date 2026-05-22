package handler

import (
	"fmt"
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

func (h *StudentHandler) GetByID(c *gin.Context) {
	studentID := c.Param("student_id")

	student, err := h.service.GetByID(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    student,
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

func (h *StudentHandler) GetAll(c *gin.Context) {
	isActiveStr := c.Query("is_active")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var isActive *bool
	if isActiveStr == "true" {
		t := true
		isActive = &t
	} else if isActiveStr == "false" {
		f := false
		isActive = &f
	}

	page := 1
	limit := 10
	if pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	}
	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	students, total, err := h.service.GetAll(isActive, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    students,
		"pagination": gin.H{
			"page":            page,
			"pageSize":        limit,
			"totalItems":      total,
			"totalPages":      totalPages,
			"hasNextPage":     page < totalPages,
			"hasPreviousPage": page > 1,
		},
	})
}

func (h *StudentHandler) Delete(c *gin.Context) {
	id := c.Param("student_id")

	if err := h.service.SoftDelete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Student deleted successfully",
	})
}

func (h *StudentHandler) Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Student updated (stub)",
	})
}

func (h *StudentHandler) GetQRCode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Student QR Code generated (stub)",
	})
}
