// Package logger holds custom, concise loggers for infrastructure
// components (GORM, HTTP requests), replacing their verbose stdlib/GORM
// defaults with a single colored line per event.
package logger

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/utils"

	gormlogger "gorm.io/gorm/logger"
)

const moduleRoot = "fm_expenses_be/"

// gormDBLogger is a compact gorm.io/gorm/logger.Interface implementation.
// Unlike GORM's default logger it:
//   - never logs gorm.ErrRecordNotFound, which this codebase uses as
//     expected control flow (category/expense "not found" lookups) — the
//     default logger prints a full error block for every one of these
//   - only prints a line for slow queries or real errors; normal queries
//     stay silent
//   - shortens the caller file path to be relative to the module instead
//     of an absolute filesystem path
type gormDBLogger struct {
	slowThreshold time.Duration
}

// New returns a gorm.io/gorm/logger.Interface for gorm.Config{Logger: ...}.
func New(slowThreshold time.Duration) gormlogger.Interface {
	return &gormDBLogger{slowThreshold: slowThreshold}
}

func (l *gormDBLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *gormDBLogger) Info(context.Context, string, ...interface{})  {}
func (l *gormDBLogger) Warn(context.Context, string, ...interface{})  {}
func (l *gormDBLogger) Error(context.Context, string, ...interface{}) {}

func (l *gormDBLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	elapsed := time.Since(begin)
	caller := shortenCaller(utils.FileWithLineNum())

	switch {
	case err != nil:
		sql, _ := fc()
		fmt.Printf("%s[db error]%s %s  %v\n  %s\n", colorRed, colorReset, caller, err, truncate(sql, 160))
	case l.slowThreshold > 0 && elapsed > l.slowThreshold:
		sql, rows := fc()
		fmt.Printf("%s[db slow %s]%s %s  rows=%d\n  %s\n", colorYellow, elapsed.Round(time.Millisecond), colorReset, caller, rows, truncate(sql, 160))
	}
}

func shortenCaller(fileWithLine string) string {
	if i := strings.Index(fileWithLine, moduleRoot); i != -1 {
		return fileWithLine[i+len(moduleRoot):]
	}
	return fileWithLine
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
