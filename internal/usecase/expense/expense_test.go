package expense_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/usecase/expense"
)

// mockExpenseRepo is an in-memory stand-in for repo.ExpenseRepo.
type mockExpenseRepo struct {
	byID map[string]*entity.Expense
}

func newMockExpenseRepo() *mockExpenseRepo {
	return &mockExpenseRepo{byID: map[string]*entity.Expense{}}
}

func (m *mockExpenseRepo) Create(_ context.Context, e *entity.Expense) error {
	e.ID = "generated-expense-id"
	m.byID[e.ID] = e
	return nil
}

func (m *mockExpenseRepo) ListByUser(_ context.Context, userID string, _ entity.ExpenseFilter) ([]entity.Expense, error) {
	var out []entity.Expense
	for _, e := range m.byID {
		if e.UserID == userID {
			out = append(out, *e)
		}
	}
	return out, nil
}

func (m *mockExpenseRepo) FindByID(_ context.Context, id, userID string) (entity.Expense, error) {
	e, ok := m.byID[id]
	if !ok || e.UserID != userID {
		return entity.Expense{}, errors.New("record not found")
	}
	return *e, nil
}

func (m *mockExpenseRepo) Update(_ context.Context, e *entity.Expense) error {
	if _, ok := m.byID[e.ID]; !ok {
		return errors.New("record not found")
	}
	m.byID[e.ID] = e
	return nil
}

func (m *mockExpenseRepo) Delete(_ context.Context, id, _ string) error {
	delete(m.byID, id)
	return nil
}

// stubCategoryRepo is a minimal repo.CategoryRepo stand-in: it only tracks
// which user owns which category id, which is all the Expense usecase needs.
type stubCategoryRepo struct {
	owners map[string]string // categoryID -> userID
}

func (s *stubCategoryRepo) FindByID(_ context.Context, id, userID string) (entity.Category, error) {
	owner, ok := s.owners[id]
	if !ok || owner != userID {
		return entity.Category{}, errors.New("record not found")
	}
	return entity.Category{ID: id, UserID: userID}, nil
}

func (s *stubCategoryRepo) Create(context.Context, *entity.Category) error { return nil }
func (s *stubCategoryRepo) ListByUser(context.Context, string) ([]entity.Category, error) {
	return nil, nil
}
func (s *stubCategoryRepo) Update(context.Context, *entity.Category) error { return nil }
func (s *stubCategoryRepo) Delete(context.Context, string, string) error   { return nil }

func TestExpenseUsecase_Create_RejectsCategoryFromAnotherUser(t *testing.T) {
	categoryRepo := &stubCategoryRepo{owners: map[string]string{"cat-1": "user-2"}}
	uc := expense.New(newMockExpenseRepo(), categoryRepo)

	_, err := uc.Create(context.Background(), entity.CreateExpenseInput{
		UserID:     "user-1",
		CategoryID: "cat-1",
		Amount:     50000,
		Currency:   "VND",
		SpentAt:    time.Now(),
	})
	if !errors.Is(err, entity.ErrCategoryInvalid) {
		t.Fatalf("expected ErrCategoryInvalid, got %v", err)
	}
}

func TestExpenseUsecase_Create_Succeeds(t *testing.T) {
	categoryRepo := &stubCategoryRepo{owners: map[string]string{"cat-1": "user-1"}}
	uc := expense.New(newMockExpenseRepo(), categoryRepo)

	spentAt := time.Now()
	e, err := uc.Create(context.Background(), entity.CreateExpenseInput{
		UserID:     "user-1",
		CategoryID: "cat-1",
		Amount:     50000,
		Currency:   "VND",
		Note:       "Phở",
		SpentAt:    spentAt,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.UserID != "user-1" || e.CategoryID != "cat-1" || e.Amount != 50000 {
		t.Fatalf("unexpected expense: %+v", e)
	}
}

func TestExpenseUsecase_Update_NotFoundMapsToErrNotFound(t *testing.T) {
	categoryRepo := &stubCategoryRepo{owners: map[string]string{"cat-1": "user-1"}}
	uc := expense.New(newMockExpenseRepo(), categoryRepo)

	_, err := uc.Update(context.Background(), "missing-id", "user-1", entity.UpdateExpenseInput{
		CategoryID: "cat-1",
		Amount:     1000,
		SpentAt:    time.Now(),
	})
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestExpenseUsecase_Update_RejectsCategoryFromAnotherUser(t *testing.T) {
	categoryRepo := &stubCategoryRepo{owners: map[string]string{
		"cat-1": "user-1",
		"cat-2": "user-2",
	}}
	expenseRepo := newMockExpenseRepo()
	uc := expense.New(expenseRepo, categoryRepo)

	created, err := uc.Create(context.Background(), entity.CreateExpenseInput{
		UserID:     "user-1",
		CategoryID: "cat-1",
		Amount:     1000,
		SpentAt:    time.Now(),
	})
	if err != nil {
		t.Fatalf("setup: unexpected error: %v", err)
	}

	_, err = uc.Update(context.Background(), created.ID, "user-1", entity.UpdateExpenseInput{
		CategoryID: "cat-2", // belongs to user-2, not user-1
		Amount:     2000,
		SpentAt:    time.Now(),
	})
	if !errors.Is(err, entity.ErrCategoryInvalid) {
		t.Fatalf("expected ErrCategoryInvalid, got %v", err)
	}
}
