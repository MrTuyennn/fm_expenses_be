// Package expense implements the Postgres-backed Expense repository.
package expense

import (
	"context"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/repo"

	"gorm.io/gorm"
)

// Repo -.
type Repo struct {
	db *gorm.DB
}

// New returns the Expense repository backed by Postgres/GORM.
func New(db *gorm.DB) repo.ExpenseRepo {
	return &Repo{db: db}
}

// Create -.
func (r *Repo) Create(ctx context.Context, e *entity.Expense) error {
	m := modelFromEntity(e)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	*e = m.toEntity()
	return nil
}

// ListByUser -.
func (r *Repo) ListByUser(ctx context.Context, userID string, filter entity.ExpenseFilter) ([]entity.Expense, error) {
	q := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filter.CategoryID != "" {
		q = q.Where("category_id = ?", filter.CategoryID)
	}
	if filter.From != nil {
		q = q.Where("spent_at >= ?", *filter.From)
	}
	if filter.To != nil {
		q = q.Where("spent_at <= ?", *filter.To)
	}

	var rows []Model
	if err := q.Order("spent_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entity.Expense, 0, len(rows))
	for _, m := range rows {
		out = append(out, m.toEntity())
	}
	return out, nil
}

// FindByID -.
func (r *Repo) FindByID(ctx context.Context, id, userID string) (entity.Expense, error) {
	var m Model
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&m).Error; err != nil {
		return entity.Expense{}, err
	}
	return m.toEntity(), nil
}

// Update -.
func (r *Repo) Update(ctx context.Context, e *entity.Expense) error {
	m := modelFromEntity(e)
	return r.db.WithContext(ctx).Save(m).Error
}

// Delete -.
func (r *Repo) Delete(ctx context.Context, id, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Model{}).Error
}
