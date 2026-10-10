// Package response holds v1 API response DTOs, kept separate from
// internal/entity so the domain entity never carries JSON tags.
package response

import (
	"time"

	"fm_expenses_be/internal/entity"
)

// User is the wire format for a profile.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromUser maps a domain entity to its wire format.
func FromUser(u entity.User) User {
	return User{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
