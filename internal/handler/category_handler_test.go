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

// setupCategoryTestRouter adalah fungsi helper untuk menyiapkan router Gin dalam mode tes.
func setupCategoryTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

// setAuthContext adalah fungsi helper untuk menyetel userID dalam konteks Gin,
// mensimulasikan pengguna yang sudah terautentikasi.
func setAuthContext(c *gin.Context, userID uuid.UUID) {
	c.Set("userID", userID)
}

// TestCategoryHandler_CreateCategory menguji fungsionalitas handler CreateCategory.
func TestCategoryHandler_CreateCategory(t *testing.T) {
	mockService := mocks.NewMockCategoryService(t)
	handler := NewCategoryHandler(mockService)

	testUserID := uuid.New()

	// Menguji kasus sukses pembuatan kategori.
	t.Run("Success", func(t *testing.T) {
		router := setupRouter()
		router.Use(func(c *gin.Context) {
			setAuthContext(c, testUserID)
			c.Next()
		})

		router.POST("/categories", handler.CreateCategory)

		reqBody := models.UpsertCategoryRequest{Name: "Transportasi"}
		jsonBody, _ := json.Marshal(reqBody)

		mockResponse := &models.Category{
			ID:     1,
			UserID: testUserID,
			Name:   "Transportasi",
		}

		mockService.EXPECT().
			CreateCategory(mock.Anything, reqBody, testUserID).
			Return(mockResponse, nil).
			Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.Category
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Transportasi", resp.Name)
	})

	// Menguji kasus di mana request body berisi JSON yang tidak valid.
	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		router := setupRouter()
		router.Use(func(c *gin.Context) {
			setAuthContext(c, testUserID)
			c.Next()
		})
		router.POST("/categories", handler.CreateCategory)

		invalidJsonBody := `{"name":`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString(invalidJsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CreateCategory")
	})
}

// TestCategoryHandler_UpdateCategory menguji fungsionalitas handler UpdateCategory.
func TestCategoryHandler_UpdateCategory(t *testing.T) {
	mockService := mocks.NewMockCategoryService(t)
	handler := NewCategoryHandler(mockService)

	testUserID := uuid.New()
	categoryID := int64(1)

	// Menguji kasus sukses pembaruan kategori.
	t.Run("Success", func(t *testing.T) {
		router := setupCategoryTestRouter()
		router.Use(func(c *gin.Context) {
			setAuthContext(c, testUserID)
			c.Next()
		})
		router.PUT("/categories/:id", handler.UpdateCategory)

		reqBody := models.UpsertCategoryRequest{Name: "Transportasi Baru"}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.EXPECT().
			UpdateCategory(mock.Anything, categoryID, reqBody, testUserID).
			Return(nil).
			Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/categories/"+strconv.FormatInt(categoryID, 10), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Menguji kasus di mana pengguna mencoba memperbarui kategori yang bukan miliknya.
	t.Run("Forbidden - Not Owner", func(t *testing.T) {
		router := setupCategoryTestRouter()
		router.Use(func(c *gin.Context) {
			setAuthContext(c, testUserID)
			c.Next()
		})
		router.PUT("/categories/:id", handler.UpdateCategory)

		reqBody := models.UpsertCategoryRequest{Name: "Transportasi Baru"}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.EXPECT().
			UpdateCategory(mock.Anything, categoryID, reqBody, testUserID).
			Return(service.ErrForbidden).
			Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/categories/"+strconv.FormatInt(categoryID, 10), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "not allowed")
	})

	// Menguji kasus di mana ID kategori yang diberikan di URL tidak valid.
	t.Run("Bad Request - Invalid ID", func(t *testing.T) {
		router := setupCategoryTestRouter()
		router.Use(func(c *gin.Context) {
			setAuthContext(c, testUserID)
			c.Next()
		})
		router.PUT("/categories/:id", handler.UpdateCategory)

		reqBody := models.UpsertCategoryRequest{Name: "Transportasi Baru"}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/categories/abc", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid category ID")
		mockService.AssertNotCalled(t, "UpdateCategory")
	})
}
