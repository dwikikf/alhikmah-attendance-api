package handler

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/dto"
	"alhikmah-attendance-api/internal/mapper"
	"alhikmah-attendance-api/pkg/response"
	"alhikmah-attendance-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	service      domain.StudentService
	classService domain.ClassService
}

func NewStudentHandler(service domain.StudentService, classService domain.ClassService) *StudentHandler {
	return &StudentHandler{service: service, classService: classService}
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req dto.CreateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	class, err := h.classService.GetByID(req.ClassID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid class ID"))
		return
	}

	role, _ := c.Get("role")
	if role == "teacher" {
		userID, _ := c.Get("userID")
		if class.TeacherID != userID.(string) {
			c.JSON(http.StatusForbidden, response.Error("You are not authorized to add students to this class"))
			return
		}
	}

	if err := h.service.Create(req.NISN, req.FullName, req.ClassID, req.DOB, req.Gender); err != nil {
		HandleAppError(c, err, "Failed to add student")
		return
	}

	c.JSON(http.StatusCreated, response.Success("Student created successfully", nil))
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	studentID := c.Param("student_id")

	student, err := h.service.GetByID(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Student not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Student fetched successfully", mapper.ToStudentDTO(student)))
}

func (h *StudentHandler) GetByClass(c *gin.Context) {
	classID := c.Param("class_id")

	students, err := h.service.GetByClassID(classID)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch student data")
		return
	}

	var dtoStudents []*dto.StudentResponse
	for _, student := range students {
		dtoStudents = append(dtoStudents, mapper.ToStudentDTO(student))
	}

	c.JSON(http.StatusOK, response.Success("Students fetched successfully", dtoStudents))
}

func (h *StudentHandler) GetAll(c *gin.Context) {
	classID := c.Query("class_id")
	search := c.Query("search")
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

	role, _ := c.Get("role")
	userID, _ := c.Get("userID")
	teacherID := ""
	if role == "teacher" {
		teacherID = userID.(string)
	}

	students, total, err := h.service.GetAll(teacherID, isActive, classID, search, page, limit)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch student data")
		return
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	var dtoStudents []*dto.StudentResponse
	for _, student := range students {
		dtoStudents = append(dtoStudents, mapper.ToStudentDTO(student))
	}

	c.JSON(http.StatusOK, response.SuccessWithPagination("Students fetched successfully", dtoStudents, gin.H{
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

	role, _ := c.Get("role")
	if role == "teacher" {
		userID, _ := c.Get("userID")
		studentToDelete, err := h.service.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, response.Error("Student not found"))
			return
		}
		class, err := h.classService.GetByID(studentToDelete.ClassID)
		if err != nil || class.TeacherID != userID.(string) {
			c.JSON(http.StatusForbidden, response.Error("You are not authorized to delete this student"))
			return
		}
	}

	if err := h.service.SoftDelete(id); err != nil {
		HandleAppError(c, err, "Failed to delete student")
		return
	}

	c.JSON(http.StatusOK, response.Success("Student deleted successfully", nil))
}

func (h *StudentHandler) Update(c *gin.Context) {
	studentID := c.Param("student_id")
	var req dto.UpdateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	role, _ := c.Get("role")
	if role == "teacher" {
		userID, _ := c.Get("userID")

		if req.ClassID != "" {
			class, err := h.classService.GetByID(req.ClassID)
			if err != nil || class.TeacherID != userID.(string) {
				c.JSON(http.StatusForbidden, response.Error("You are not authorized to move student to this class"))
				return
			}
		}

		studentToUpdate, err := h.service.GetByID(studentID)
		if err != nil {
			c.JSON(http.StatusNotFound, response.Error("Student not found"))
			return
		}
		oldClass, err := h.classService.GetByID(studentToUpdate.ClassID)
		if err != nil || oldClass.TeacherID != userID.(string) {
			c.JSON(http.StatusForbidden, response.Error("You are not authorized to modify this student"))
			return
		}
	}

	// Get current student to retrieve NISN
	currentStudent, err := h.service.GetByID(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Student not found"))
		return
	}

	// Resolve class name for QR code update
	var className string
	if req.ClassID != "" {
		targetClass, err := h.classService.GetByID(req.ClassID)
		if err == nil {
			className = targetClass.DisplayName
		}
	} else {
		className = currentStudent.ClassName
	}

	if err := h.service.Update(studentID, currentStudent.NISN, req.FullName, req.ClassID, className, req.DOB, req.Gender, *req.IsActive); err != nil {
		HandleAppError(c, err, "Failed to update student data")
		return
	}

	c.JSON(http.StatusOK, response.Success("Student updated successfully", nil))
}

func (h *StudentHandler) GetQRCode(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("Student QR Code generated (stub)", nil))
}

func (h *StudentHandler) ImportCSV(c *gin.Context) {
	classID := c.PostForm("class_id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, response.Error("class_id is required"))
		return
	}

	class, err := h.classService.GetByID(classID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid class ID"))
		return
	}

	role, _ := c.Get("role")
	if role == "teacher" {
		userID, _ := c.Get("userID")
		if class.TeacherID != userID.(string) {
			c.JSON(http.StatusForbidden, response.Error("You are not authorized to add students to this class"))
			return
		}
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("CSV file is required"))
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Failed to read CSV header"))
		return
	}

	if len(header) < 3 || strings.ToLower(strings.TrimSpace(header[0])) != "nisn" {
		c.JSON(http.StatusBadRequest, response.Error("Invalid CSV format. Expected header: nisn, nama_siswa, gender"))
		return
	}

	var entries []domain.StudentRawEntry
	rowIndex := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowIndex++
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(fmt.Sprintf("Failed to read CSV at row %d: %v", rowIndex, err)))
			return
		}

		if len(record) < 3 {
			continue
		}

		nisn := strings.TrimSpace(record[0])
		nama := strings.TrimSpace(record[1])
		genderRaw := strings.TrimSpace(record[2])

		if nisn == "" || nama == "" {
			continue
		}

		if len(nisn) != 10 {
			c.JSON(http.StatusBadRequest, response.Error(fmt.Sprintf("Validation failed at row %d: NISN must be 10 characters.", rowIndex)))
			return
		}

		if len(nama) < 3 {
			c.JSON(http.StatusBadRequest, response.Error(fmt.Sprintf("Validation failed at row %d: Student name (NISN %s) must be at least 3 characters.", rowIndex, nisn)))
			return
		}

		var gender string
		g := strings.ToUpper(genderRaw)
		if g == "L" {
			gender = "laki-laki"
		} else if g == "P" {
			gender = "perempuan"
		} else {
			c.JSON(http.StatusBadRequest, response.Error(fmt.Sprintf("Invalid gender format at row %d (NISN %s). Use only 'L' or 'P'.", rowIndex, nisn)))
			return
		}

		entries = append(entries, domain.StudentRawEntry{
			NISN:      nisn,
			FullName:  nama,
			ClassID:   classID,
			ClassName: class.DisplayName,
			Gender:    &gender,
		})
	}

	if len(entries) == 0 {
		c.JSON(http.StatusBadRequest, response.Error("No valid student data found in CSV"))
		return
	}

	if err := h.service.CreateBulkFromRaw(entries); err != nil {
		HandleAppError(c, err, "Failed to import student data")
		return
	}

	c.JSON(http.StatusCreated, response.Success(fmt.Sprintf("%d students imported successfully", len(entries)), nil))
}
