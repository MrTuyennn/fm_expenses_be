package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/usecase/user"
)

// mockUserRepo is an in-memory stand-in for repo.UserRepo.
type mockUserRepo struct {
	byID    map[string]entity.User
	byEmail map[string]string // email -> id
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{byID: map[string]entity.User{}, byEmail: map[string]string{}}
}

func (m *mockUserRepo) Create(_ context.Context, u *entity.User) error {
	if _, ok := m.byEmail[u.Email]; ok {
		return entity.ErrUserAlreadyExists
	}
	u.ID = "generated-id"
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

func TestUserUsecase_GetProfile_ReturnsStoredProfile(t *testing.T) {
	repo := newMockUserRepo()
	repo.byID["user-1"] = entity.User{
		ID:        "user-1",
		Email:     "a@example.com",
		FullName:  "A B",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	uc := user.New(repo)

	got, err := uc.GetProfile(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Email != "a@example.com" || got.FullName != "A B" {
		t.Fatalf("unexpected profile: %+v", got)
	}
}

func TestUserUsecase_GetProfile_NotFoundMapsToErrNotFound(t *testing.T) {
	uc := user.New(newMockUserRepo())

	_, err := uc.GetProfile(context.Background(), "missing-id")
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
