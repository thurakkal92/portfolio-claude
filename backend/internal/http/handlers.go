package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thurakkal92/portfolio/backend/internal/contact"
	"github.com/thurakkal92/portfolio/backend/internal/content"
	"github.com/thurakkal92/portfolio/backend/internal/i18n"
)

func handleContent(svc *content.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.URL.Query().Get("locale")
		if locale == "" {
			locale = i18n.Default
		}
		if !i18n.IsSupported(locale) {
			writeError(w, http.StatusBadRequest, "unsupported_locale", "Locale is not supported")
			return
		}
		payload, err := svc.Get(r.Context(), locale)
		if err != nil {
			if errors.Is(err, content.ErrLocaleNotFound) {
				writeError(w, http.StatusNotFound, "not_found", "Content for locale not seeded")
				return
			}
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=60, s-maxage=300")
		writeJSON(w, http.StatusOK, payload)
	}
}

func handleContact(svc *contact.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contact.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_body", "Could not parse JSON")
			return
		}
		req.IPAddress = r.RemoteAddr
		req.UserAgent = r.UserAgent()

		if err := svc.Submit(r.Context(), &req); err != nil {
			var valErr *contact.ValidationError
			switch {
			case errors.As(err, &valErr):
				writeValidationError(w, valErr.Fields)
				return
			case errors.Is(err, contact.ErrRateLimited):
				writeError(w, http.StatusTooManyRequests, "rate_limited", "Too many submissions. Try again later.")
				return
			case errors.Is(err, contact.ErrHoneypot):
				// Silently accept honeypot trips so bots don't learn anything.
				writeJSON(w, http.StatusOK, map[string]any{"ok": true})
				return
			default:
				writeError(w, http.StatusInternalServerError, "internal_error", "Could not send message")
				return
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{
		"error":   code,
		"message": msg,
	})
}

func writeValidationError(w http.ResponseWriter, fields map[string]string) {
	writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
		"error":  "validation_failed",
		"fields": fields,
	})
}
