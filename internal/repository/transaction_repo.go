package repository

import (
	"context"
	"time"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	CreateTx(ctx context.Context, tx pgx.Tx, transaction *models.Transaction) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
	// TODO: tambahkan GetByID, Update, Delete
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTx(ctx context.Context, tx pgx.Tx, t *models.Transaction) error {
	query := `INSERT INTO transactions 
	          (user_id, wallet_id, category_id, amount, type, description, transaction_date)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          RETURNING id, created_at, updated_at`

	if t.TransactionDate.IsZero() {
		t.TransactionDate = time.Now()
	}

	return tx.QueryRow(ctx, query,
		t.UserID, t.WalletID, t.CategoryID, t.Amount, t.Type, t.Description, t.TransactionDate,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *transactionRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	query := `SELECT id, wallet_id, category_id, amount, type, description, transaction_date, created_at, updated_at 
	          FROM transactions 
	          WHERE user_id = $1 
	          ORDER BY transaction_date DESC, created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID, &t.WalletID, &t.CategoryID, &t.Amount, &t.Type,
			&t.Description, &t.TransactionDate, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
