// Package handlers contains the HTTP handlers for the API endpoints.
// These handlers follow standard Go HTTP patterns for educational clarity.
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/user/go-backend/internal/models"
	"github.com/user/go-backend/internal/repository"
)

// ProjectHandler handles HTTP requests for project readiness endpoints
type ProjectHandler struct {
	repo   repository.ProjectRepository
	logger *slog.Logger
}

// NewProjectHandler creates a new project handler instance
func NewProjectHandler(repo repository.ProjectRepository, logger *slog.Logger) *ProjectHandler {
	return &ProjectHandler{
		repo:   repo,
		logger: logger,
	}
}

// ListProjects handles GET /api/v1/projects
// It returns a paginated list of projects
func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters for pagination
	limit := 50 // Default limit
	offset := 0 // Default offset

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100 // Max limit
			}
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get projects from repository
	projects, err := h.repo.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list projects", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve projects")
		return
	}

	// Get total count for pagination metadata
	total, err := h.repo.Count(ctx)
	if err != nil {
		h.logger.Error("failed to count projects", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to count projects")
		return
	}

	// Create response with pagination metadata
	response := map[string]interface{}{
		"data": projects,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetProject handles GET /api/v1/projects/{id}
// It returns a single project by ID
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := chi.URLParam(r, "id")

	if projectID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	project, err := h.repo.GetByID(ctx, projectID)
	if err != nil {
		if err.Error() == "project not found" {
			h.respondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.logger.Error("failed to get project", "error", err, "project_id", projectID)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
		return
	}

	// Include additional computed fields in response
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"project":           project,
			"is_production_ready": project.IsProductionReady(),
			"failed_checks":      project.FailedChecks(),
		},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// CreateProject handles POST /api/v1/projects
// It creates a new project
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if project.ProjectID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	// Check if project already exists
	existing, err := h.repo.GetByID(ctx, project.ProjectID)
	if err == nil && existing != nil {
		h.respondWithError(w, http.StatusConflict, "Project already exists")
		return
	}

	// Create the project
	if err := h.repo.Create(ctx, &project); err != nil {
		h.logger.Error("failed to create project", "error", err, "project_id", project.ProjectID)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	h.logger.Info("project created", "project_id", project.ProjectID)
	h.respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"data": project,
	})
}

// UpdateProject handles PUT /api/v1/projects/{id}
// It updates an existing project
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := chi.URLParam(r, "id")

	if projectID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Ensure the project ID matches the URL parameter
	project.ProjectID = projectID

	// Update the project
	if err := h.repo.Update(ctx, &project); err != nil {
		if err.Error() == "project not found" {
			h.respondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.logger.Error("failed to update project", "error", err, "project_id", projectID)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	h.logger.Info("project updated", "project_id", projectID)
	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": project,
	})
}

// DeleteProject handles DELETE /api/v1/projects/{id}
// It deletes a project
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := chi.URLParam(r, "id")

	if projectID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	if err := h.repo.Delete(ctx, projectID); err != nil {
		if err.Error() == "project not found" {
			h.respondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.logger.Error("failed to delete project", "error", err, "project_id", projectID)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	h.logger.Info("project deleted", "project_id", projectID)
	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck handles GET /api/v1/health
// It returns the health status of the API
func (h *ProjectHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "healthy",
		"service": "gitlab-readiness-api",
	})
}

// Helper methods for consistent JSON responses

func (h *ProjectHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *ProjectHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
		},
	})
}