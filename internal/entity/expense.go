package entity

import "time"

// Expense is the domain entity for a spending/income transaction.
type Expense struct {
	ID         string
	UserID     string
	CategoryID string
	Amount     float64
	Currency   string
	Note       string
	SpentAt    time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ExpenseFilter narrows an Expense List query.
type ExpenseFilter struct {
	CategoryID string
	From       *time.Time
	To         *time.Time
}

// CreateExpenseInput carries the fields needed to create an Expense.
type CreateExpenseInput struct {
	UserID     string
	CategoryID string
	Amount     float64
	Currency   string
	Note       string
	SpentAt    time.Time
}

// UpdateExpenseInput carries the fields needed to replace an Expense.
type UpdateExpenseInput struct {
	CategoryID string
	Amount     float64
	Currency   string
	Note       string
	SpentAt    time.Time
}
