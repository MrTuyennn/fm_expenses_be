// Package category implements the Category usecase.
package category

import (
	"context"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/repo"
	"fm_expenses_be/internal/usecase"
)

// UseCase -.
type UseCase struct {
	repo repo.CategoryRepo
}

// New returns the Category usecase.
func New(r repo.CategoryRepo) usecase.Category {
	return &UseCase{repo: r}
}

func isValidType(t string) bool {
	return t == entity.CategoryTypeIncome || t == entity.CategoryTypeExpense
}

// Create -.
func (uc *UseCase) Create(ctx context.Context, userID, name, icon, catType string) (entity.Category, error) {
	if !isValidType(catType) {
		return entity.Category{}, entity.ErrInvalidCategoryType
	}

	c := &entity.Category{UserID: userID, Name: name, Icon: icon, Type: catType}
	if err := uc.repo.Create(ctx, c); err != nil {
		return entity.Category{}, err
	}
	return *c, nil
}

// List -.
func (uc *UseCase) List(ctx context.Context, userID string) ([]entity.Category, error) {
	return uc.repo.ListByUser(ctx, userID)
}

// Get -.
func (uc *UseCase) Get(ctx context.Context, id, userID string) (entity.Category, error) {
	return uc.repo.FindByID(ctx, id, userID)
}

// Update -.
func (uc *UseCase) Update(ctx context.Context, id, userID, name, icon, catType string) (entity.Category, error) {
	if !isValidType(catType) {
		return entity.Category{}, entity.ErrInvalidCategoryType
	}

	existing, err := uc.repo.FindByID(ctx, id, userID)
	if err != nil {
		return entity.Category{}, entity.ErrNotFound
	}
	existing.Name = name
	existing.Icon = icon
	existing.Type = catType
	if err := uc.repo.Update(ctx, &existing); err != nil {
		return entity.Category{}, err
	}
	return existing, nil
}

// Delete -.
func (uc *UseCase) Delete(ctx context.Context, id, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
