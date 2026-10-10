package user

import (
	"time"

	"fm_expenses_be/internal/entity"
)

// Model is the GORM-tagged persistence row for an account, kept separate
// from entity.User so the domain entity never carries ORM tags.
type Model struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	FullName     string    `gorm:"column:full_name"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (Model) TableName() string { return "users" }

func (m Model) toEntity() entity.User {
	return entity.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		FullName:     m.FullName,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func modelFromEntity(u *entity.User) *Model {
	return &Model{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     u.FullName,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
