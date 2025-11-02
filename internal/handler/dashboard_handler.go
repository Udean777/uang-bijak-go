package handler

import (
	"net/http"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/Udean777/uang-bijak-go/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: svc}
}

func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	userID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var query models.DashboardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	startTime, endTime := query.GetDateRange()

	summary, err := h.dashboardService.GetDashboardSummary(c.Request.Context(), userID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
