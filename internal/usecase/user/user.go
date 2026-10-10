// Package user implements the User usecase.
package user

import (
	"context"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/repo"
	"fm_expenses_be/internal/usecase"
)

// UseCase -.
type UseCase struct {
	repo repo.UserRepo
}

// New returns the User usecase.
func New(r repo.UserRepo) usecase.User {
	return &UseCase{repo: r}
}

// GetProfile -.
func (uc *UseCase) GetProfile(ctx context.Context, id string) (entity.User, error) {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, entity.ErrNotFound
	}
	return u, nil
}
