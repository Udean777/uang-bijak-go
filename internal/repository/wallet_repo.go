package repository

import (
	"context"
	"time"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *models.Wallet) (int64, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Wallet, error)
	GetByID(ctx context.Context, id int64) (*models.Wallet, error)
	Update(ctx context.Context, id int64, name string) error
	Delete(ctx context.Context, id int64) error
	CheckOwnership(ctx context.Context, walletID int64, userID uuid.UUID) (*models.Wallet, error)
	// TODO: Tambahkan UpdateBalance(ctx, id, amount) di sini untuk dipakai oleh service Transaksi
}

type walletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(ctx context.Context, wallet *models.Wallet) (int64, error) {
	query := `INSERT INTO wallets (user_id, name, balance) VALUES ($1, $2, $3) 
	          RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, wallet.UserID, wallet.Name, wallet.Balance).Scan(
		&wallet.ID,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		return 0, err
	}

	return wallet.ID, nil
}

func (r *walletRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Wallet, error) {
	query := `SELECT id, name, balance, created_at, updated_at FROM wallets 
	          WHERE user_id = $1 ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []models.Wallet
	for rows.Next() {
		var w models.Wallet
		if err := rows.Scan(&w.ID, &w.Name, &w.Balance, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}

	return wallets, nil
}

func (r *walletRepository) GetByID(ctx context.Context, id int64) (*models.Wallet, error) {
	query := `SELECT id, user_id, name, balance, created_at, updated_at FROM wallets WHERE id = $1`
	var w models.Wallet

	err := r.db.QueryRow(ctx, query, id).Scan(
		&w.ID, &w.UserID, &w.Name, &w.Balance, &w.CreatedAt, &w.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *walletRepository) Update(ctx context.Context, id int64, name string) error {
	query := `UPDATE wallets SET name = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, name, time.Now(), id)
	return err
}

func (r *walletRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM wallets WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *walletRepository) CheckOwnership(ctx context.Context, walletID int64, userID uuid.UUID) (*models.Wallet, error) {
	query := `SELECT id, user_id, name, balance, created_at, updated_at FROM wallets 
	          WHERE id = $1 AND user_id = $2`
	var w models.Wallet

	err := r.db.QueryRow(ctx, query, walletID, userID).Scan(
		&w.ID, &w.UserID, &w.Name, &w.Balance, &w.CreatedAt, &w.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &w, nil
}
