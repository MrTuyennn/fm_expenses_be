package entity

import "time"

// RefreshToken is a long-lived, revocable credential used to mint new access
// tokens. Only its SHA-256 hash is ever persisted — the plaintext value is
// handed to the client once, at issuance, and never stored.
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// TokenPair is what Register/Login/Refresh hand back to the client.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // access token lifetime, in seconds
}
