// Package expense implements the Expense usecase.
package expense

import (
	"context"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/repo"
	"fm_expenses_be/internal/usecase"
)

// UseCase -.
type UseCase struct {
	repo         repo.ExpenseRepo
	categoryRepo repo.CategoryRepo
}

// New returns the Expense usecase. It also depends on repo.CategoryRepo to
// validate that a category belongs to the requesting user.
func New(r repo.ExpenseRepo, categoryRepo repo.CategoryRepo) usecase.Expense {
	return &UseCase{repo: r, categoryRepo: categoryRepo}
}

// Create -.
func (uc *UseCase) Create(ctx context.Context, in entity.CreateExpenseInput) (entity.Expense, error) {
	if _, err := uc.categoryRepo.FindByID(ctx, in.CategoryID, in.UserID); err != nil {
		return entity.Expense{}, entity.ErrCategoryInvalid
	}

	e := &entity.Expense{
		UserID:     in.UserID,
		CategoryID: in.CategoryID,
		Amount:     in.Amount,
		Currency:   in.Currency,
		Note:       in.Note,
		SpentAt:    in.SpentAt,
	}
	if err := uc.repo.Create(ctx, e); err != nil {
		return entity.Expense{}, err
	}
	return *e, nil
}

// List -.
func (uc *UseCase) List(ctx context.Context, userID string, filter entity.ExpenseFilter) ([]entity.Expense, error) {
	return uc.repo.ListByUser(ctx, userID, filter)
}

// Get -.
func (uc *UseCase) Get(ctx context.Context, id, userID string) (entity.Expense, error) {
	return uc.repo.FindByID(ctx, id, userID)
}

// Update -.
func (uc *UseCase) Update(ctx context.Context, id, userID string, in entity.UpdateExpenseInput) (entity.Expense, error) {
	existing, err := uc.repo.FindByID(ctx, id, userID)
	if err != nil {
		return entity.Expense{}, entity.ErrNotFound
	}

	if _, err := uc.categoryRepo.FindByID(ctx, in.CategoryID, userID); err != nil {
		return entity.Expense{}, entity.ErrCategoryInvalid
	}

	existing.CategoryID = in.CategoryID
	existing.Amount = in.Amount
	existing.Currency = in.Currency
	existing.Note = in.Note
	existing.SpentAt = in.SpentAt

	if err := uc.repo.Update(ctx, &existing); err != nil {
		return entity.Expense{}, err
	}
	return existing, nil
}

// Delete -.
func (uc *UseCase) Delete(ctx context.Context, id, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
