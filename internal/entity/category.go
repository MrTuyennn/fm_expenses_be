package entity

import "time"

const (
	CategoryTypeIncome  = "income"
	CategoryTypeExpense = "expense"
)

// Category is the domain entity for an income/expense category.
type Category struct {
	ID        string
	UserID    string
	Name      string
	Icon      string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
