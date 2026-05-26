package handler

import (
	"net/http"
	"time"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/dto"
	"alhikmah-attendance-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service domain.AttendanceService
}

func NewAttendanceHandler(service domain.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) ScanQR(c *gin.Context) {
	var req dto.ScanQRRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	if err := h.service.ProcessQRScan(req.NISN, userID.(string), role.(string)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Attendance recorded successfully", nil))
}

func (h *AttendanceHandler) ManualInput(c *gin.Context) {
	var req dto.ManualAttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	if err := h.service.ProcessManualAttendance(req.StudentID, req.Status, req.Notes, userID.(string), role.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Manual attendance recorded successfully", nil))
}

func (h *AttendanceHandler) GetClassAttendanceForToday(c *gin.Context) {
	classID := c.Param("class_id")

	attendances, err := h.service.GetClassAttendanceForToday(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to fetch attendance"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Attendance fetched successfully", attendances))
}

func (h *AttendanceHandler) GetByClassAndDate(c *gin.Context) {
	classID := c.Param("class_id")
	dateStr := c.Param("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid date format, expected YYYY-MM-DD"))
		return
	}

	attendances, err := h.service.GetByClassAndDate(classID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to fetch attendance"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Attendance fetched successfully", attendances))
}

func (h *AttendanceHandler) Update(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("Attendance updated (stub)", nil))
}
