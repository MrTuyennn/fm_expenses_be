package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"fm_expenses_be/internal/controller/restapi/v1/request"
	"fm_expenses_be/internal/controller/restapi/v1/response"
	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/middleware"
	envelope "fm_expenses_be/pkg/response"

	"github.com/go-chi/chi/v5"
)

const timeLayout = time.RFC3339

func writeExpenseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrNotFound):
		envelope.Error(w, http.StatusNotFound, "expense not found")
	case errors.Is(err, entity.ErrCategoryInvalid):
		envelope.Error(w, http.StatusBadRequest, err.Error())
	default:
		envelope.Error(w, http.StatusInternalServerError, "failed to save expense")
	}
}

func (h *V1) createExpense(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var body request.Expense
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	spentAt, err := time.Parse(timeLayout, body.SpentAt)
	if err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid spent_at, expected RFC3339")
		return
	}

	e, err := h.expenseUsecase.Create(r.Context(), entity.CreateExpenseInput{
		UserID:     userID,
		CategoryID: body.CategoryID,
		Amount:     body.Amount,
		Currency:   body.Currency,
		Note:       body.Note,
		SpentAt:    spentAt,
	})
	if err != nil {
		writeExpenseError(w, err)
		return
	}

	envelope.JSON(w, http.StatusCreated, response.FromExpense(e))
}

func (h *V1) listExpenses(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	filter := entity.ExpenseFilter{
		CategoryID: r.URL.Query().Get("category_id"),
	}
	if from := r.URL.Query().Get("from"); from != "" {
		if t, err := time.Parse(timeLayout, from); err == nil {
			filter.From = &t
		}
	}
	if to := r.URL.Query().Get("to"); to != "" {
		if t, err := time.Parse(timeLayout, to); err == nil {
			filter.To = &t
		}
	}

	items, err := h.expenseUsecase.List(r.Context(), userID, filter)
	if err != nil {
		envelope.Error(w, http.StatusInternalServerError, "failed to list expenses")
		return
	}

	envelope.JSON(w, http.StatusOK, response.FromExpenseList(items))
}

func (h *V1) getExpense(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	e, err := h.expenseUsecase.Get(r.Context(), id, userID)
	if err != nil {
		envelope.Error(w, http.StatusNotFound, "expense not found")
		return
	}

	envelope.JSON(w, http.StatusOK, response.FromExpense(e))
}

func (h *V1) updateExpense(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var body request.Expense
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	spentAt, err := time.Parse(timeLayout, body.SpentAt)
	if err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid spent_at, expected RFC3339")
		return
	}

	e, err := h.expenseUsecase.Update(r.Context(), id, userID, entity.UpdateExpenseInput{
		CategoryID: body.CategoryID,
		Amount:     body.Amount,
		Currency:   body.Currency,
		Note:       body.Note,
		SpentAt:    spentAt,
	})
	if err != nil {
		writeExpenseError(w, err)
		return
	}

	envelope.JSON(w, http.StatusOK, response.FromExpense(e))
}

func (h *V1) deleteExpense(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.expenseUsecase.Delete(r.Context(), id, userID); err != nil {
		envelope.Error(w, http.StatusInternalServerError, "failed to delete expense")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
