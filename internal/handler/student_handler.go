package handler

import (
	"fmt"
	"net/http"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/dto"
	"alhikmah-attendance-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	service domain.StudentService
}

func NewStudentHandler(service domain.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req dto.CreateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	student := &domain.Student{
		NISN:     req.NISN,
		FullName: req.FullName,
		ClassID:  req.ClassID,
		Gender:   req.Gender,
		DOB:      req.DOB,
	}

	if err := h.service.Create(student); err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success("Student created successfully", student))
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	studentID := c.Param("student_id")

	student, err := h.service.GetByID(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Student not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Student fetched successfully", student))
}

func (h *StudentHandler) GetByClass(c *gin.Context) {
	classID := c.Param("class_id")

	students, err := h.service.GetByClassID(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to fetch students"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Students fetched successfully", students))
}

func (h *StudentHandler) GetAll(c *gin.Context) {
	isActiveStr := c.Query("is_active")
	classID := c.Query("class_id")
	search := c.Query("search")
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

	students, total, err := h.service.GetAll(isActive, classID, search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.SuccessWithPagination("Students fetched successfully", students, gin.H{
		"page":            page,
		"pageSize":        limit,
		"totalItems":      total,
		"totalPages":      totalPages,
		"hasNextPage":     page < totalPages,
		"hasPreviousPage": page > 1,
	}))
}

func (h *StudentHandler) Delete(c *gin.Context) {
	id := c.Param("student_id")

	if err := h.service.SoftDelete(id); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Student deleted successfully", nil))
}

func (h *StudentHandler) Update(c *gin.Context) {
	studentID := c.Param("student_id")
	var req dto.UpdateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	student := &domain.Student{
		ID:       studentID,
		FullName: req.FullName,
		ClassID:  req.ClassID,
		DOB:      req.DOB,
		Gender:   req.Gender,
		IsActive: req.IsActive,
	}

	if err := h.service.Update(student); err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("Student updated successfully", nil))
}

func (h *StudentHandler) GetQRCode(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("Student QR Code generated (stub)", nil))
}
