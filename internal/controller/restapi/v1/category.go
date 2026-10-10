package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"fm_expenses_be/internal/controller/restapi/v1/request"
	"fm_expenses_be/internal/controller/restapi/v1/response"
	"fm_expenses_be/internal/entity"
	"fm_expenses_be/internal/middleware"
	envelope "fm_expenses_be/pkg/response"

	"github.com/go-chi/chi/v5"
)

// writeCategoryError maps domain sentinel errors to the right HTTP status.
// Any other error (unexpected DB failure, etc.) falls back to 500.
func writeCategoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrNotFound):
		envelope.Error(w, http.StatusNotFound, "category not found")
	case errors.Is(err, entity.ErrInvalidCategoryType):
		envelope.Error(w, http.StatusBadRequest, err.Error())
	default:
		envelope.Error(w, http.StatusInternalServerError, "failed to save category")
	}
}

func (h *V1) createCategory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var body request.Category
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.categoryUsecase.Create(r.Context(), userID, body.Name, body.Icon, body.Type)
	if err != nil {
		writeCategoryError(w, err)
		return
	}

	envelope.JSON(w, http.StatusCreated, response.FromCategory(c))
}

func (h *V1) listCategories(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	items, err := h.categoryUsecase.List(r.Context(), userID)
	if err != nil {
		envelope.Error(w, http.StatusInternalServerError, "failed to list categories")
		return
	}

	envelope.JSON(w, http.StatusOK, response.FromCategoryList(items))
}

func (h *V1) getCategory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	c, err := h.categoryUsecase.Get(r.Context(), id, userID)
	if err != nil {
		envelope.Error(w, http.StatusNotFound, "category not found")
		return
	}

	envelope.JSON(w, http.StatusOK, response.FromCategory(c))
}

func (h *V1) updateCategory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var body request.Category
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.categoryUsecase.Update(r.Context(), id, userID, body.Name, body.Icon, body.Type)
	if err != nil {
		writeCategoryError(w, err)
		return
	}

	envelope.JSON(w, http.StatusOK, response.FromCategory(c))
}

func (h *V1) deleteCategory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.categoryUsecase.Delete(r.Context(), id, userID); err != nil {
		envelope.Error(w, http.StatusInternalServerError, "failed to delete category")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
