package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Udean777/uang-bijak-go/internal/models"

	repoMocks "github.com/Udean777/uang-bijak-go/internal/repository/mocks"
)

func setupTransactionService(t *testing.T) (TransactionService, *repoMocks.MockTransactionRepository, *repoMocks.MockWalletRepository, *repoMocks.MockCategoryRepository) {

	mockTrxRepo := repoMocks.NewMockTransactionRepository(t)
	mockWalletRepo := repoMocks.NewMockWalletRepository(t)
	mockCategoryRepo := repoMocks.NewMockCategoryRepository(t)

	service := NewTransactionService(nil, mockTrxRepo, mockWalletRepo, mockCategoryRepo)
	return service, mockTrxRepo, mockWalletRepo, mockCategoryRepo
}

func TestTransactionService_GetUserTransactions(t *testing.T) {
	service, mockTrxRepo, _, _ := setupTransactionService(t)
	ctx := context.Background()
	testUserID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// 1. Setup
		mockResponse := []models.Transaction{
			{ID: 1, Amount: 10000, Type: "expense"},
		}

		mockTrxRepo.EXPECT().
			GetAllByUserID(ctx, testUserID).
			Return(mockResponse, nil).
			Once()

		// 2. Act
		trxs, err := service.GetUserTransactions(ctx, testUserID)

		// 3. Assert
		assert.NoError(t, err)
		assert.NotNil(t, trxs)
		assert.Equal(t, 1, len(trxs))
	})
}

// Kita hanya bisa menguji skenario GAGAL (Forbidden) secara unit test
// Skenario Sukses (Success) untuk CreateTransaction adalah Integration Test
func TestTransactionService_CreateTransaction_Failure_Forbidden(t *testing.T) {
	service, _, mockWalletRepo, mockCategoryRepo := setupTransactionService(t)
	ctx := context.Background()
	testUserID := uuid.New()
	req := models.CreateTransactionRequest{
		WalletID:   1,
		CategoryID: 1,
		Amount:     100,
		Type:       "expense",
	}

	t.Run("Fail - Wallet Ownership", func(t *testing.T) {
		// 1. Setup
		// Simulasikan walletRepo.CheckOwnership GAGAL
		mockWalletRepo.EXPECT().
			CheckOwnership(ctx, req.WalletID, testUserID).
			Return(nil, errors.New("not found")).
			Once()

		// 2. Act
		_, err := service.CreateTransaction(ctx, req, testUserID)

		// 3. Assert
		assert.Error(t, err)
		// Pastikan kita mengembalikan error 'forbidden'
		assert.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("Fail - Category Ownership", func(t *testing.T) {
		// 1. Setup
		// Simulasikan walletRepo.CheckOwnership SUKSES
		mockWalletRepo.EXPECT().
			CheckOwnership(ctx, req.WalletID, testUserID).
			Return(&models.Wallet{}, nil).
			Once()

		// Simulasikan categoryRepo.CheckOwnership GAGAL
		mockCategoryRepo.EXPECT().
			CheckOwnership(ctx, req.CategoryID, testUserID).
			Return(nil, errors.New("not found")).
			Once()

		// 2. Act
		_, err := service.CreateTransaction(ctx, req, testUserID)

		// 3. Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrForbidden)
	})
}
