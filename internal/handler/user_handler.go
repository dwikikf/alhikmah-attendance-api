package handler

import (
	"net/http"
	"strconv"

	"alhikmah-attendance-api/internal/domain"

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	user, err := h.service.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
		"pagination": gin.H{
			"page":            page,
			"pageSize":        limit,
			"totalItems":      total,
			"totalPages":      totalPages,
			"hasNextPage":     page < totalPages,
			"hasPreviousPage": page > 1,
		},
	})
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required,min=3"`
		Role     string `json:"role" binding:"required,oneof=admin teacher"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password, // Service will hash this
		FullName:     req.FullName,
		Role:         req.Role,
	}

	if err := h.service.Create(user); err != nil {
		handleDBError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
		"message": "User created successfully",
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("user_id")

	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For security, normally users can only update themselves unless they are an admin.
	// We'll trust the route middleware to have checked this.

	user := &domain.User{
		ID:       id,
		FullName: req.FullName,
		Email:    req.Email,
	}

	if err := h.service.Update(user); err != nil {
		handleDBError(c, err)
		return
	}

	// Fetch updated user to return complete data
	updatedUser, _ := h.service.GetByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedUser,
		"message": "User updated successfully",
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("user_id")

	if err := h.service.SoftDelete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}
