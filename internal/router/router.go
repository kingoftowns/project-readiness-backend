package router

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/user/go-backend/internal/handlers"
	"github.com/user/go-backend/internal/models"

	_ "github.com/user/go-backend/docs" // This is required for Swagger
)

func New(projectHandler *handlers.ProjectHandler, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)                 // Add request ID for tracing
	r.Use(middleware.RealIP)                    // Get real IP from headers
	r.Use(middleware.Recoverer)                 // Recover from panics
	r.Use(LoggerMiddleware(logger))             // Custom logging middleware
	r.Use(middleware.Timeout(60 * time.Second)) // Request timeout

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Use relative URL instead of absolute
	))

	r.Get("/api/v1/health", projectHandler.HealthCheck)

	r.Route("/api/v1/gitlab/projects", func(r chi.Router) {
		r.Get("/", projectHandler.ListProjects)         // GET /api/v1/gitlab/projects
		r.Post("/", projectHandler.CreateProject)       // POST /api/v1/gitlab/projects
		r.Get("/{id}", projectHandler.GetProject)       // GET /api/v1/gitlab/projects/{id}
		r.Put("/{id}", projectHandler.UpdateProject)    // PUT /api/v1/gitlab/projects/{id}
		r.Delete("/{id}", projectHandler.DeleteProject) // DELETE /api/v1/gitlab/projects/{id}
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		response := models.NewErrorResponse(http.StatusNotFound, "Route not found")
		json.NewEncoder(w).Encode(response)
	})

	return r
}

func LoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

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

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
