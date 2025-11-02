package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/Udean777/uang-bijak-go/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, req models.CreateTransactionRequest, userID uuid.UUID) (*models.Transaction, error)
	GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
}

type transactionService struct {
	db           *pgxpool.Pool
	trxRepo      repository.TransactionRepository
	walletRepo   repository.WalletRepository
	categoryRepo repository.CategoryRepository
}

func NewTransactionService(db *pgxpool.Pool, trxRepo repository.TransactionRepository, walletRepo repository.WalletRepository, categoryRepo repository.CategoryRepository) TransactionService {
	return &transactionService{
		db:           db,
		trxRepo:      trxRepo,
		walletRepo:   walletRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, req models.CreateTransactionRequest, userID uuid.UUID) (*models.Transaction, error) {

	if _, err := s.walletRepo.CheckOwnership(ctx, req.WalletID, userID); err != nil {
		return nil, fmt.Errorf("wallet ownership validation failed: %w", ErrForbidden)
	}
	if _, err := s.categoryRepo.CheckOwnership(ctx, req.CategoryID, userID); err != nil {
		return nil, fmt.Errorf("category ownership validation failed: %w", ErrForbidden)
	}

	balanceChange := req.Amount
	if models.TransactionType(req.Type) == models.TransactionExpense {
		balanceChange = -req.Amount
	}

	t := &models.Transaction{
		UserID:      userID,
		WalletID:    req.WalletID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Type:        models.TransactionType(req.Type),
		Description: req.Description,
	}
	if req.TransactionDate != nil {
		t.TransactionDate = *req.TransactionDate
	} else {
		t.TransactionDate = time.Now()
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	if err := s.walletRepo.UpdateBalanceTx(ctx, tx, req.WalletID, balanceChange); err != nil {
		return nil, err
	}

	if err := s.trxRepo.CreateTx(ctx, tx, t); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *transactionService) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {

	return s.trxRepo.GetAllByUserID(ctx, userID)
}
