package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Helper untuk membuat token
func generateTestToken(t *testing.T, secret string, userID uuid.UUID, tokenType string, ttl time.Duration) string {
	claims := jwt.MapClaims{
		"sub":        userID.String(),
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(ttl).Unix(),
		"token_type": tokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)
	return signedToken
}

func TestAuthMiddleware(t *testing.T) {
	testSecret := "my-secret-key"
	testUserID := uuid.New()

	// Setup router dengan middleware
	router := gin.Default()
	router.Use(AuthMiddleware(testSecret))
	// Buat dummy handler yang hanya bisa diakses jika middleware lolos
	router.GET("/protected", func(c *gin.Context) {
		// Cek apakah userID di-set di context
		id, exists := c.Get("userID")
		assert.True(t, exists)
		assert.Equal(t, testUserID, id)

		c.JSON(http.StatusOK, gin.H{"message": "WELCOME_BACK"})
	})

	t.Run("Success - Valid Access Token", func(t *testing.T) {
		token := generateTestToken(t, testSecret, testUserID, "access", time.Minute*15)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "WELCOME_BACK")
	})

	t.Run("Fail - No Header", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Fail - Invalid Header Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bear token") // Format salah

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Fail - Expired Token", func(t *testing.T) {
		// Buat token yang sudah kadaluarsa 1 jam lalu
		token := generateTestToken(t, testSecret, testUserID, "access", -time.Hour)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Fail - Wrong Token Type (Refresh Token)", func(t *testing.T) {
		// Gunakan 'refresh' token untuk akses
		token := generateTestToken(t, testSecret, testUserID, "refresh", time.Hour)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "expected 'access' token")
	})
}
