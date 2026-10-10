// Package v1 implements the v1 REST API controllers.
package v1

import "fm_expenses_be/internal/usecase"

// V1 holds every usecase this API version depends on.
type V1 struct {
	authUsecase     usecase.Auth
	userUsecase     usecase.User
	categoryUsecase usecase.Category
	expenseUsecase  usecase.Expense
}
