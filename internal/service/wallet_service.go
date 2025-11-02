package service

import (
	"context"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/Udean777/uang-bijak-go/internal/repository"
	"github.com/google/uuid"
)

type WalletService interface {
	CreateWallet(ctx context.Context, req models.CreateWalletRequest, userID uuid.UUID) (*models.Wallet, error)
	GetUserWallets(ctx context.Context, userID uuid.UUID) ([]models.Wallet, error)
	UpdateWallet(ctx context.Context, walletID int64, req models.UpdateWalletRequest, userID uuid.UUID) error
	DeleteWallet(ctx context.Context, walletID int64, userID uuid.UUID) error
}

type walletService struct {
	walletRepo repository.WalletRepository
	// TODO: Tambahkan transactionRepo di sini
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{walletRepo: repo}
}

func (s *walletService) CreateWallet(ctx context.Context, req models.CreateWalletRequest, userID uuid.UUID) (*models.Wallet, error) {
	wallet := &models.Wallet{
		UserID:  userID,
		Name:    req.Name,
		Balance: req.InitialBalance,
	}

	_, err := s.walletRepo.Create(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]models.Wallet, error) {
	return s.walletRepo.GetAllByUserID(ctx, userID)
}

func (s *walletService) UpdateWallet(ctx context.Context, walletID int64, req models.UpdateWalletRequest, userID uuid.UUID) error {
	_, err := s.walletRepo.CheckOwnership(ctx, walletID, userID)
	if err != nil {
		return ErrForbidden
	}

	return s.walletRepo.Update(ctx, walletID, req.Name)
}

func (s *walletService) DeleteWallet(ctx context.Context, walletID int64, userID uuid.UUID) error {
	_, err := s.walletRepo.CheckOwnership(ctx, walletID, userID)
	if err != nil {
		return ErrForbidden
	}

	// TODO: Tambahkan pengecekan di sini
	// "Apakah dompet ini masih memiliki transaksi?"
	// "Apakah saldonya nol?"
	// Untuk saat ini, kita izinkan delete.

	return s.walletRepo.Delete(ctx, walletID)
}
