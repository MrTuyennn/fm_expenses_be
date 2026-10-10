package response

import (
	"time"

	"fm_expenses_be/internal/entity"
)

// Expense is the wire format for an expense.
type Expense struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CategoryID string    `json:"category_id"`
	Amount     float64   `json:"amount"`
	Currency   string    `json:"currency"`
	Note       string    `json:"note"`
	SpentAt    time.Time `json:"spent_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FromExpense maps a domain entity to its wire format.
func FromExpense(e entity.Expense) Expense {
	return Expense{
		ID:         e.ID,
		UserID:     e.UserID,
		CategoryID: e.CategoryID,
		Amount:     e.Amount,
		Currency:   e.Currency,
		Note:       e.Note,
		SpentAt:    e.SpentAt,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

// FromExpenseList maps a slice of domain entities to their wire format.
func FromExpenseList(items []entity.Expense) []Expense {
	out := make([]Expense, 0, len(items))
	for _, e := range items {
		out = append(out, FromExpense(e))
	}
	return out
}
