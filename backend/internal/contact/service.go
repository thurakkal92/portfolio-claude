package contact

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Mailer interface {
	Send(ctx context.Context, e Email) (string, error)
}

type Service struct {
	pool           *pgxpool.Pool
	mailer         Mailer
	fromAddr       string
	toAddr         string
	rateLimitPerHr int
	logger         *slog.Logger
}

func NewService(pool *pgxpool.Pool, mailer Mailer, from, to string, rateLimit int, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		pool:           pool,
		mailer:         mailer,
		fromAddr:       from,
		toAddr:         to,
		rateLimitPerHr: rateLimit,
		logger:         logger,
	}
}

func (s *Service) Submit(ctx context.Context, req *Request) error {
	// Honeypot: silently reject.
	if strings.TrimSpace(req.Website) != "" {
		return ErrHoneypot
	}

	if err := validate(req); err != nil {
		return err
	}

	ip := normalizeIP(req.IPAddress)
	if ip != "" && s.rateLimitPerHr > 0 {
		ok, err := s.checkRateLimit(ctx, ip)
		if err != nil {
			s.logger.Warn("rate limit check failed", "err", err)
		} else if !ok {
			return ErrRateLimited
		}
	}

	id, err := s.persist(ctx, req, ip)
	if err != nil {
		return fmt.Errorf("persist: %w", err)
	}

	if s.mailer == nil {
		s.logger.Warn("contact submission stored but mailer not configured", "id", id)
		return nil
	}

	providerID, err := s.mailer.Send(ctx, Email{
		From:    s.fromAddr,
		To:      s.toAddr,
		ReplyTo: req.Email,
		Subject: fmt.Sprintf("Portfolio contact from %s", req.Name),
		Text:    buildText(req),
	})
	if err != nil {
		s.logger.Error("send email failed", "err", err, "id", id)
		return err
	}

	_, _ = s.pool.Exec(ctx, `
		UPDATE contact_submissions SET email_sent = TRUE, email_provider_id = $1 WHERE id = $2`,
		providerID, id)
	return nil
}

func (s *Service) persist(ctx context.Context, req *Request, ip string) (int64, error) {
	var ipParam any
	if ip != "" {
		ipParam = ip
	}
	var id int64
	err := s.pool.QueryRow(ctx, `
		INSERT INTO contact_submissions (name, email, message, locale, ip_address, user_agent)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6)
		RETURNING id`,
		req.Name, req.Email, req.Message, req.Locale, ipParam, req.UserAgent,
	).Scan(&id)
	return id, err
}

func (s *Service) checkRateLimit(ctx context.Context, ip string) (bool, error) {
	window := time.Hour
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	var count int
	var windowStart time.Time
	err = tx.QueryRow(ctx, `
		SELECT attempt_count, window_start FROM contact_rate_limit WHERE ip_address = $1 FOR UPDATE`, ip).
		Scan(&count, &windowStart)

	now := time.Now()
	switch {
	case errors.Is(err, pgxNoRows()):
		_, err = tx.Exec(ctx, `
			INSERT INTO contact_rate_limit (ip_address, attempt_count, window_start)
			VALUES ($1, 1, $2)`, ip, now)
		if err != nil {
			return false, err
		}
	case err != nil:
		return false, err
	default:
		if now.Sub(windowStart) > window {
			count = 0
			windowStart = now
		}
		count++
		if count > s.rateLimitPerHr {
			_ = tx.Commit(ctx)
			return false, nil
		}
		_, err = tx.Exec(ctx, `
			UPDATE contact_rate_limit SET attempt_count = $1, window_start = $2 WHERE ip_address = $3`,
			count, windowStart, ip)
		if err != nil {
			return false, err
		}
	}
	return true, tx.Commit(ctx)
}

func buildText(req *Request) string {
	return fmt.Sprintf("From: %s <%s>\nLocale: %s\n\n%s\n", req.Name, req.Email, req.Locale, req.Message)
}

func normalizeIP(addr string) string {
	if addr == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}
