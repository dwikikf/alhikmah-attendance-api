package handler

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"alhikmah-attendance-api/core/domain"
	"alhikmah-attendance-api/core/dto"
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
		c.JSON(http.StatusBadRequest, response.ValidationError("Validasi gagal", utils.FormatValidationError(err)))
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

	student := &domain.Student{
		NISN:      req.NISN,
		FullName:  req.FullName,
		ClassID:   req.ClassID,
		ClassName: class.ClassName,
		Gender:    req.Gender,
		DOB:       req.DOB,
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
		handleDBError(c, err)
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
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Student deleted successfully", nil))
}

func (h *StudentHandler) Update(c *gin.Context) {
	studentID := c.Param("student_id")
	var req dto.UpdateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validasi gagal", utils.FormatValidationError(err)))
		return
	}

	role, _ := c.Get("role")
	if role == "teacher" {
		userID, _ := c.Get("userID")
		
		// Check target class ownership
		if req.ClassID != "" {
			class, err := h.classService.GetByID(req.ClassID)
			if err != nil || class.TeacherID != userID.(string) {
				c.JSON(http.StatusForbidden, response.Error("You are not authorized to move student to this class"))
				return
			}
		}

		// Check current student class ownership
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

	// Always get current student to get NISN
	studentToUpdate, err := h.service.GetByID(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Student not found"))
		return
	}

	// Get target class to get ClassName for QR Code
	var className string
	if req.ClassID != "" {
		targetClass, err := h.classService.GetByID(req.ClassID)
		if err == nil {
			className = targetClass.ClassName
		}
	} else {
		className = studentToUpdate.ClassName
	}

	student := &domain.Student{
		ID:        studentID,
		NISN:      studentToUpdate.NISN,
		FullName:  req.FullName,
		ClassID:   req.ClassID,
		ClassName: className,
		Gender:    req.Gender,
		DOB:       req.DOB,
		IsActive:  *req.IsActive,
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

func (h *StudentHandler) ImportCSV(c *gin.Context) {
	// Parse multipart form
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
	// Read header
	header, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Failed to read CSV header"))
		return
	}

	// Validate header (expected: nisn, nama_siswa, gender)
	if len(header) < 3 || strings.ToLower(strings.TrimSpace(header[0])) != "nisn" {
		c.JSON(http.StatusBadRequest, response.Error("Invalid CSV format. Expected header: nisn, nama_siswa, gender"))
		return
	}

	var students []*domain.Student
	rowIndex := 1 // header is row 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowIndex++
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(fmt.Sprintf("Gagal membaca CSV pada baris %d: %v", rowIndex, err)))
			return
		}

		if len(record) < 3 {
			continue // skip invalid rows
		}

		nisn := strings.TrimSpace(record[0])
		nama := strings.TrimSpace(record[1])
		genderRaw := strings.TrimSpace(record[2])

		if nisn == "" || nama == "" {
			continue // skip empty rows
		}

		if len(nisn) != 10 {
			c.JSON(http.StatusBadRequest, response.Error(fmt.Sprintf("Validasi gagal pada baris %d: NISN harus 10 karakter.", rowIndex)))
			return
		}

		if len(nama) < 3 {
			c.JSON(http.StatusBadRequest, response.Error(fmt.Sprintf("Validasi gagal pada baris %d: Nama siswa (NISN %s) minimal 3 karakter.", rowIndex, nisn)))
			return
		}

		var gender string
		g := strings.ToUpper(genderRaw)
		if g == "L" {
			gender = "laki-laki"
		} else if g == "P" {
			gender = "perempuan"
		} else {
			c.JSON(http.StatusBadRequest, response.Error(fmt.Sprintf("Format gender tidak valid pada baris %d (NISN %s). Hanya gunakan 'L' atau 'P'.", rowIndex, nisn)))
			return
		}

		student := &domain.Student{
			NISN:      nisn,
			FullName:  nama,
			ClassID:   classID,
			ClassName: class.ClassName,
			Gender:    &gender,
		}
		students = append(students, student)
	}

	if len(students) == 0 {
		c.JSON(http.StatusBadRequest, response.Error("No valid student data found in CSV"))
		return
	}

	if err := h.service.CreateBulk(students); err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success(fmt.Sprintf("%d students imported successfully", len(students)), nil))
}

