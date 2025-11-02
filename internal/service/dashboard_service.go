package service

import (
	"context"
	"time"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/Udean777/uang-bijak-go/internal/repository"
	"github.com/google/uuid"
)

// DashboardService interface
type DashboardService interface {
	GetDashboardSummary(ctx context.Context, userID uuid.UUID, startTime time.Time, endTime time.Time) (*models.DashboardSummary, error)
}

// dashboardService struct
type dashboardService struct {
	walletRepo repository.WalletRepository
	trxRepo    repository.TransactionRepository
}

// NewDashboardService constructor
func NewDashboardService(walletRepo repository.WalletRepository, trxRepo repository.TransactionRepository) DashboardService {
	return &dashboardService{
		walletRepo: walletRepo,
		trxRepo:    trxRepo,
	}
}

// GetDashboardSummary implementation
func (s *dashboardService) GetDashboardSummary(ctx context.Context, userID uuid.UUID, startTime time.Time, endTime time.Time) (*models.DashboardSummary, error) {

	// 1. Ambil Total Saldo
	totalBalance, err := s.walletRepo.GetTotalBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Ambil Total Pemasukan & Pengeluaran
	totalIncome, totalExpense, err := s.trxRepo.GetTotalIncomeAndExpense(ctx, userID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 3. Gabungkan hasilnya
	summary := &models.DashboardSummary{
		TotalBalance: totalBalance,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
	}

	return summary, nil
}
