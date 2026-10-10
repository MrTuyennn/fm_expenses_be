package request

// Expense is the request body for creating/updating an expense.
type Expense struct {
	CategoryID string  `json:"category_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Note       string  `json:"note"`
	SpentAt    string  `json:"spent_at"`
}
