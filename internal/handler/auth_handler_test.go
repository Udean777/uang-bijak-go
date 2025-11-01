package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Udean777/uang-bijak-go/internal/models"

	mocks "github.com/Udean777/uang-bijak-go/internal/service/mocks"
)

func setupRouter() *gin.Engine {
	// Matikan output logging Gin agar tidak mengganggu tes
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func TestAuthHandler_Register(t *testing.T) {
	// 1. Setup
	mockAuthService := mocks.NewMockAuthService(t)
	handler := NewAuthHandler(mockAuthService)

	router := setupRouter()
	router.POST("/register", handler.Register)

	t.Run("Success", func(t *testing.T) {
		// Data mock yang akan dikembalikan service
		mockUser := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		// Setup body request
		reqBody := RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		// Setup mock expectation
		mockAuthService.EXPECT().
			Register(mock.Anything, "Test User", "test@example.com", "password123").
			Return(mockUser, nil).
			Once()

		// 2. Act
		// Buat request dan response recorder
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req) // Jalankan request

		// 3. Assert
		assert.Equal(t, http.StatusCreated, w.Code) // Cek status code

		// Cek body response
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "User registered successfully", response["message"])
		assert.Equal(t, mockUser.ID.String(), response["user_id"])
	})

	t.Run("Invalid Input (Bad Request)", func(t *testing.T) {
		// Setup body request yang tidak valid (password hilang)
		reqBody := `{"name": "Test User", "email": "test@example.com"}`

		// 2. Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusBadRequest, w.Code) // Harapannya 400 Bad Request

		// Pastikan service TIDAK dipanggil
		mockAuthService.AssertNotCalled(t, "Register")
	})
}
