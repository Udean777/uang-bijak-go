package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/Udean777/uang-bijak-go/internal/service"

	mocks "github.com/Udean777/uang-bijak-go/internal/service/mocks"
)

func TestWalletHandler_CreateWallet(t *testing.T) {
	mockService := mocks.NewMockWalletService(t)
	handler := NewWalletHandler(mockService)
	testUserID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		router := setupRouter()
		router.Use(func(c *gin.Context) { setAuthContext(c, testUserID) })
		router.POST("/wallets", handler.CreateWallet)

		reqBody := models.CreateWalletRequest{Name: "Dompet OVO", InitialBalance: 100000} // Rp 1000
		jsonBody, _ := json.Marshal(reqBody)

		mockResponse := &models.Wallet{
			ID:      1,
			UserID:  testUserID,
			Name:    "Dompet OVO",
			Balance: 100000,
		}

		mockService.EXPECT().
			CreateWallet(mock.Anything, reqBody, testUserID).
			Return(mockResponse, nil).
			Once()

		// 2. Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/wallets", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.Wallet
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Dompet OVO", resp.Name)
		assert.Equal(t, int64(100000), resp.Balance)
	})

	t.Run("Bad Request - Negative Balance", func(t *testing.T) {
		// 1. Setup
		router := setupRouter()
		router.Use(func(c *gin.Context) { setAuthContext(c, testUserID) })
		router.POST("/wallets", handler.CreateWallet)

		// Saldo awal negatif (gte=0 akan gagal)
		reqBody := `{"name": "Dompet Aneh", "initial_balance": -100}`

		// 2. Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/wallets", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		// Pastikan service tidak dipanggil
		mockService.AssertNotCalled(t, "CreateWallet")
	})
}

func TestWalletHandler_UpdateWallet(t *testing.T) {
	mockService := mocks.NewMockWalletService(t)
	handler := NewWalletHandler(mockService)
	testUserID := uuid.New()
	walletID := int64(1)

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		router := setupRouter()
		router.Use(func(c *gin.Context) { setAuthContext(c, testUserID) })
		router.PUT("/wallets/:id", handler.UpdateWallet)

		reqBody := models.UpdateWalletRequest{Name: "Dompet GoPay"}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.EXPECT().
			UpdateWallet(mock.Anything, walletID, reqBody, testUserID).
			Return(nil).
			Once()

		// 2. Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/wallets/"+strconv.FormatInt(walletID, 10), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Wallet updated successfully")
	})

	t.Run("Forbidden - Not Owner", func(t *testing.T) {
		// 1. Setup
		router := setupRouter()
		router.Use(func(c *gin.Context) { setAuthContext(c, testUserID) })
		router.PUT("/wallets/:id", handler.UpdateWallet)

		reqBody := models.UpdateWalletRequest{Name: "Dompet GoPay"}
		jsonBody, _ := json.Marshal(reqBody)

		// Simulasikan service mengembalikan error 'forbidden'
		mockService.EXPECT().
			UpdateWallet(mock.Anything, walletID, reqBody, testUserID).
			Return(service.ErrForbidden).
			Once()

		// 2. Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/wallets/"+strconv.FormatInt(walletID, 10), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// 3. Assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "not allowed")
	})
}
