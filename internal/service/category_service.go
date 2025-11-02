package service

import (
	"context"
	"errors"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/Udean777/uang-bijak-go/internal/repository"
	"github.com/google/uuid"
)

var ErrForbidden = errors.New("forbidden access")

type CategoryService interface {
	CreateCategory(ctx context.Context, req models.UpsertCategoryRequest, userID uuid.UUID) (*models.Category, error)
	GetUserCategories(ctx context.Context, userID uuid.UUID) ([]models.Category, error)
	UpdateCategory(ctx context.Context, categoryID int64, req models.UpsertCategoryRequest, userID uuid.UUID) error
	DeleteCategory(ctx context.Context, categoryID int64, userID uuid.UUID) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, req models.UpsertCategoryRequest, userID uuid.UUID) (*models.Category, error) {
	cat := &models.Category{
		UserID: userID,
		Name:   req.Name,
	}

	_, err := s.categoryRepo.Create(ctx, cat)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) GetUserCategories(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	return s.categoryRepo.GetAllByUserID(ctx, userID)
}

func (s *categoryService) checkOwnership(ctx context.Context, categoryID int64, userID uuid.UUID) error {
	_, err := s.categoryRepo.CheckOwnership(ctx, categoryID, userID)
	if err != nil {
		return ErrForbidden
	}
	return nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, categoryID int64, req models.UpsertCategoryRequest, userID uuid.UUID) error {
	if err := s.checkOwnership(ctx, categoryID, userID); err != nil {
		return err
	}

	return s.categoryRepo.Update(ctx, categoryID, req.Name)
}

func (s *categoryService) DeleteCategory(ctx context.Context, categoryID int64, userID uuid.UUID) error {
	if err := s.checkOwnership(ctx, categoryID, userID); err != nil {
		return err
	}

	return s.categoryRepo.Delete(ctx, categoryID)
}
