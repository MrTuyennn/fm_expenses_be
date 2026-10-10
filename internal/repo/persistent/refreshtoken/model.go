package refreshtoken

import (
	"time"

	"fm_expenses_be/internal/entity"
)

// Model is the GORM-tagged persistence row for a refresh token.
type Model struct {
	ID        string     `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string     `gorm:"column:user_id;type:uuid;index;not null"`
	TokenHash string     `gorm:"column:token_hash;uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

func (Model) TableName() string { return "refresh_tokens" }

func (m Model) toEntity() entity.RefreshToken {
	return entity.RefreshToken{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		RevokedAt: m.RevokedAt,
		CreatedAt: m.CreatedAt,
	}
}

func modelFromEntity(rt *entity.RefreshToken) *Model {
	return &Model{
		ID:        rt.ID,
		UserID:    rt.UserID,
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
		RevokedAt: rt.RevokedAt,
		CreatedAt: rt.CreatedAt,
	}
}
