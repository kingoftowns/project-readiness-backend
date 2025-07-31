// Package router sets up the HTTP routing for the API.
// It uses chi router for its simplicity and standard library compatibility.
package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/user/go-backend/internal/handlers"
)

// New creates and configures a new router with all routes and middleware
func New(projectHandler *handlers.ProjectHandler, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)     // Add request ID for tracing
	r.Use(middleware.RealIP)        // Get real IP from headers
	r.Use(middleware.Recoverer)     // Recover from panics
	r.Use(LoggerMiddleware(logger)) // Custom logging middleware
	r.Use(middleware.Timeout(60 * time.Second)) // Request timeout

	// Health check endpoint (no auth required)
	r.Get("/api/v1/health", projectHandler.HealthCheck)

	// API routes
	r.Route("/api/v1/projects", func(r chi.Router) {
		r.Get("/", projectHandler.ListProjects)       // GET /api/v1/projects
		r.Post("/", projectHandler.CreateProject)     // POST /api/v1/projects
		r.Get("/{id}", projectHandler.GetProject)     // GET /api/v1/projects/{id}
		r.Put("/{id}", projectHandler.UpdateProject)  // PUT /api/v1/projects/{id}
		r.Delete("/{id}", projectHandler.DeleteProject) // DELETE /api/v1/projects/{id}
	})

	// Catch-all for undefined routes
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":{"message":"Route not found"}}`))
	})

	return r
}

// LoggerMiddleware creates a custom logging middleware using slog
func LoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap ResponseWriter to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Process request
			next.ServeHTTP(wrapped, r)

			// Log request details
			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration", time.Since(start).String(),
				"request_id", middleware.GetReqID(r.Context()),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}