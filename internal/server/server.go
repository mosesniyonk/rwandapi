package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/mosesniyonk/rwandapi/internal/middleware"
)

// New creates and configures the HTTP server with all middleware and routes.
func New() http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)
	r.Use(chimw.SetHeader("Content-Type", "application/json"))

	// Rate limiting: 100 requests per minute per IP
	rl := middleware.NewRateLimiter(100, time.Minute)
	r.Use(rl.Handler)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Root
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"name":"rwandapi","version":"1.0.0","docs":"https://github.com/mosesniyonk/rwandapi"}`))
	})

	// Register API routes
	registerRoutes(r)

	return r
}
