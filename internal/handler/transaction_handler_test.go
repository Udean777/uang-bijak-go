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
	"github.com/Udean777/uang-bijak-go/internal/service"

	mocks "github.com/Udean777/uang-bijak-go/internal/service/mocks"
)

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	mockService := mocks.NewMockTransactionService(t)
	handler := NewTransactionHandler(mockService)
	testUserID := uuid.New()

	t.Run("Success - Expense", func(t *testing.T) {
		// 1. Setup
		router := setupRouter()
		router.Use(func(c *gin.Context) { setAuthContext(c, testUserID) })
		router.POST("/transactions", handler.CreateTransaction)

		reqBody := models.CreateTransactionRequest{
			WalletID:   1,
			CategoryID: 1,
			Amount:     20000, // Rp 200
			Type:       "expense",
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockResponse := &models.Transaction{
			ID:         1,
			UserID:     testUserID,
			WalletID:   1,
			CategoryID: 1,
			Amount:     20000,
			Type:       "expense",
		}

		mockService.EXPECT().
			CreateTransaction(mock.Anything, reqBody, testUserID).
			Return(mockResponse, nil).
			Once()

		// 2. Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.Transaction
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, int64(20000), resp.Amount)
		assert.Equal(t, models.TransactionType("expense"), resp.Type)
	})

	t.Run("Bad Request - Invalid Type", func(t *testing.T) {
		// 1. Setup
		router := setupRouter()
		router.Use(func(c *gin.Context) { setAuthContext(c, testUserID) })
		router.POST("/transactions", handler.CreateTransaction)

		// Tipe 'transfer' tidak valid (oneof=expense income)
		reqBody := `{"wallet_id": 1, "category_id": 1, "amount": 100, "type": "transfer"}`

		// 2. Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CreateTransaction")
	})

	t.Run("Forbidden - Invalid Wallet/Category", func(t *testing.T) {
		// 1. Setup
		router := setupRouter()
		router.Use(func(c *gin.Context) { setAuthContext(c, testUserID) })
		router.POST("/transactions", handler.CreateTransaction)

		reqBody := models.CreateTransactionRequest{
			WalletID:   99, // ID dompet milik user lain
			CategoryID: 1,
			Amount:     20000,
			Type:       "expense",
		}
		jsonBody, _ := json.Marshal(reqBody)

		// Simulasikan service mengembalikan error 'forbidden'
		mockService.EXPECT().
			CreateTransaction(mock.Anything, reqBody, testUserID).
			Return(nil, service.ErrForbidden).
			Once()

		// 2. Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid wallet or category ID")
	})
}
