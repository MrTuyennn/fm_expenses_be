// Package repo implements application outer layer logic (persistence). Each
// logic group has its own subpackage under repo/persistent; this file only
// declares the interfaces the usecase layer depends on.
package repo

import (
	"context"

	"fm_expenses_be/internal/entity"
)

type (
	// UserRepo -.
	UserRepo interface {
		// Create returns entity.ErrUserAlreadyExists if the email is taken.
		Create(ctx context.Context, u *entity.User) error
		FindByID(ctx context.Context, id string) (entity.User, error)
		FindByEmail(ctx context.Context, email string) (entity.User, error)
	}

	// RefreshTokenRepo -.
	RefreshTokenRepo interface {
		Create(ctx context.Context, rt *entity.RefreshToken) error
		// FindValidByHash returns entity.ErrRefreshTokenInvalid if no
		// matching token exists, or it is expired/revoked.
		FindValidByHash(ctx context.Context, tokenHash string) (entity.RefreshToken, error)
		Revoke(ctx context.Context, id string) error
	}

	// CategoryRepo -.
	CategoryRepo interface {
		Create(ctx context.Context, c *entity.Category) error
		ListByUser(ctx context.Context, userID string) ([]entity.Category, error)
		FindByID(ctx context.Context, id, userID string) (entity.Category, error)
		Update(ctx context.Context, c *entity.Category) error
		Delete(ctx context.Context, id, userID string) error
	}

	// ExpenseRepo -.
	ExpenseRepo interface {
		Create(ctx context.Context, e *entity.Expense) error
		ListByUser(ctx context.Context, userID string, filter entity.ExpenseFilter) ([]entity.Expense, error)
		FindByID(ctx context.Context, id, userID string) (entity.Expense, error)
		Update(ctx context.Context, e *entity.Expense) error
		Delete(ctx context.Context, id, userID string) error
	}
)
