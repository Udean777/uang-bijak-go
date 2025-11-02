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
func setupWalletService(t *testing.T) (WalletService, *mocks.MockWalletRepository) {
	mockRepo := mocks.NewMockWalletRepository(t)
	service := NewWalletService(mockRepo)
	return service, mockRepo
}

func TestWalletService_CreateWallet(t *testing.T) {
	service, mockRepo := setupWalletService(t)
	ctx := context.Background()
	testUserID := uuid.New()
	req := models.CreateWalletRequest{Name: "Dompet Tunai", InitialBalance: 50000} // 50000 sen = Rp 500

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*models.Wallet")).
			Run(func(ctx context.Context, w *models.Wallet) {
				// Cek data yang dikirim ke repo
				assert.Equal(t, testUserID, w.UserID)
				assert.Equal(t, "Dompet Tunai", w.Name)
				assert.Equal(t, int64(50000), w.Balance)
				// Simulasikan repo mengatur ID
				w.ID = 1
			}).
			Return(int64(1), nil).
			Once()

		// 2. Act
		wallet, err := service.CreateWallet(ctx, req, testUserID)

		// 3. Assert
		assert.NoError(t, err)
		assert.NotNil(t, wallet)
		assert.Equal(t, "Dompet Tunai", wallet.Name)
		assert.Equal(t, int64(1), wallet.ID)
		assert.Equal(t, int64(50000), wallet.Balance)
	})
}

func TestWalletService_UpdateWallet(t *testing.T) {
	service, mockRepo := setupWalletService(t)
	ctx := context.Background()

	testUserID := uuid.New()
	walletID := int64(1)
	req := models.UpdateWalletRequest{Name: "Dompet BCA"}

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		// Harapkan panggilan ke checkOwnership (sukses)
		mockRepo.EXPECT().
			CheckOwnership(ctx, walletID, testUserID).
			Return(&models.Wallet{ID: walletID, UserID: testUserID}, nil).
			Once()

		// Harapkan panggilan ke Update (sukses)
		mockRepo.EXPECT().
			Update(ctx, walletID, "Dompet BCA").
			Return(nil).
			Once()

		// 2. Act
		err := service.UpdateWallet(ctx, walletID, req, testUserID)

		// 3. Assert
		assert.NoError(t, err)
	})

	t.Run("Fail - Forbidden (Not Owner)", func(t *testing.T) {
		// 1. Setup
		otherUserID := uuid.New()
		// Harapkan panggilan ke checkOwnership (gagal)
		mockRepo.EXPECT().
			CheckOwnership(ctx, walletID, otherUserID).
			Return(nil, errors.New("not found")). // Simulasikan 'not found'
			Once()

		// 2. Act
		err := service.UpdateWallet(ctx, walletID, req, otherUserID)

		// 3. Assert
		assert.Error(t, err)
		assert.Equal(t, ErrForbidden, err)

		// Pastikan Update TIDAK pernah dipanggil
		mockRepo.AssertNotCalled(t, "Update")
	})
}

func TestWalletService_DeleteWallet(t *testing.T) {
	service, mockRepo := setupWalletService(t)
	ctx := context.Background()

	testUserID := uuid.New()
	walletID := int64(1)

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		mockRepo.EXPECT().
			CheckOwnership(ctx, walletID, testUserID).
			Return(&models.Wallet{ID: walletID, UserID: testUserID}, nil).
			Once()

		mockRepo.EXPECT().
			Delete(ctx, walletID).
			Return(nil).
			Once()

		// 2. Act
		err := service.DeleteWallet(ctx, walletID, testUserID)

		// 3. Assert
		assert.NoError(t, err)
	})

	t.Run("Fail - Forbidden (Not Owner)", func(t *testing.T) {
		// 1. Setup
		otherUserID := uuid.New()
		mockRepo.EXPECT().
			CheckOwnership(ctx, walletID, otherUserID).
			Return(nil, errors.New("not found")).
			Once()

		// 2. Act
		err := service.DeleteWallet(ctx, walletID, otherUserID)

		// 3. Assert
		assert.Error(t, err)
		assert.Equal(t, ErrForbidden, err)

		mockRepo.AssertNotCalled(t, "Delete")
	})
}
