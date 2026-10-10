// Package app wires every dependency and runs the HTTP server.
package app

import (
	"log"
	"net/http"

	"fm_expenses_be/internal/config"
	"fm_expenses_be/internal/controller/restapi"
	"fm_expenses_be/internal/middleware"
	"fm_expenses_be/internal/repo/persistent"
	persistCategoryRepo "fm_expenses_be/internal/repo/persistent/category"
	persistExpenseRepo "fm_expenses_be/internal/repo/persistent/expense"
	persistRefreshTokenRepo "fm_expenses_be/internal/repo/persistent/refreshtoken"
	persistUserRepo "fm_expenses_be/internal/repo/persistent/user"
	"fm_expenses_be/internal/usecase/auth"
	"fm_expenses_be/internal/usecase/category"
	"fm_expenses_be/internal/usecase/expense"
	"fm_expenses_be/internal/usecase/user"
	"fm_expenses_be/pkg/jwt"
	"fm_expenses_be/pkg/postgres"
)

// Run connects to the database, wires every dependency (repo → usecase →
// controller), and starts the HTTP server. It blocks until the server exits.
func Run(cfg *config.Config) {
	db, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := persistent.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	authMW := middleware.NewAuth(cfg.JWTSecret)
	jwtManager := jwt.New(cfg.JWTSecret, cfg.AccessTokenTTL)

	userRepo := persistUserRepo.New(db)
	refreshTokenRepo := persistRefreshTokenRepo.New(db)
	authUC := auth.New(userRepo, refreshTokenRepo, jwtManager, cfg.RefreshTokenTTL)

	userUC := user.New(userRepo)

	categoryRepo := persistCategoryRepo.New(db)
	categoryUC := category.New(categoryRepo)

	expenseUC := expense.New(persistExpenseRepo.New(db), categoryRepo)

	handler := restapi.NewRouter(authMW, authUC, userUC, categoryUC, expenseUC)

	addr := ":" + cfg.Port
	log.Printf("fm_expenses_be listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
