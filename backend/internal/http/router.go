package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/thurakkal92/portfolio/backend/internal/config"
	"github.com/thurakkal92/portfolio/backend/internal/contact"
	"github.com/thurakkal92/portfolio/backend/internal/content"
)

type Deps struct {
	Config         *config.Config
	ContentService *content.Service
	ContactService *contact.Service
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(corsMiddleware(d.Config.AllowedOrigins))

	r.Get("/healthz", handleHealth)
	r.Get("/api/content", handleContent(d.ContentService))
	r.Post("/api/contact", handleContact(d.ContactService))

	return r
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
