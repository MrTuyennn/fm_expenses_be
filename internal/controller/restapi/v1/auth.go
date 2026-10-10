package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"fm_expenses_be/internal/controller/restapi/v1/request"
	"fm_expenses_be/internal/controller/restapi/v1/response"
	"fm_expenses_be/internal/entity"
	envelope "fm_expenses_be/pkg/response"
)

func (h *V1) register(w http.ResponseWriter, r *http.Request) {
	var body request.Register
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, tokens, err := h.authUsecase.Register(r.Context(), body.Email, body.Password, body.FullName)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserAlreadyExists):
			envelope.Error(w, http.StatusConflict, "user already exists")
		case errors.Is(err, entity.ErrPasswordTooShort):
			envelope.Error(w, http.StatusBadRequest, err.Error())
		default:
			envelope.Error(w, http.StatusInternalServerError, "failed to register")
		}
		return
	}

	envelope.JSON(w, http.StatusCreated, response.RegisterResult{
		User:   response.FromUser(user),
		Tokens: response.FromTokenPair(tokens),
	})
}

func (h *V1) login(w http.ResponseWriter, r *http.Request) {
	var body request.Login
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := h.authUsecase.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			envelope.Error(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		envelope.Error(w, http.StatusInternalServerError, "failed to login")
		return
	}

	envelope.JSON(w, http.StatusOK, response.FromTokenPair(tokens))
}

func (h *V1) refresh(w http.ResponseWriter, r *http.Request) {
	var body request.Refresh
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := h.authUsecase.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		if errors.Is(err, entity.ErrRefreshTokenInvalid) {
			envelope.Error(w, http.StatusUnauthorized, "refresh token invalid or expired")
			return
		}
		envelope.Error(w, http.StatusInternalServerError, "failed to refresh")
		return
	}

	envelope.JSON(w, http.StatusOK, response.FromTokenPair(tokens))
}

func (h *V1) logout(w http.ResponseWriter, r *http.Request) {
	var body request.Refresh
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		envelope.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authUsecase.Logout(r.Context(), body.RefreshToken); err != nil {
		envelope.Error(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
