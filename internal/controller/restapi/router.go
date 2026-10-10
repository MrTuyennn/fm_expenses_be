// Package restapi builds the top-level HTTP router.
package restapi

import (
	"net/http"

	v1 "fm_expenses_be/internal/controller/restapi/v1"
	"fm_expenses_be/internal/middleware"
	"fm_expenses_be/internal/usecase"
	"fm_expenses_be/pkg/logger"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// NewRouter builds the full router: logging/recovery middleware, an
// unauthenticated /health check and /api/v1/auth/*, and the authenticated
// rest of /api/v1.
func NewRouter(authMW *middleware.Auth, authUC usecase.Auth, userUC usecase.User, categoryUC usecase.Category, expenseUC usecase.Expense) http.Handler {
	r := chi.NewRouter()
	r.Use(logger.HTTPRequestLogger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Group(func(public chi.Router) {
			v1.NewAuthRoutes(public, authUC)
		})

		api.Group(func(protected chi.Router) {
			protected.Use(authMW.Middleware)
			v1.NewRoutes(protected, userUC, categoryUC, expenseUC)
		})
	})

	return r
}
