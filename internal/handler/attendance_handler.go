package handler

import (
	"net/http"
	"time"

	"alhikmah-attendance-api/internal/domain"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service domain.AttendanceService
}

func NewAttendanceHandler(service domain.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) ScanQR(c *gin.Context) {
	var req struct {
		QRCodeData string `json:"qr_code_data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	if err := h.service.ScanQR(req.QRCodeData, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Attendance recorded successfully",
	})
}

func (h *AttendanceHandler) ManualInput(c *gin.Context) {
	var req struct {
		ClassID    string   `json:"class_id" binding:"required"`
		StudentIDs []string `json:"student_ids" binding:"required"`
		Status     string   `json:"status" binding:"required"`
		Notes      string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	if err := h.service.ManualInput(req.ClassID, req.StudentIDs, req.Status, req.Notes, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Manual attendance recorded successfully",
	})
}

func (h *AttendanceHandler) GetClassAttendanceForToday(c *gin.Context) {
	classID := c.Param("class_id")

	attendances, err := h.service.GetClassAttendanceForToday(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    attendances,
	})
}

func (h *AttendanceHandler) GetByClassAndDate(c *gin.Context) {
	classID := c.Param("class_id")
	dateStr := c.Param("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected YYYY-MM-DD"})
		return
	}

	attendances, err := h.service.GetByClassAndDate(classID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    attendances,
	})
}

func (h *AttendanceHandler) Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Attendance updated (stub)",
	})
}
