package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	repoMocks "github.com/Udean777/uang-bijak-go/internal/repository/mocks"
)

// Helper setup
func setupDashboardService(t *testing.T) (DashboardService, *repoMocks.MockWalletRepository, *repoMocks.MockTransactionRepository) {
	mockWalletRepo := repoMocks.NewMockWalletRepository(t)
	mockTrxRepo := repoMocks.NewMockTransactionRepository(t)
	service := NewDashboardService(mockWalletRepo, mockTrxRepo)
	return service, mockWalletRepo, mockTrxRepo
}

func TestDashboardService_GetDashboardSummary(t *testing.T) {
	service, mockWalletRepo, mockTrxRepo := setupDashboardService(t)
	ctx := context.Background()
	testUserID := uuid.New()
	startTime := time.Now()
	endTime := time.Now()

	t.Run("Success", func(t *testing.T) {
		// 1. Setup Mock
		// Harapkan panggilan ke WalletRepo, kembalikan total saldo 1.000.000
		mockWalletRepo.EXPECT().
			GetTotalBalanceByUserID(ctx, testUserID).
			Return(int64(1000000), nil).
			Once()

		// Harapkan panggilan ke TrxRepo, kembalikan income 500.000, expense 150.000
		mockTrxRepo.EXPECT().
			GetTotalIncomeAndExpense(ctx, testUserID, startTime, endTime).
			Return(int64(500000), int64(150000), nil).
			Once()

		// 2. Act
		summary, err := service.GetDashboardSummary(ctx, testUserID, startTime, endTime)

		// 3. Assert
		assert.NoError(t, err)
		assert.NotNil(t, summary)
		assert.Equal(t, int64(1000000), summary.TotalBalance)
		assert.Equal(t, int64(500000), summary.TotalIncome)
		assert.Equal(t, int64(150000), summary.TotalExpense)
	})

	t.Run("Fail - WalletRepo Fails", func(t *testing.T) {
		// 1. Setup Mock
		mockWalletRepo.EXPECT().
			GetTotalBalanceByUserID(ctx, testUserID).
			Return(int64(0), errors.New("db error")).
			Once()

		// 2. Act
		summary, err := service.GetDashboardSummary(ctx, testUserID, startTime, endTime)

		// 3. Assert
		assert.Error(t, err)
		assert.Nil(t, summary)
		// Pastikan TrxRepo tidak dipanggil jika WalletRepo gagal
		mockTrxRepo.AssertNotCalled(t, "GetTotalIncomeAndExpense")
	})

	t.Run("Fail - TrxRepo Fails", func(t *testing.T) {
		// 1. Setup Mock
		mockWalletRepo.EXPECT().
			GetTotalBalanceByUserID(ctx, testUserID).
			Return(int64(1000000), nil). // Ini sukses
			Once()

		mockTrxRepo.EXPECT().
			GetTotalIncomeAndExpense(ctx, testUserID, startTime, endTime).
			Return(int64(0), int64(0), errors.New("db error")). // Ini gagal
			Once()

		// 2. Act
		summary, err := service.GetDashboardSummary(ctx, testUserID, startTime, endTime)

		// 3. Assert
		assert.Error(t, err)
		assert.Nil(t, summary)
	})
}
