package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/thurakkal92/portfolio/backend/internal/config"
	"github.com/thurakkal92/portfolio/backend/internal/db"
	"github.com/thurakkal92/portfolio/backend/internal/seed"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	config.LoadDotenv()
	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "err", err)
		os.Exit(1)
	}
	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("db connect failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := seed.Run(ctx, pool); err != nil {
		logger.Error("seed failed", "err", err)
		os.Exit(1)
	}
	logger.Info("seed complete")
}
