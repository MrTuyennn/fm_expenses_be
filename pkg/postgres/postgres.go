// Package postgres implements the Postgres connection, reusable outside this repo.
package postgres

import (
	"time"

	"fm_expenses_be/pkg/logger"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const slowQueryThreshold = 200 * time.Millisecond

// New opens a GORM/Postgres connection with a compact custom logger (see
// pkg/logger) instead of GORM's default — see that package for why.
func New(dsn string) (*gorm.DB, error) {
	return gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		Logger: logger.New(slowQueryThreshold),
	})
}
