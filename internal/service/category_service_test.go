package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Udean777/uang-bijak-go/internal/models"
	mocks "github.com/Udean777/uang-bijak-go/internal/repository/mocks"
)

// Helper setup
func setupCategoryService(t *testing.T) (CategoryService, *mocks.MockCategoryRepository) {
	mockRepo := mocks.NewMockCategoryRepository(t)
	service := NewCategoryService(mockRepo)
	return service, mockRepo
}

func TestCategoryService_CreateCategory(t *testing.T) {
	service, mockRepo := setupCategoryService(t)
	ctx := context.Background()
	testUserID := uuid.New()
	req := models.UpsertCategoryRequest{Name: "Makanan"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*models.Category")).
			// Tambahkan Run() untuk menyimulasikan modifikasi pointer
			Run(func(ctx context.Context, cat *models.Category) {
				// Cek data yang dikirim ke repo
				assert.Equal(t, testUserID, cat.UserID)
				assert.Equal(t, "Makanan", cat.Name)

				// SIMULASIKAN REPO: Set ID pada pointer
				cat.ID = 1
			}).
			Return(int64(1), nil). // Return ini tetap diperlukan untuk mencocokkan signature
			Once()

		// 2. Act
		cat, err := service.CreateCategory(ctx, req, testUserID)

		// 3. Assert
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, "Makanan", cat.Name)
		// Baris ini sekarang akan berhasil
		assert.Equal(t, int64(1), cat.ID)
	})

	t.Run("Fail - Duplicate", func(t *testing.T) {
		// 1. Setup
		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*models.Category")).
			Return(int64(0), errors.New("unique constraint violation")).
			Once()

		// 2. Act
		_, err := service.CreateCategory(ctx, req, testUserID)

		// 3. Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unique constraint")
	})
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	service, mockRepo := setupCategoryService(t)
	ctx := context.Background()

	testUserID := uuid.New()
	otherUserID := uuid.New()
	categoryID := int64(1)
	req := models.UpsertCategoryRequest{Name: "Updated Makanan"}

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		// Harapkan panggilan ke checkOwnership (sukses)
		mockRepo.EXPECT().
			CheckOwnership(ctx, categoryID, testUserID).
			Return(&models.Category{ID: categoryID, UserID: testUserID}, nil).
			Once()

		// Harapkan panggilan ke Update (sukses)
		mockRepo.EXPECT().
			Update(ctx, categoryID, "Updated Makanan").
			Return(nil).
			Once()

		// 2. Act
		err := service.UpdateCategory(ctx, categoryID, req, testUserID)

		// 3. Assert
		assert.NoError(t, err)
	})

	t.Run("Fail - Forbidden (Not Owner)", func(t *testing.T) {
		// 1. Setup
		// Harapkan panggilan ke checkOwnership (gagal)
		mockRepo.EXPECT().
			CheckOwnership(ctx, categoryID, otherUserID).
			Return(nil, errors.New("not found")). // Simulasikan 'not found' atau error
			Once()

		// 2. Act
		err := service.UpdateCategory(ctx, categoryID, req, otherUserID)

		// 3. Assert
		assert.Error(t, err)
		assert.Equal(t, ErrForbidden, err)

		// Pastikan Update TIDAK pernah dipanggil
		mockRepo.AssertNotCalled(t, "Update")
	})
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	service, mockRepo := setupCategoryService(t)
	ctx := context.Background()

	testUserID := uuid.New()
	otherUserID := uuid.New()
	categoryID := int64(1)

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		// Harapkan panggilan ke checkOwnership (sukses)
		mockRepo.EXPECT().
			CheckOwnership(ctx, categoryID, testUserID).
			Return(&models.Category{ID: categoryID, UserID: testUserID}, nil).
			Once()

		// Harapkan panggilan ke Delete (sukses)
		mockRepo.EXPECT().
			Delete(ctx, categoryID).
			Return(nil).
			Once()

		// 2. Act
		err := service.DeleteCategory(ctx, categoryID, testUserID)

		// 3. Assert
		assert.NoError(t, err)
	})

	t.Run("Fail - Forbidden (Not Owner)", func(t *testing.T) {
		// 1. Setup
		// Harapkan panggilan ke checkOwnership (gagal)
		mockRepo.EXPECT().
			CheckOwnership(ctx, categoryID, otherUserID).
			Return(nil, errors.New("not found")).
			Once()

		// 2. Act
		err := service.DeleteCategory(ctx, categoryID, otherUserID)

		// 3. Assert
		assert.Error(t, err)
		assert.Equal(t, ErrForbidden, err)

		// Pastikan Delete TIDAK pernah dipanggil
		mockRepo.AssertNotCalled(t, "Delete")
	})
}
