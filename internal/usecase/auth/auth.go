// Package auth implements the Auth usecase: registration, login, and
// refresh-token rotation. Access tokens are short-lived, stateless JWTs
// (pkg/jwt). Refresh tokens are opaque random strings whose SHA-256 hash is
// stored in repo.RefreshTokenRepo, so they can be revoked and rotated —
// something a stateless JWT alone cannot do.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/repo"
	"fm_expenses_be/internal/usecase"
	"fm_expenses_be/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength = 8

// UseCase -.
type UseCase struct {
	userRepo         repo.UserRepo
	refreshTokenRepo repo.RefreshTokenRepo
	jwtManager       *jwt.Manager
	refreshTokenTTL  time.Duration
}

// New returns the Auth usecase.
func New(userRepo repo.UserRepo, refreshTokenRepo repo.RefreshTokenRepo, jwtManager *jwt.Manager, refreshTokenTTL time.Duration) usecase.Auth {
	return &UseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		refreshTokenTTL:  refreshTokenTTL,
	}
}

// Register -.
func (uc *UseCase) Register(ctx context.Context, email, password, fullName string) (entity.User, entity.TokenPair, error) {
	if len(password) < minPasswordLength {
		return entity.User{}, entity.TokenPair{}, entity.ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, entity.TokenPair{}, fmt.Errorf("auth - Register - bcrypt.GenerateFromPassword: %w", err)
	}

	u := &entity.User{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
	}
	if err := uc.userRepo.Create(ctx, u); err != nil {
		if errors.Is(err, entity.ErrUserAlreadyExists) {
			return entity.User{}, entity.TokenPair{}, entity.ErrUserAlreadyExists
		}
		return entity.User{}, entity.TokenPair{}, fmt.Errorf("auth - Register - userRepo.Create: %w", err)
	}

	tokens, err := uc.issueTokenPair(ctx, u.ID)
	if err != nil {
		return entity.User{}, entity.TokenPair{}, err
	}
	return *u, tokens, nil
}

// Login -.
func (uc *UseCase) Login(ctx context.Context, email, password string) (entity.TokenPair, error) {
	u, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return entity.TokenPair{}, entity.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return entity.TokenPair{}, entity.ErrInvalidCredentials
	}

	return uc.issueTokenPair(ctx, u.ID)
}

// Refresh validates and rotates a refresh token: the old one is revoked and
// a brand new access/refresh pair is issued. Rotation means a leaked,
// already-used refresh token is worthless to an attacker.
func (uc *UseCase) Refresh(ctx context.Context, refreshToken string) (entity.TokenPair, error) {
	hash := hashToken(refreshToken)

	rt, err := uc.refreshTokenRepo.FindValidByHash(ctx, hash)
	if err != nil {
		return entity.TokenPair{}, entity.ErrRefreshTokenInvalid
	}

	if err := uc.refreshTokenRepo.Revoke(ctx, rt.ID); err != nil {
		return entity.TokenPair{}, fmt.Errorf("auth - Refresh - refreshTokenRepo.Revoke: %w", err)
	}

	return uc.issueTokenPair(ctx, rt.UserID)
}

// Logout revokes a refresh token. It is deliberately silent about whether
// the token existed, so a client can always treat logout as succeeding.
func (uc *UseCase) Logout(ctx context.Context, refreshToken string) error {
	hash := hashToken(refreshToken)

	rt, err := uc.refreshTokenRepo.FindValidByHash(ctx, hash)
	if err != nil {
		return nil
	}
	return uc.refreshTokenRepo.Revoke(ctx, rt.ID)
}

func (uc *UseCase) issueTokenPair(ctx context.Context, userID string) (entity.TokenPair, error) {
	accessToken, err := uc.jwtManager.GenerateToken(userID)
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("auth - issueTokenPair - jwtManager.GenerateToken: %w", err)
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("auth - issueTokenPair - generateRefreshToken: %w", err)
	}

	rt := &entity.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: time.Now().Add(uc.refreshTokenTTL),
	}
	if err := uc.refreshTokenRepo.Create(ctx, rt); err != nil {
		return entity.TokenPair{}, fmt.Errorf("auth - issueTokenPair - refreshTokenRepo.Create: %w", err)
	}

	return entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    uc.jwtManager.ExpiresIn(),
	}, nil
}

// generateRefreshToken returns a 256-bit cryptographically random,
// URL-safe string — opaque to the client, unrelated to the access token.
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// hashToken returns the hex-encoded SHA-256 hash of a refresh token. Only
// this hash is ever persisted, so a leaked database dump doesn't hand out
// usable tokens.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
