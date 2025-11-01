package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	// Import mock kita
	"github.com/Udean777/uang-bijak-go/internal/models"
	mocks "github.com/Udean777/uang-bijak-go/internal/repository/mocks"
)

func setupAuthService(t *testing.T) (AuthService, *mocks.MockUserRepository) {
	mockUserRepo := mocks.NewMockUserRepository(t)

	testSecret := "test_secret_key"
	testAccessTTL := time.Minute * 15
	testRefreshTTL := time.Hour * 24

	service := NewAuthService(mockUserRepo, testSecret, testAccessTTL, testRefreshTTL)
	return service, mockUserRepo
}

func TestAuthService_Register(t *testing.T) {
	service, mockUserRepo := setupAuthService(t)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		testUUID := uuid.New()

		// Kita harus "mengharapkan" (expect) panggilan ke CreateUser
		// Kita tidak bisa tahu persis hashed password-nya, jadi kita pakai 'mock.Anything'
		mockUserRepo.EXPECT().
			CreateUser(ctx, mock.AnythingOfType("*models.User")).
			Run(func(ctx context.Context, user *models.User) {
				// Cek apakah data yang dikirim ke repo sudah benar
				assert.Equal(t, "Test User", user.Name)
				assert.Equal(t, "test@example.com", user.Email)
				// Cek apakah password-nya di-hash (bukan plain text)
				assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123")))
			}).
			Return(testUUID, nil). // Kembalikan ID sukses
			Once()                 // Harapkan dipanggil 1x

		// 2. Act
		user, err := service.Register(ctx, "Test User", "test@example.com", "password123")

		// 3. Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUUID, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		// 1. Setup
		mockUserRepo.EXPECT().
			CreateUser(ctx, mock.AnythingOfType("*models.User")).
			Return(uuid.Nil, errors.New("unique constraint violation")). // Simulasikan error DB
			Once()

		// 2. Act
		user, err := service.Register(ctx, "Test User", "test@example.com", "password123")

		// 3. Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "unique constraint")
	})
}

func TestAuthService_Login(t *testing.T) {
	service, mockUserRepo := setupAuthService(t)
	ctx := context.Background()

	// Buat hash password yang valid untuk tes
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &models.User{
		ID:           uuid.New(),
		Email:        "user@example.com",
		PasswordHash: string(hashedPassword),
	}

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		mockUserRepo.EXPECT().
			GetUserByEmail(ctx, "user@example.com").
			Return(testUser, nil).
			Once()

		// 2. Act
		accessToken, refreshToken, err := service.Login(ctx, "user@example.com", "password123")

		// 3. Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)

		// Verifikasi token (opsional tapi bagus)
		accessID, err := service.ValidateToken(accessToken, "access")
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID, accessID)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// 1. Setup
		mockUserRepo.EXPECT().
			GetUserByEmail(ctx, "wrong@example.com").
			Return(nil, errors.New("not found")). // Simulasikan user tidak ada
			Once()

		// 2. Act
		_, _, err := service.Login(ctx, "wrong@example.com", "password123")

		// 3. Assert
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("Wrong Password", func(t *testing.T) {
		// 1. Setup
		mockUserRepo.EXPECT().
			GetUserByEmail(ctx, "user@example.com").
			Return(testUser, nil). // User ditemukan...
			Once()

		// 2. Act
		_, _, err := service.Login(ctx, "user@example.com", "wrongpassword") // ...tapi password salah

		// 3. Assert
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}
