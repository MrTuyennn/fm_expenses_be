// Package user implements the Postgres-backed User repository.
package user

import (
	"context"
	"errors"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/repo"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repo -.
type Repo struct {
	db *gorm.DB
}

// New returns the User repository backed by Postgres/GORM.
func New(db *gorm.DB) repo.UserRepo {
	return &Repo{db: db}
}

// Create -.
func (r *Repo) Create(ctx context.Context, u *entity.User) error {
	m := modelFromEntity(u)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.ErrUserAlreadyExists
		}
		return err
	}
	*u = m.toEntity()
	return nil
}

// FindByID -.
func (r *Repo) FindByID(ctx context.Context, id string) (entity.User, error) {
	var m Model
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return entity.User{}, err
	}
	return m.toEntity(), nil
}

// FindByEmail -.
func (r *Repo) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var m Model
	if err := r.db.WithContext(ctx).First(&m, "email = ?", email).Error; err != nil {
		return entity.User{}, err
	}
	return m.toEntity(), nil
}
