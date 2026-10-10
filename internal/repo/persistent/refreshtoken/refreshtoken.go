// Package refreshtoken implements the Postgres-backed RefreshToken repository.
package refreshtoken

import (
	"context"
	"time"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/repo"

	"gorm.io/gorm"
)

// Repo -.
type Repo struct {
	db *gorm.DB
}

// New returns the RefreshToken repository backed by Postgres/GORM.
func New(db *gorm.DB) repo.RefreshTokenRepo {
	return &Repo{db: db}
}

// Create -.
func (r *Repo) Create(ctx context.Context, rt *entity.RefreshToken) error {
	m := modelFromEntity(rt)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	*rt = m.toEntity()
	return nil
}

// FindValidByHash returns entity.ErrRefreshTokenInvalid if no row matches,
// or the matching row is already revoked or past its expiry.
func (r *Repo) FindValidByHash(ctx context.Context, tokenHash string) (entity.RefreshToken, error) {
	var m Model
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&m).Error
	if err != nil {
		return entity.RefreshToken{}, entity.ErrRefreshTokenInvalid
	}
	return m.toEntity(), nil
}

// Revoke -.
func (r *Repo) Revoke(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Update("revoked_at", time.Now()).Error
}
