// Package jwt issues and verifies short-lived access tokens signed by this
// backend (HS256). Refresh tokens are handled separately — see
// internal/usecase/auth — because they must be revocable, which a stateless
// JWT cannot be.
package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// ErrUnexpectedSigningMethod is returned when a token's signing method
// doesn't match what this Manager expects.
var ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

// Manager issues and parses access tokens.
type Manager struct {
	secret   string
	duration time.Duration
}

// New returns a Manager that signs tokens with secret and sets them to
// expire after duration.
func New(secret string, duration time.Duration) *Manager {
	return &Manager{secret: secret, duration: duration}
}

// GenerateToken creates a signed access token for the given user ID.
func (m *Manager) GenerateToken(userID string) (string, error) {
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(m.duration)),
	})

	signed, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("jwt - GenerateToken - token.SignedString: %w", err)
	}
	return signed, nil
}

// ExpiresIn returns the access token lifetime in seconds, for clients that
// need to know when to proactively refresh.
func (m *Manager) ExpiresIn() int {
	return int(m.duration.Seconds())
}
