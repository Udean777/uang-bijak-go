package repository

import (
	"context"
	"time"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) (int64, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Category, error)
	GetByID(ctx context.Context, id int64) (*models.Category, error)
	Update(ctx context.Context, id int64, name string) error
	Delete(ctx context.Context, id int64) error

	// Helper untuk mengecek kepemilikan
	CheckOwnership(ctx context.Context, categoryID int64, userID uuid.UUID) (*models.Category, error)
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) (int64, error) {
	query := `INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, category.UserID, category.Name).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return 0, err
	}
	return category.ID, nil
}

func (r *categoryRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE user_id = $1 ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*models.Category, error) {
	query := `SELECT id, user_id, name, created_at, updated_at FROM categories WHERE id = $1`
	var cat models.Category

	err := r.db.QueryRow(ctx, query, id).Scan(
		&cat.ID,
		&cat.UserID,
		&cat.Name,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) Update(ctx context.Context, id int64, name string) error {
	query := `UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, name, time.Now(), id)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *categoryRepository) CheckOwnership(ctx context.Context, categoryID int64, userID uuid.UUID) (*models.Category, error) {
	query := `SELECT id, user_id, name, created_at, updated_at FROM categories WHERE id = $1 AND user_id = $2`
	var cat models.Category

	err := r.db.QueryRow(ctx, query, categoryID, userID).Scan(
		&cat.ID,
		&cat.UserID,
		&cat.Name,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &cat, nil
}
