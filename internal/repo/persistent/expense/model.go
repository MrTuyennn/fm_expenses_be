package expense

import (
	"time"

	"fm_expenses_be/internal/entity"
)

// Model is the GORM-tagged persistence row for an expense.
type Model struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     string    `gorm:"column:user_id;type:uuid;index;not null"`
	CategoryID string    `gorm:"column:category_id;type:uuid;index;not null"`
	Amount     float64   `gorm:"column:amount;type:numeric(14,2);not null"`
	Currency   string    `gorm:"column:currency;size:3;not null;default:'VND'"`
	Note       string    `gorm:"column:note"`
	SpentAt    time.Time `gorm:"column:spent_at;not null"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Model) TableName() string { return "expenses" }

func (m Model) toEntity() entity.Expense {
	return entity.Expense{
		ID:         m.ID,
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		Amount:     m.Amount,
		Currency:   m.Currency,
		Note:       m.Note,
		SpentAt:    m.SpentAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func modelFromEntity(e *entity.Expense) *Model {
	return &Model{
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
