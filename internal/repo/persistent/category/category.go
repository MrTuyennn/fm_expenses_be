// Package category implements the Postgres-backed Category repository.
package category

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

// New returns the Category repository backed by Postgres/GORM.
func New(db *gorm.DB) repo.CategoryRepo {
	return &Repo{db: db}
}

// Create -.
func (r *Repo) Create(ctx context.Context, c *entity.Category) error {
	m := modelFromEntity(c)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	*c = m.toEntity()
	return nil
}

// ListByUser -.
func (r *Repo) ListByUser(ctx context.Context, userID string) ([]entity.Category, error) {
	var rows []Model
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entity.Category, 0, len(rows))
	for _, m := range rows {
		out = append(out, m.toEntity())
	}
	return out, nil
}

// FindByID -.
func (r *Repo) FindByID(ctx context.Context, id, userID string) (entity.Category, error) {
	var m Model
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&m).Error; err != nil {
		return entity.Category{}, err
	}
	return m.toEntity(), nil
}

// Update -.
func (r *Repo) Update(ctx context.Context, c *entity.Category) error {
	m := modelFromEntity(c)
	return r.db.WithContext(ctx).Save(m).Error
}

// Delete -.
func (r *Repo) Delete(ctx context.Context, id, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Model{}).Error
}
