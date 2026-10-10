package v1

import (
	"net/http"

	"fm_expenses_be/internal/controller/restapi/v1/response"
	"fm_expenses_be/internal/middleware"
	envelope "fm_expenses_be/pkg/response"
)

func (h *V1) getMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	profile, err := h.userUsecase.GetProfile(r.Context(), userID)
	if err != nil {
		envelope.Error(w, http.StatusNotFound, "user not found")
		return
	}

	envelope.JSON(w, http.StatusOK, response.FromUser(profile))
}
