package category_test

import (
	"context"
	"errors"
	"testing"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/usecase/category"
)

// mockCategoryRepo is an in-memory stand-in for repo.CategoryRepo. It exists
// to prove the usecase layer is testable without a real Postgres connection.
type mockCategoryRepo struct {
	byID map[string]*entity.Category
}

func newMockCategoryRepo() *mockCategoryRepo {
	return &mockCategoryRepo{byID: map[string]*entity.Category{}}
}

func (m *mockCategoryRepo) Create(_ context.Context, c *entity.Category) error {
	c.ID = "generated-id"
	m.byID[c.ID] = c
	return nil
}

func (m *mockCategoryRepo) ListByUser(_ context.Context, userID string) ([]entity.Category, error) {
	var out []entity.Category
	for _, c := range m.byID {
		if c.UserID == userID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (m *mockCategoryRepo) FindByID(_ context.Context, id, userID string) (entity.Category, error) {
	c, ok := m.byID[id]
	if !ok || c.UserID != userID {
		return entity.Category{}, errors.New("record not found")
	}
	return *c, nil
}

func (m *mockCategoryRepo) Update(_ context.Context, c *entity.Category) error {
	if _, ok := m.byID[c.ID]; !ok {
		return errors.New("record not found")
	}
	m.byID[c.ID] = c
	return nil
}

func (m *mockCategoryRepo) Delete(_ context.Context, id, _ string) error {
	delete(m.byID, id)
	return nil
}

func TestCategoryUsecase_Create_RejectsInvalidType(t *testing.T) {
	uc := category.New(newMockCategoryRepo())

	_, err := uc.Create(context.Background(), "user-1", "Ăn uống", "", "not-a-real-type")
	if !errors.Is(err, entity.ErrInvalidCategoryType) {
		t.Fatalf("expected ErrInvalidCategoryType, got %v", err)
	}
}

func TestCategoryUsecase_Create_Succeeds(t *testing.T) {
	uc := category.New(newMockCategoryRepo())

	c, err := uc.Create(context.Background(), "user-1", "Ăn uống", "🍜", entity.CategoryTypeExpense)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.UserID != "user-1" || c.Type != entity.CategoryTypeExpense {
		t.Fatalf("unexpected category: %+v", c)
	}
}

func TestCategoryUsecase_Update_NotFoundMapsToErrNotFound(t *testing.T) {
	uc := category.New(newMockCategoryRepo())

	_, err := uc.Update(context.Background(), "missing-id", "user-1", "x", "", entity.CategoryTypeIncome)
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCategoryUsecase_Update_OtherUsersCategoryIsNotFound(t *testing.T) {
	repo := newMockCategoryRepo()
	uc := category.New(repo)

	c, err := uc.Create(context.Background(), "user-1", "Lương", "", entity.CategoryTypeIncome)
	if err != nil {
		t.Fatalf("setup: unexpected error: %v", err)
	}

	_, err = uc.Update(context.Background(), c.ID, "user-2", "x", "", entity.CategoryTypeIncome)
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for another user's category, got %v", err)
	}
}
