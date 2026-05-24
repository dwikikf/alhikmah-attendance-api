package handler

import (
	"net/http"
	"strconv"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service domain.ReportService
}

func NewReportHandler(service domain.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetDailyReport(c *gin.Context) {
	classID := c.Query("class_id")
	dateStr := c.Query("date")

	if classID == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, response.Error("class_id and date are required"))
		return
	}

	report, err := h.service.GetDailyReport(classID, dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Daily report fetched successfully", report))
}

func (h *ReportHandler) GetMonthlyReport(c *gin.Context) {
	classID := c.Query("class_id")
	monthStr := c.Query("month")

	if classID == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, response.Error("class_id and month are required"))
		return
	}

	report, err := h.service.GetMonthlyReport(classID, monthStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Monthly report fetched successfully", report))
}

func (h *ReportHandler) GetSemesterReport(c *gin.Context) {
	classID := c.Query("class_id")
	semesterStr := c.Query("semester")
	academicYear := c.Query("academic_year")

	if classID == "" || semesterStr == "" || academicYear == "" {
		c.JSON(http.StatusBadRequest, response.Error("class_id, semester, and academic_year are required"))
		return
	}

	semester, err := strconv.Atoi(semesterStr)
	if err != nil || (semester != 1 && semester != 2) {
		c.JSON(http.StatusBadRequest, response.Error("invalid semester value, must be 1 or 2"))
		return
	}

	report, err := h.service.GetSemesterReport(classID, academicYear, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Semester report fetched successfully", report))
}

func (h *ReportHandler) GetStudentReport(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("Student report fetched (stub)", nil))
}

func (h *ReportHandler) Export(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("Report exported (stub)", nil))
}
