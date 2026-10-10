package main

import (
	"fm_expenses_be/internal/app"
	"fm_expenses_be/internal/config"
)

func main() {
	cfg := config.Load()
	app.Run(cfg)
}
