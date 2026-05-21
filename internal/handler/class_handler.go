package handler

import (
	"net/http"
	"strconv"

	"alhikmah-attendance-api/internal/domain"

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

	classes, total, err := h.service.GetAll(teacherID, academicYear, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch classes"})
		return
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 20
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    classes,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *ClassHandler) GetByID(c *gin.Context) {
	id := c.Param("class_id")

	class, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch class"})
		return
	}

	// For PRD, the endpoint also needs to return students if `include_attendance` is true.
	// Since we are building iteratively, we will just return the class for now.
	// We'll update this in Phase 3.

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    class,
	})
}
