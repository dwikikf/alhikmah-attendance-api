package handler

import (
	"archive/zip"
	"fmt"
	"net/http"
	"strconv"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/dto"
	"alhikmah-attendance-api/internal/mapper"
	"alhikmah-attendance-api/pkg/response"
	"alhikmah-attendance-api/pkg/utils"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

type ClassHandler struct {
	service        domain.ClassService
	studentService domain.StudentService
}

func NewClassHandler(service domain.ClassService, studentService domain.StudentService) *ClassHandler {
	return &ClassHandler{
		service:        service,
		studentService: studentService,
	}
}

func (h *ClassHandler) GetAll(c *gin.Context) {
	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	teacherID := ""
	if role == "teacher" {
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
		HandleAppError(c, err, "Failed to fetch classes")
		return
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	var dtoClasses []*dto.ClassResponse
	for _, class := range classes {
		dtoClasses = append(dtoClasses, mapper.ToClassDTO(class))
	}

	c.JSON(http.StatusOK, response.SuccessWithPagination("Classes fetched successfully", dtoClasses, gin.H{
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
		HandleAppError(c, err, "Failed to fetch class")
		return
	}

	c.JSON(http.StatusOK, response.Success("Class fetched successfully", mapper.ToClassDTO(class)))
}

func (h *ClassHandler) Create(c *gin.Context) {
	var req dto.CreateClassRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	role, _ := c.Get("role")
	if role == "teacher" {
		userID, _ := c.Get("userID")
		req.TeacherID = userID.(string)
	}

	if err := h.service.Create(req.RoomName, req.Grade, req.Section, req.TeacherID, req.AcademicYear, req.Capacity, req.Description); err != nil {
		HandleAppError(c, err, "Failed to create class")
		return
	}

	c.JSON(http.StatusCreated, response.Success("Class created successfully", nil))
}

func (h *ClassHandler) Update(c *gin.Context) {
	id := c.Param("class_id")

	role, _ := c.Get("role")
	if role == "teacher" {
		userID, _ := c.Get("userID")
		class, err := h.service.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, response.Error("Class not found"))
			return
		}
		if class.TeacherID != userID.(string) {
			c.JSON(http.StatusForbidden, response.Error("You are not authorized to modify this class"))
			return
		}
	}

	var req dto.UpdateClassRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	if err := h.service.Update(id, req.RoomName, req.Grade, req.Section, req.TeacherID, req.Capacity, req.Description); err != nil {
		HandleAppError(c, err, "Failed to update class")
		return
	}

	// Fetch updated class
	updatedClass, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusOK, response.Success("Class updated successfully", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("Class updated successfully", mapper.ToClassDTO(updatedClass)))
}

func (h *ClassHandler) Delete(c *gin.Context) {
	id := c.Param("class_id")

	role, _ := c.Get("role")
	if role == "teacher" {
		userID, _ := c.Get("userID")
		class, err := h.service.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, response.Error("Class not found"))
			return
		}
		if class.TeacherID != userID.(string) {
			c.JSON(http.StatusForbidden, response.Error("You are not authorized to delete this class"))
			return
		}
	}

	if err := h.service.SoftDelete(id); err != nil {
		HandleAppError(c, err, "Failed to delete class")
		return
	}

	c.JSON(http.StatusOK, response.Success("Class deleted successfully", nil))
}

func (h *ClassHandler) ExportExcel(c *gin.Context) {
	classID := c.Param("class_id")

	students, err := h.studentService.GetByClassID(classID)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch students")
		return
	}

	class, err := h.service.GetByID(classID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Class not found"))
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Data Siswa"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"No", "NISN", "Nama Lengkap"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, student := range students {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), student.NISN)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), student.FullName)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=Data_Siswa_%s.xlsx", class.DisplayName))

	if err := f.Write(c.Writer); err != nil {
		return
	}
}

func (h *ClassHandler) ExportQRCode(c *gin.Context) {
	classID := c.Param("class_id")

	students, err := h.studentService.GetByClassID(classID)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch students")
		return
	}

	class, err := h.service.GetByID(classID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Class not found"))
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=QRCode_%s.zip", class.DisplayName))

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	for _, student := range students {
		qrData := student.QRCodeData
		if qrData == "" {
			qrData = h.studentService.GenerateQRCodeData(student.NISN, student.FullName, student.ClassID)
		}

		png, err := qrcode.Encode(qrData, qrcode.Medium, 256)
		if err != nil {
			continue
		}

		fileName := fmt.Sprintf("%s_%s.png", student.NISN, student.FullName)
		f, err := zipWriter.Create(fileName)
		if err != nil {
			continue
		}
		f.Write(png)
	}
}
