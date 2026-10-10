package v1

import (
	"fm_expenses_be/internal/usecase"

	"github.com/go-chi/chi/v5"
)

// NewAuthRoutes registers the public (no bearer token required) /auth/*
// routes onto r.
func NewAuthRoutes(r chi.Router, authUC usecase.Auth) {
	h := &V1{authUsecase: authUC}

	r.Route("/auth", func(sub chi.Router) {
		sub.Post("/register", h.register)
		sub.Post("/login", h.login)
		sub.Post("/refresh", h.refresh)
		sub.Post("/logout", h.logout)
	})
}

// NewRoutes registers every authenticated v1 route onto r.
func NewRoutes(r chi.Router, userUC usecase.User, categoryUC usecase.Category, expenseUC usecase.Expense) {
	h := &V1{userUsecase: userUC, categoryUsecase: categoryUC, expenseUsecase: expenseUC}

	r.Route("/users", func(sub chi.Router) {
		sub.Get("/me", h.getMe)
	})

	r.Route("/categories", func(sub chi.Router) {
		sub.Get("/", h.listCategories)
		sub.Post("/", h.createCategory)
		sub.Get("/{id}", h.getCategory)
		sub.Put("/{id}", h.updateCategory)
		sub.Delete("/{id}", h.deleteCategory)
	})

	r.Route("/expenses", func(sub chi.Router) {
		sub.Get("/", h.listExpenses)
		sub.Post("/", h.createExpense)
		sub.Get("/{id}", h.getExpense)
		sub.Put("/{id}", h.updateExpense)
		sub.Delete("/{id}", h.deleteExpense)
	})
}
