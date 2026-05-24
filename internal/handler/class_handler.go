package handler

import (
	"net/http"
	"strconv"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/dto"
	"alhikmah-attendance-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	service domain.ClassService
}

func NewClassHandler(service domain.ClassService) *ClassHandler {
	return &ClassHandler{service: service}
}

func (h *ClassHandler) GetAll(c *gin.Context) {
	// A teacher should only see their own classes, admin sees all
	role, _ := c.Get("role")
	userID, _ := c.Get("userID")
	
	teacherID := ""
	if role == "guru" {
		teacherID = userID.(string)
	}

	academicYear := c.Query("academic_year")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	classes, total, err := h.service.GetAll(teacherID, academicYear, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to fetch classes"))
		return
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.SuccessWithPagination("Classes fetched successfully", classes, gin.H{
		"page":            page,
		"pageSize":        limit,
		"totalItems":      total,
		"totalPages":      totalPages,
		"hasNextPage":     page < totalPages,
		"hasPreviousPage": page > 1,
	}))
}

func (h *ClassHandler) GetByID(c *gin.Context) {
	id := c.Param("class_id")

	class, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to fetch class"))
		return
	}

	// For PRD, the endpoint also needs to return students if `include_attendance` is true.
	// Since we are building iteratively, we will just return the class for now.
	// We'll update this in Phase 3.

	c.JSON(http.StatusOK, response.Success("Class fetched successfully", class))
}

func (h *ClassHandler) Create(c *gin.Context) {
	var req dto.CreateClassRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	class := &domain.Class{
		ClassName:    req.ClassName,
		TeacherID:    req.TeacherID,
		AcademicYear: req.AcademicYear,
		Capacity:     req.Capacity,
		Description:  req.Description,
	}

	if err := h.service.Create(class); err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success("Class created successfully", class))
}

func (h *ClassHandler) Update(c *gin.Context) {
	id := c.Param("class_id")

	var req dto.UpdateClassRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	class := &domain.Class{
		ID:          id,
		ClassName:   req.ClassName,
		TeacherID:   req.TeacherID,
		Capacity:    req.Capacity,
		Description: req.Description,
	}

	if err := h.service.Update(class); err != nil {
		if err.Error() == "class not found" {
			c.JSON(http.StatusNotFound, response.Error("Class not found"))
			return
		}
		handleDBError(c, err)
		return
	}

	// Fetch updated class
	updatedClass, err := h.service.GetByID(id)
	if err != nil {
		updatedClass = class
	}

	c.JSON(http.StatusOK, response.Success("Class updated successfully", updatedClass))
}

func (h *ClassHandler) Delete(c *gin.Context) {
	id := c.Param("class_id")

	if err := h.service.SoftDelete(id); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Class deleted successfully", nil))
}
