// Package usecase implements application business logic. Each logic group
// has its own subpackage (user, category, expense); this file only declares
// the interfaces the delivery layer (controller) depends on.
package usecase

import (
	"context"

	"fm_expenses_be/internal/entity"
)

type (
	// User -.
	User interface {
		GetProfile(ctx context.Context, id string) (entity.User, error)
	}

	// Auth -.
	Auth interface {
		Register(ctx context.Context, email, password, fullName string) (entity.User, entity.TokenPair, error)
		Login(ctx context.Context, email, password string) (entity.TokenPair, error)
		Refresh(ctx context.Context, refreshToken string) (entity.TokenPair, error)
		Logout(ctx context.Context, refreshToken string) error
	}

	// Category -.
	Category interface {
		Create(ctx context.Context, userID, name, icon, catType string) (entity.Category, error)
		List(ctx context.Context, userID string) ([]entity.Category, error)
		Get(ctx context.Context, id, userID string) (entity.Category, error)
		Update(ctx context.Context, id, userID, name, icon, catType string) (entity.Category, error)
		Delete(ctx context.Context, id, userID string) error
	}

	// Expense -.
	Expense interface {
		Create(ctx context.Context, in entity.CreateExpenseInput) (entity.Expense, error)
		List(ctx context.Context, userID string, filter entity.ExpenseFilter) ([]entity.Expense, error)
		Get(ctx context.Context, id, userID string) (entity.Expense, error)
		Update(ctx context.Context, id, userID string, in entity.UpdateExpenseInput) (entity.Expense, error)
		Delete(ctx context.Context, id, userID string) error
	}
)
