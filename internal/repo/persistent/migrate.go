// Package persistent aggregates every Postgres-backed repository
// (user, category, expense) and the migration that creates their tables.
package persistent

import (
	"gorm.io/gorm"

	"fm_expenses_be/internal/repo/persistent/category"
	"fm_expenses_be/internal/repo/persistent/expense"
	"fm_expenses_be/internal/repo/persistent/refreshtoken"
	"fm_expenses_be/internal/repo/persistent/user"
)

// Migrate creates/updates tables for every persistence model.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&user.Model{},
		&refreshtoken.Model{},
		&category.Model{},
		&expense.Model{},
	)
}
