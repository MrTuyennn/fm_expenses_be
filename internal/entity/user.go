package entity

import "time"

// User is the domain entity for an account. Plain Go — no ORM or JSON tags.
// PasswordHash is never serialized to a response (see response.User, which
// deliberately omits it).
type User struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
