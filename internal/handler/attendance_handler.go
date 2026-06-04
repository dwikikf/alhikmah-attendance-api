package handler

import (
	"net/http"
	"time"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/dto"
	"alhikmah-attendance-api/internal/mapper"
	"alhikmah-attendance-api/pkg/response"
	"alhikmah-attendance-api/pkg/utils"

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
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	if err := h.service.ProcessQRScan(req.NISN, userID.(string), role.(string), req.Subject); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Attendance recorded successfully", nil))
}

func (h *AttendanceHandler) ManualInput(c *gin.Context) {
	var req dto.ManualAttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	if err := h.service.ProcessManualAttendance(req.StudentID, req.Status, req.Notes, userID.(string), role.(string), req.Subject); err != nil {
		HandleAppError(c, err, "")
		return
	}

	c.JSON(http.StatusOK, response.Success("Manual attendance recorded successfully", nil))
}

func (h *AttendanceHandler) GetClassAttendanceForToday(c *gin.Context) {
	classID := c.Param("class_id")
	subject := c.Query("subject")

	attendances, err := h.service.GetClassAttendanceForToday(classID, subject)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch attendance data")
		return
	}

	var dtoAttendances []*dto.AttendanceResponse
	for _, a := range attendances {
		dtoAttendances = append(dtoAttendances, mapper.ToAttendanceDTO(a))
	}

	c.JSON(http.StatusOK, response.Success("Attendance fetched successfully", dtoAttendances))
}

func (h *AttendanceHandler) GetByClassAndDate(c *gin.Context) {
	classID := c.Param("class_id")
	dateStr := c.Param("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid date format, expected YYYY-MM-DD"))
		return
	}

	subject := c.Query("subject")

	attendances, err := h.service.GetByClassAndDate(classID, date, subject)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch attendance data")
		return
	}

	var dtoAttendances []*dto.AttendanceResponse
	for _, a := range attendances {
		dtoAttendances = append(dtoAttendances, mapper.ToAttendanceDTO(a))
	}

	c.JSON(http.StatusOK, response.Success("Attendance fetched successfully", dtoAttendances))
}

func (h *AttendanceHandler) Update(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("Attendance updated (stub)", nil))
}
