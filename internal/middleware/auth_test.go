package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fm_expenses_be/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func signToken(t *testing.T, secret, sub string, expiresAt time.Time) string {
	t.Helper()

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return signed
}

func newProtectedHandler() (http.Handler, *string) {
	var gotUserID string
	h := middleware.NewAuth(testSecret).Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID = middleware.UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))
	return h, &gotUserID
}

func TestAuthMiddleware_MissingHeader_Returns401(t *testing.T) {
	h, _ := newProtectedHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_NonBearerHeader_Returns401(t *testing.T) {
	h, _ := newProtectedHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic sometoken")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_WrongSecret_Returns401(t *testing.T) {
	h, _ := newProtectedHandler()

	token := signToken(t, "a-different-secret", "user-1", time.Now().Add(time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ExpiredToken_Returns401(t *testing.T) {
	h, _ := newProtectedHandler()

	token := signToken(t, testSecret, "user-1", time.Now().Add(-time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ValidToken_InjectsUserID(t *testing.T) {
	h, gotUserID := newProtectedHandler()

	token := signToken(t, testSecret, "user-123", time.Now().Add(time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if *gotUserID != "user-123" {
		t.Fatalf("expected userID user-123, got %q", *gotUserID)
	}
}
