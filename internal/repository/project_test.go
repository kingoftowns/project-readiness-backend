package repository

import (
	"context"
	"testing"
	"time"

	"github.com/user/go-backend/internal/database"
	"github.com/user/go-backend/internal/models"
)

// Note: These tests require a running PostgreSQL instance
// Run: docker-compose up -d postgres
// Or use the test database from devcontainer setup
func setupTestDB(t *testing.T) *database.DB {
	// Skip tests if no test database is available
	testURL := "postgres://postgres:postgres@localhost:5432/gitlab_readiness_test?sslmode=disable"
	cfg := database.Config{
		URL: testURL,
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		t.Skipf("Skipping test - PostgreSQL not available: %v", err)
	}

	_, _ = db.Exec("DROP TABLE IF EXISTS gitlab_projects")

	schema := `
		CREATE TABLE gitlab_projects (
			project_id VARCHAR(255) PRIMARY KEY,
			project_present BOOLEAN DEFAULT FALSE,
			app_name_set BOOLEAN DEFAULT FALSE,
			moab_id_set BOOLEAN DEFAULT FALSE,
			codeowners_exists BOOLEAN DEFAULT FALSE,
			branch_protection_enabled BOOLEAN DEFAULT FALSE,
			codeowner_approval_required BOOLEAN DEFAULT FALSE,
			push_merge_restricted BOOLEAN DEFAULT FALSE,
			force_push_disabled BOOLEAN DEFAULT FALSE,
			push_rules_enabled BOOLEAN DEFAULT FALSE,
			min_approvals_required BOOLEAN DEFAULT FALSE,
			author_approval_prevented BOOLEAN DEFAULT FALSE,
			committer_approval_prevented BOOLEAN DEFAULT FALSE,
			approvals_removed_on_commit BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestProjectRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)
	ctx := context.Background()

	project := &models.Project{
		ProjectID:        "test-123",
		ProjectPresent:   true,
		AppNameSet:       true,
		MoabIDSet:        false,
		CodeownersExists: true,
	}

	err := repo.Create(ctx, project)
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	retrieved, err := repo.GetByID(ctx, "test-123")
	if err != nil {
		t.Fatalf("failed to retrieve project: %v", err)
	}

	if retrieved.ProjectID != project.ProjectID {
		t.Errorf("ProjectID = %v, want %v", retrieved.ProjectID, project.ProjectID)
	}
	if retrieved.ProjectPresent != project.ProjectPresent {
		t.Errorf("ProjectPresent = %v, want %v", retrieved.ProjectPresent, project.ProjectPresent)
	}
	if retrieved.AppNameSet != project.AppNameSet {
		t.Errorf("AppNameSet = %v, want %v", retrieved.AppNameSet, project.AppNameSet)
	}
	if retrieved.MoabIDSet != project.MoabIDSet {
		t.Errorf("MoabIDSet = %v, want %v", retrieved.MoabIDSet, project.MoabIDSet)
	}

	err = repo.Create(ctx, project)
	if err == nil {
		t.Error("expected error when creating duplicate project")
	}
}

func TestProjectRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)
	ctx := context.Background()

	project := &models.Project{
		ProjectID:      "update-test",
		ProjectPresent: true,
		AppNameSet:     false,
		MoabIDSet:      false,
	}

	if err := repo.Create(ctx, project); err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	project.AppNameSet = true
	project.MoabIDSet = true

	if err := repo.Update(ctx, project); err != nil {
		t.Fatalf("failed to update project: %v", err)
	}

	retrieved, err := repo.GetByID(ctx, "update-test")
	if err != nil {
		t.Fatalf("failed to retrieve project: %v", err)
	}

	if !retrieved.AppNameSet {
		t.Error("AppNameSet should be true after update")
	}
	if !retrieved.MoabIDSet {
		t.Error("MoabIDSet should be true after update")
	}

	nonExistent := &models.Project{
		ProjectID: "does-not-exist",
	}
	err = repo.Update(ctx, nonExistent)
	if err == nil {
		t.Error("expected error when updating non-existent project")
	}
}

func TestProjectRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)
	ctx := context.Background()

	project := &models.Project{
		ProjectID:      "delete-test",
		ProjectPresent: true,
	}

	if err := repo.Create(ctx, project); err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	if err := repo.Delete(ctx, "delete-test"); err != nil {
		t.Fatalf("failed to delete project: %v", err)
	}

	_, err := repo.GetByID(ctx, "delete-test")
	if err == nil {
		t.Error("expected error when getting deleted project")
	}

	err = repo.Delete(ctx, "does-not-exist")
	if err == nil {
		t.Error("expected error when deleting non-existent project")
	}
}

func TestProjectRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)
	ctx := context.Background()

	projects := []models.Project{
		{ProjectID: "proj-1", ProjectPresent: true},
		{ProjectID: "proj-2", ProjectPresent: true},
		{ProjectID: "proj-3", ProjectPresent: true},
		{ProjectID: "proj-4", ProjectPresent: true},
		{ProjectID: "proj-5", ProjectPresent: true},
	}

	for i := range projects {
		time.Sleep(1 * time.Millisecond)
		if err := repo.Create(ctx, &projects[i]); err != nil {
			t.Fatalf("failed to create project %s: %v", projects[i].ProjectID, err)
		}
	}

	tests := []struct {
		name   string
		limit  int
		offset int
		want   int
	}{
		{"first page", 2, 0, 2},
		{"second page", 2, 2, 2},
		{"third page", 2, 4, 1},
		{"all items", 10, 0, 5},
		{"offset beyond total", 10, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := repo.List(ctx, tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("failed to list projects: %v", err)
			}

			if len(results) != tt.want {
				t.Errorf("List() returned %d items, want %d", len(results), tt.want)
			}
		})
	}

	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("failed to count projects: %v", err)
	}

	if count != 5 {
		t.Errorf("Count() = %d, want 5", count)
	}
}

func TestProjectRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProjectRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "does-not-exist")
	if err == nil {
		t.Error("expected error when getting non-existent project")
	}

	if err.Error() != "project not found" {
		t.Errorf("unexpected error message: %v", err)
	}
}
