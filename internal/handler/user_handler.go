package handler

import (
	"net/http"
	"strconv"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/internal/dto"
	"alhikmah-attendance-api/internal/mapper"
	"alhikmah-attendance-api/pkg/response"
	"alhikmah-attendance-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error("User ID not found in context"))
		return
	}

	user, err := h.service.GetByID(userID.(string))
	if err != nil {
		HandleAppError(c, err, "Failed to fetch user profile data")
		return
	}

	c.JSON(http.StatusOK, response.Success("User fetched successfully", mapper.ToUserDTO(user)))
}

func (h *UserHandler) GetAll(c *gin.Context) {
	role := c.Query("role")
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

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	users, total, err := h.service.GetAll(role, isActive, page, limit)
	if err != nil {
		HandleAppError(c, err, "Failed to fetch user data")
		return
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	var dtoUsers []*dto.UserResponse
	for _, user := range users {
		dtoUsers = append(dtoUsers, mapper.ToUserDTO(user))
	}

	c.JSON(http.StatusOK, response.SuccessWithPagination("Users fetched successfully", dtoUsers, gin.H{
		"page":            page,
		"pageSize":        limit,
		"totalItems":      total,
		"totalPages":      totalPages,
		"hasNextPage":     page < totalPages,
		"hasPreviousPage": page > 1,
	}))
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	if err := h.service.Create(req.Username, req.Email, req.Password, req.FullName, req.Role); err != nil {
		HandleAppError(c, err, "Failed to create new user")
		return
	}

	c.JSON(http.StatusCreated, response.Success("User created successfully", nil))
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("user_id")

	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
		return
	}

	if err := h.service.Update(id, req.FullName, req.Email, req.Password); err != nil {
		HandleAppError(c, err, "Failed to update user")
		return
	}

	// Fetch updated user to return complete data
	updatedUser, _ := h.service.GetByID(id)

	c.JSON(http.StatusOK, response.Success("User updated successfully", mapper.ToUserDTO(updatedUser)))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("user_id")

	if err := h.service.SoftDelete(id); err != nil {
		HandleAppError(c, err, "Failed to delete user")
		return
	}

	c.JSON(http.StatusOK, response.Success("User deleted successfully", nil))
}
