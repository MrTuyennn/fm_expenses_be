package auth_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/usecase"
	"fm_expenses_be/internal/usecase/auth"
	"fm_expenses_be/pkg/jwt"
)

type mockUserRepo struct {
	byID    map[string]entity.User
	byEmail map[string]string
	nextID  int
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{byID: map[string]entity.User{}, byEmail: map[string]string{}}
}

func (m *mockUserRepo) Create(_ context.Context, u *entity.User) error {
	if _, ok := m.byEmail[u.Email]; ok {
		return entity.ErrUserAlreadyExists
	}
	m.nextID++
	u.ID = "user-" + strconv.Itoa(m.nextID)
	m.byID[u.ID] = *u
	m.byEmail[u.Email] = u.ID
	return nil
}

func (m *mockUserRepo) FindByID(_ context.Context, id string) (entity.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return entity.User{}, errors.New("record not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindByEmail(_ context.Context, email string) (entity.User, error) {
	id, ok := m.byEmail[email]
	if !ok {
		return entity.User{}, errors.New("record not found")
	}
	return m.byID[id], nil
}

type mockRefreshTokenRepo struct {
	byHash map[string]*entity.RefreshToken
	nextID int
}

func newMockRefreshTokenRepo() *mockRefreshTokenRepo {
	return &mockRefreshTokenRepo{byHash: map[string]*entity.RefreshToken{}}
}

func (m *mockRefreshTokenRepo) Create(_ context.Context, rt *entity.RefreshToken) error {
	m.nextID++
	rt.ID = "rt-" + strconv.Itoa(m.nextID)
	cp := *rt
	m.byHash[rt.TokenHash] = &cp
	return nil
}

func (m *mockRefreshTokenRepo) FindValidByHash(_ context.Context, tokenHash string) (entity.RefreshToken, error) {
	rt, ok := m.byHash[tokenHash]
	if !ok || rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return entity.RefreshToken{}, entity.ErrRefreshTokenInvalid
	}
	return *rt, nil
}

func (m *mockRefreshTokenRepo) Revoke(_ context.Context, id string) error {
	for _, rt := range m.byHash {
		if rt.ID == id {
			now := time.Now()
			rt.RevokedAt = &now
			return nil
		}
	}
	return errors.New("record not found")
}

func newTestUseCase() usecase.Auth {
	userRepo := newMockUserRepo()
	refreshRepo := newMockRefreshTokenRepo()
	jwtManager := jwt.New("test-secret", 15*time.Minute)
	return auth.New(userRepo, refreshRepo, jwtManager, 30*24*time.Hour)
}

func TestAuthUsecase_Register_RejectsShortPassword(t *testing.T) {
	uc := newTestUseCase()

	_, _, err := uc.Register(context.Background(), "a@example.com", "short", "A")
	if !errors.Is(err, entity.ErrPasswordTooShort) {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestAuthUsecase_Register_RejectsDuplicateEmail(t *testing.T) {
	uc := newTestUseCase()

	_, _, err := uc.Register(context.Background(), "a@example.com", "password1", "A")
	if err != nil {
		t.Fatalf("setup: unexpected error: %v", err)
	}

	_, _, err = uc.Register(context.Background(), "a@example.com", "password2", "B")
	if !errors.Is(err, entity.ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthUsecase_Register_ReturnsUsableTokenPair(t *testing.T) {
	uc := newTestUseCase()

	_, tokens, err := uc.Register(context.Background(), "a@example.com", "password1", "A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens, got %+v", tokens)
	}
}

func TestAuthUsecase_Login_RejectsWrongPassword(t *testing.T) {
	uc := newTestUseCase()

	_, _, err := uc.Register(context.Background(), "a@example.com", "password1", "A")
	if err != nil {
		t.Fatalf("setup: unexpected error: %v", err)
	}

	_, err = uc.Login(context.Background(), "a@example.com", "wrong-password")
	if !errors.Is(err, entity.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUsecase_Login_RejectsUnknownEmail(t *testing.T) {
	uc := newTestUseCase()

	_, err := uc.Login(context.Background(), "nobody@example.com", "password1")
	if !errors.Is(err, entity.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUsecase_Login_Succeeds(t *testing.T) {
	uc := newTestUseCase()

	_, _, err := uc.Register(context.Background(), "a@example.com", "password1", "A")
	if err != nil {
		t.Fatalf("setup: unexpected error: %v", err)
	}

	tokens, err := uc.Login(context.Background(), "a@example.com", "password1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens, got %+v", tokens)
	}
}

func TestAuthUsecase_Refresh_RejectsUnknownToken(t *testing.T) {
	uc := newTestUseCase()

	_, err := uc.Refresh(context.Background(), "not-a-real-refresh-token")
	if !errors.Is(err, entity.ErrRefreshTokenInvalid) {
		t.Fatalf("expected ErrRefreshTokenInvalid, got %v", err)
	}
}

func TestAuthUsecase_Refresh_RotatesToken(t *testing.T) {
	uc := newTestUseCase()

	_, tokens, err := uc.Register(context.Background(), "a@example.com", "password1", "A")
	if err != nil {
		t.Fatalf("setup: unexpected error: %v", err)
	}

	rotated, err := uc.Refresh(context.Background(), tokens.RefreshToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated.RefreshToken == tokens.RefreshToken {
		t.Fatalf("expected a new refresh token, got the same one back")
	}

	// The old refresh token must now be rejected (rotation = single use).
	if _, err := uc.Refresh(context.Background(), tokens.RefreshToken); !errors.Is(err, entity.ErrRefreshTokenInvalid) {
		t.Fatalf("expected old refresh token to be invalid after rotation, got %v", err)
	}
}

func TestAuthUsecase_Logout_RevokesToken(t *testing.T) {
	uc := newTestUseCase()

	_, tokens, err := uc.Register(context.Background(), "a@example.com", "password1", "A")
	if err != nil {
		t.Fatalf("setup: unexpected error: %v", err)
	}

	if err := uc.Logout(context.Background(), tokens.RefreshToken); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := uc.Refresh(context.Background(), tokens.RefreshToken); !errors.Is(err, entity.ErrRefreshTokenInvalid) {
		t.Fatalf("expected refresh token to be invalid after logout, got %v", err)
	}
}

func TestAuthUsecase_Logout_UnknownTokenIsSilentlyOK(t *testing.T) {
	uc := newTestUseCase()

	if err := uc.Logout(context.Background(), "never-issued"); err != nil {
		t.Fatalf("expected no error for unknown token, got %v", err)
	}
}
