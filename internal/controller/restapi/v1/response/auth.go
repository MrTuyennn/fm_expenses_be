package response

import "fm_expenses_be/internal/entity"

// TokenPair is the wire format for an access/refresh token pair.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// FromTokenPair maps a domain entity to its wire format.
func FromTokenPair(t entity.TokenPair) TokenPair {
	return TokenPair{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresIn:    t.ExpiresIn,
		TokenType:    "Bearer",
	}
}

// RegisterResult is the response body for POST /auth/register: the new
// user plus a token pair so the client is logged in immediately.
type RegisterResult struct {
	User   User      `json:"user"`
	Tokens TokenPair `json:"tokens"`
}
