package handler

import (
	"net/http"
	"strconv"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/mapper"
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
	subject := c.Query("subject")
	forceRefresh := c.Query("force_refresh") == "true"

	if classID == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, response.Error("class_id and date are required"))
		return
	}

	userID, _ := c.Get("userID")
	generatedBy := ""
	if uid, ok := userID.(string); ok {
		generatedBy = uid
	}

	report, err := h.service.GetDailyReport(classID, dateStr, subject, forceRefresh, generatedBy)
	if err != nil {
		HandleAppError(c, err, "")
		return
	}

	dtoReport := mapper.ToDailyReportDTO(report)
	c.JSON(http.StatusOK, response.Success("Daily report fetched successfully", dtoReport))
}

func (h *ReportHandler) GetMonthlyReport(c *gin.Context) {
	classID := c.Query("class_id")
	monthStr := c.Query("month")
	subject := c.Query("subject")
	forceRefresh := c.Query("force_refresh") == "true"

	if classID == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, response.Error("class_id and month are required"))
		return
	}

	userID, _ := c.Get("userID")
	generatedBy := ""
	if uid, ok := userID.(string); ok {
		generatedBy = uid
	}

	report, err := h.service.GetMonthlyReport(classID, monthStr, subject, forceRefresh, generatedBy)
	if err != nil {
		HandleAppError(c, err, "")
		return
	}

	dtoReport := mapper.ToMonthlyReportDTO(report)
	c.JSON(http.StatusOK, response.Success("Monthly report fetched successfully", dtoReport))
}

func (h *ReportHandler) GetSemesterReport(c *gin.Context) {
	classID := c.Query("class_id")
	semesterStr := c.Query("semester")
	academicYear := c.Query("academic_year")
	subject := c.Query("subject")
	forceRefresh := c.Query("force_refresh") == "true"

	if classID == "" || semesterStr == "" || academicYear == "" {
		c.JSON(http.StatusBadRequest, response.Error("class_id, semester, and academic_year are required"))
		return
	}

	semester, err := strconv.Atoi(semesterStr)
	if err != nil || (semester != 1 && semester != 2) {
		c.JSON(http.StatusBadRequest, response.Error("invalid semester value, must be 1 or 2"))
		return
	}

	userID, _ := c.Get("userID")
	generatedBy := ""
	if uid, ok := userID.(string); ok {
		generatedBy = uid
	}

	report, err := h.service.GetSemesterReport(classID, academicYear, semester, subject, forceRefresh, generatedBy)
	if err != nil {
		HandleAppError(c, err, "")
		return
	}

	dtoReport := mapper.ToSemesterReportDTO(report)
	c.JSON(http.StatusOK, response.Success("Semester report fetched successfully", dtoReport))
}

func (h *ReportHandler) GetStudentReport(c *gin.Context) {
	studentID := c.Param("student_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if studentID == "" || fromDate == "" || toDate == "" {
		c.JSON(http.StatusBadRequest, response.Error("student_id, from_date, and to_date are required"))
		return
	}

	report, err := h.service.GetStudentReport(studentID, fromDate, toDate)
	if err != nil {
		HandleAppError(c, err, "")
		return
	}

	dtoReport := mapper.ToStudentReportDTO(report)
	c.JSON(http.StatusOK, response.Success("Student report fetched successfully", dtoReport))
}

type ExportRequest struct {
	ReportType   string `json:"report_type" binding:"required"`
	ClassID      string `json:"class_id" binding:"required"`
	Date         string `json:"date"`
	Month        string `json:"month"`
	Semester     string `json:"semester"`
	Subject      string `json:"subject"`
	AcademicYear string `json:"academic_year"`
	Format       string `json:"format" binding:"required"`
}

func (h *ReportHandler) Export(c *gin.Context) {
	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request payload"))
		return
	}

	if req.Format != "csv" && req.Format != "excel" {
		c.JSON(http.StatusBadRequest, response.Error("Unsupported format. Use 'csv' or 'excel'"))
		return
	}

	switch req.ReportType {
	case "harian":
		if req.Date == "" {
			c.JSON(http.StatusBadRequest, response.Error("date is required for daily report"))
			return
		}
		userID, _ := c.Get("userID")
		generatedBy := ""
		if uid, ok := userID.(string); ok {
			generatedBy = uid
		}
		report, err := h.service.GetDailyReport(req.ClassID, req.Date, req.Subject, false, generatedBy)
		if err != nil {
			HandleAppError(c, err, "")
			return
		}
		dtoReport := mapper.ToDailyReportDTO(report)
		h.exportDaily(c, dtoReport, req.Format)

	case "bulanan":
		if req.Month == "" {
			c.JSON(http.StatusBadRequest, response.Error("month is required for monthly report"))
			return
		}
		userID, _ := c.Get("userID")
		generatedBy := ""
		if uid, ok := userID.(string); ok {
			generatedBy = uid
		}
		report, err := h.service.GetMonthlyReport(req.ClassID, req.Month, req.Subject, false, generatedBy)
		if err != nil {
			HandleAppError(c, err, "")
			return
		}
		dtoReport := mapper.ToMonthlyReportDTO(report)
		h.exportMonthly(c, dtoReport, req.Month, req.Format)

	case "semesteran":
		if req.Semester == "" || req.AcademicYear == "" {
			c.JSON(http.StatusBadRequest, response.Error("semester and academic_year are required"))
			return
		}
		sem, err := strconv.Atoi(req.Semester)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error("invalid semester format"))
			return
		}
		userID, _ := c.Get("userID")
		generatedBy := ""
		if uid, ok := userID.(string); ok {
			generatedBy = uid
		}
		report, err := h.service.GetSemesterReport(req.ClassID, req.AcademicYear, sem, req.Subject, false, generatedBy)
		if err != nil {
			HandleAppError(c, err, "")
			return
		}
		dtoReport := mapper.ToSemesterReportDTO(report)
		h.exportSemester(c, dtoReport, req.Format)

	default:
		c.JSON(http.StatusBadRequest, response.Error("Unsupported report_type"))
	}
}
