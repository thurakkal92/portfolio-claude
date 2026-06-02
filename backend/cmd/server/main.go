package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	root "github.com/thurakkal92/portfolio/backend"
	"github.com/thurakkal92/portfolio/backend/internal/config"
	"github.com/thurakkal92/portfolio/backend/internal/contact"
	"github.com/thurakkal92/portfolio/backend/internal/content"
	"github.com/thurakkal92/portfolio/backend/internal/db"
	apphttp "github.com/thurakkal92/portfolio/backend/internal/http"
	"github.com/thurakkal92/portfolio/backend/internal/seed"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	config.LoadDotenv()
	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "err", err)
		os.Exit(1)
	}

	if cfg.RunMigrations {
		logger.Info("running migrations")
		if err := db.Migrate(cfg.DatabaseURL, root.MigrationsFS); err != nil {
			logger.Error("migrations failed", "err", err)
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("db connect failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if cfg.SeedOnStart {
		logger.Info("seeding content")
		if err := seed.Run(ctx, pool); err != nil {
			logger.Error("seed failed", "err", err)
			os.Exit(1)
		}
	}

	contentSvc := content.New(pool)

	var mailer contact.Mailer
	if cfg.ResendAPIKey != "" {
		mailer = contact.NewResend(cfg.ResendAPIKey)
	} else {
		logger.Warn("RESEND_API_KEY not set — contact submissions will be stored but not emailed")
	}
	contactSvc := contact.NewService(pool, mailer, cfg.ContactFrom, cfg.ContactTo, cfg.RateLimitPerHr, logger)

	handler := apphttp.NewRouter(apphttp.Deps{
		Config:         cfg,
		ContentService: contentSvc,
		ContactService: contactSvc,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
			cancel()
		}
	}()

	<-stop
	logger.Info("shutting down")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}
