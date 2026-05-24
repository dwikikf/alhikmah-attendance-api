package handler

import (
	"net/http"
	"strconv"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service domain.DashboardService
}

func NewDashboardHandler(service domain.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetRecentActivity(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	activities, err := h.service.GetRecentActivity(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to fetch recent activity"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Recent activity fetched successfully", activities))
}

func (h *DashboardHandler) GetAttendanceTrend(c *gin.Context) {
	daysStr := c.Query("days")
	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	trends, err := h.service.GetAttendanceTrend(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to fetch attendance trend"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Attendance trend fetched successfully", trends))
}
