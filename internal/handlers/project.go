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

type ProjectHandler struct {
	repo   repository.ProjectRepository
	logger *slog.Logger
}

func NewProjectHandler(repo repository.ProjectRepository, logger *slog.Logger) *ProjectHandler {
	return &ProjectHandler{
		repo:   repo,
		logger: logger,
	}
}

// ListProjects handles GET /api/v1/projects
// It returns a paginated list of projects
//
//	@Summary		List projects
//	@Description	Get a paginated list of projects with their readiness status
//	@Tags			gitlab
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Number of items to return (max 100)"	default(50)
//	@Param			offset	query		int	false	"Number of items to skip"				default(0)
//	@Success		200		{object}	models.PaginatedResponse	"List of projects with pagination metadata"
//	@Failure		500		{object}	models.ErrorResponse	"Internal server error"
//	@Router			/gitlab/projects [get]
func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100
			}
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	projects, err := h.repo.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list projects", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve projects")
		return
	}

	total, err := h.repo.Count(ctx)
	if err != nil {
		h.logger.Error("failed to count projects", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to count projects")
		return
	}

	pagination := &models.PaginationMeta{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
	response := models.NewPaginatedResponse(http.StatusOK, "Projects retrieved successfully", projects, pagination)

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetProject handles GET /api/v1/projects/{id}
// It returns a single project by ID
//
//	@Summary		Get project by ID
//	@Description	Get a single project with all readiness check data
//	@Tags			gitlab
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Success		200	{object}	models.SuccessResponse	"Project details with readiness status"
//	@Failure		400	{object}	models.ErrorResponse	"Bad request"
//	@Failure		404	{object}	models.ErrorResponse	"Project ID not found"
//	@Failure		500	{object}	models.ErrorResponse	"Internal server error"
//	@Router			/gitlab/projects/{id} [get]
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
			h.respondWithError(w, http.StatusNotFound, "project_id not found")
			return
		}
		h.logger.Error("failed to get project", "error", err, "project_id", projectID)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve project")
		return
	}

	response := models.NewSuccessResponse(http.StatusOK, "Project retrieved successfully", project)

	h.respondWithJSON(w, http.StatusOK, response)
}

// CreateProject handles POST /api/v1/projects
// It creates a new project
//
//	@Summary		Create a new project
//	@Description	Create a new project with initial readiness checks
//	@Tags			gitlab
//	@Accept			json
//	@Produce		json
//	@Param			project	body		models.Project			true	"Project data"
//	@Success		201		{object}	models.SuccessResponse	"Created project"
//	@Failure		400		{object}	models.ErrorResponse	"Bad request"
//	@Failure		409		{object}	models.ErrorResponse	"Project already exists"
//	@Failure		500		{object}	models.ErrorResponse	"Internal server error"
//	@Router			/gitlab/projects [post]
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if project.ProjectID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	existing, err := h.repo.GetByID(ctx, project.ProjectID)
	if err == nil && existing != nil {
		h.respondWithError(w, http.StatusConflict, "Project already exists")
		return
	}

	if err := h.repo.Create(ctx, &project); err != nil {
		h.logger.Error("failed to create project", "error", err, "project_id", project.ProjectID)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	h.logger.Info("project created", "project_id", project.ProjectID)
	response := models.NewSuccessResponse(http.StatusCreated, "Project created successfully", project)
	h.respondWithJSON(w, http.StatusCreated, response)
}

// UpdateProject handles PUT /api/v1/projects/{id}
// It updates an existing project
//
//	@Summary		Update project
//	@Description	Update an existing project's readiness checks
//	@Tags			gitlab
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Project ID"
//	@Param			project	body		models.Project		true	"Updated project data"
//	@Success		200		{object}	models.SuccessResponse	"Updated project"
//	@Failure		400		{object}	models.ErrorResponse	"Bad request"
//	@Failure		404		{object}	models.ErrorResponse	"Project not found"
//	@Failure		500		{object}	models.ErrorResponse	"Internal server error"
//	@Router			/gitlab/projects/{id} [put]
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

	project.ProjectID = projectID

	if err := h.repo.Update(ctx, &project); err != nil {
		if err.Error() == "project not found" {
			h.respondWithError(w, http.StatusNotFound, "project_id not found")
			return
		}
		h.logger.Error("failed to update project", "error", err, "project_id", projectID)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	h.logger.Info("project updated", "project_id", projectID)
	response := models.NewSuccessResponse(http.StatusOK, "Project updated successfully", project)
	h.respondWithJSON(w, http.StatusOK, response)
}

// DeleteProject handles DELETE /api/v1/projects/{id}
// It deletes a project
//
//	@Summary		Delete project
//	@Description	Delete a project by ID
//	@Tags			gitlab
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Project ID"
//	@Success		204	{object}	models.SuccessResponse	"Project deleted successfully"
//	@Failure		400	{object}	models.ErrorResponse	"Bad request"
//	@Failure		404	{object}	models.ErrorResponse	"Project ID not found"
//	@Failure		500	{object}	models.ErrorResponse	"Internal server error"
//	@Router			/gitlab/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := chi.URLParam(r, "id")

	if projectID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	if err := h.repo.Delete(ctx, projectID); err != nil {
		if err.Error() == "project not found" {
			h.respondWithError(w, http.StatusNotFound, "project_id not found")
			return
		}
		h.logger.Error("failed to delete project", "error", err, "project_id", projectID)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	h.logger.Info("project deleted", "project_id", projectID)
	response := models.NewSuccessResponse(http.StatusNoContent, "Project deleted successfully", nil)
	h.respondWithJSON(w, http.StatusNoContent, response)
}

// HealthCheck handles GET /api/v1/health
// It returns the health status of the API
//
//	@Summary		Health check
//	@Description	Check if the API is healthy and running
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.SuccessResponse	"Health status"
//	@Router			/health [get]
func (h *ProjectHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"service": "gitlab-readiness-api",
		"version": "1.0.0",
	}
	response := models.NewSuccessResponse(http.StatusOK, "Service is healthy", data)
	h.respondWithJSON(w, http.StatusOK, response)
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
	response := models.NewErrorResponse(code, message)
	h.respondWithJSON(w, code, response)
}
