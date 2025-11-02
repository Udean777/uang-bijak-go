package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Udean777/uang-bijak-go/internal/models"
	// Import mock service
	serviceMocks "github.com/Udean777/uang-bijak-go/internal/service/mocks"
)

// (Asumsikan helper setupRouter dan setAuthContext sudah ada)

func TestDashboardHandler_GetDashboardSummary(t *testing.T) {
	mockService := serviceMocks.NewMockDashboardService(t)
	handler := NewDashboardHandler(mockService)
	testUserID := uuid.New()

	mockResponse := &models.DashboardSummary{
		TotalBalance: 100000,
		TotalIncome:  50000,
		TotalExpense: 20000,
	}

	t.Run("Success - Default (Bulan Ini)", func(t *testing.T) {
		// 1. Setup
		router := setupRouter()
		router.Use(func(c *gin.Context) { setAuthContext(c, testUserID) })
		router.GET("/dashboard", handler.GetDashboardSummary)

		// Harapkan panggilan service. startTime dan endTime akan
		// di-resolve oleh GetDateRange(), jadi kita gunakan mock.Anything
		mockService.EXPECT().
			GetDashboardSummary(mock.Anything, testUserID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return(mockResponse, nil).
			Once()

		// 2. Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/dashboard", nil)

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.DashboardSummary
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, mockResponse.TotalBalance, resp.TotalBalance)
	})

	t.Run("Success - Dengan Query Parameter (Oktober 2025)", func(t *testing.T) {
		// 1. Setup
		router := setupRouter()
		router.Use(func(c *gin.Context) { setAuthContext(c, testUserID) })
		router.GET("/dashboard", handler.GetDashboardSummary)

		// Definisikan rentang waktu yang kita harapkan
		loc, _ := time.LoadLocation("Asia/Jakarta")
		expectedStart := time.Date(2025, time.October, 1, 0, 0, 0, 0, loc)
		expectedEnd := time.Date(2025, time.October, 31, 23, 59, 59, 999999999, loc)

		// Harapkan panggilan service dengan rentang waktu yang TEPAT
		mockService.EXPECT().
			GetDashboardSummary(mock.Anything, testUserID, expectedStart, expectedEnd).
			Return(mockResponse, nil).
			Once()

		// 2. Act
		w := httptest.NewRecorder()
		// Request dengan query parameter
		req, _ := http.NewRequest(http.MethodGet, "/dashboard?month=10&year=2025", nil)

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.DashboardSummary
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, mockResponse.TotalIncome, resp.TotalIncome)
	})
}
