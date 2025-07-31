package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/user/go-backend/internal/database"
	"github.com/user/go-backend/internal/models"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error

	GetByID(ctx context.Context, projectID string) (*models.Project, error)

	Update(ctx context.Context, project *models.Project) error

	Delete(ctx context.Context, projectID string) error

	List(ctx context.Context, limit, offset int) ([]*models.Project, error)

	Count(ctx context.Context) (int, error)
}

type projectRepo struct {
	db *database.DB
}

func NewProjectRepository(db *database.DB) ProjectRepository {
	return &projectRepo{db: db}
}


func (r *projectRepo) Create(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO gitlab_projects (
			project_id, project_present, app_name_set, moab_id_set,
			codeowners_exists, branch_protection_enabled, codeowner_approval_required,
			push_merge_restricted, force_push_disabled, push_rules_enabled,
			min_approvals_required, author_approval_prevented, committer_approval_prevented,
			approvals_removed_on_commit, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		project.ProjectID,
		project.ProjectPresent,
		project.AppNameSet,
		project.MoabIDSet,
		project.CodeownersExists,
		project.BranchProtectionEnabled,
		project.CodeownerApprovalRequired,
		project.PushMergeRestricted,
		project.ForcePushDisabled,
		project.PushRulesEnabled,
		project.MinApprovalsRequired,
		project.AuthorApprovalPrevented,
		project.CommitterApprovalPrevented,
		project.ApprovalsRemovedOnCommit,
		project.CreatedAt,
		project.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

func (r *projectRepo) GetByID(ctx context.Context, projectID string) (*models.Project, error) {
	query := `
		SELECT 
			project_id, project_present, app_name_set, moab_id_set,
			codeowners_exists, branch_protection_enabled, codeowner_approval_required,
			push_merge_restricted, force_push_disabled, push_rules_enabled,
			min_approvals_required, author_approval_prevented, committer_approval_prevented,
			approvals_removed_on_commit, created_at, updated_at
		FROM gitlab_projects
		WHERE project_id = $1
	`

	project := &models.Project{}
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&project.ProjectID,
		&project.ProjectPresent,
		&project.AppNameSet,
		&project.MoabIDSet,
		&project.CodeownersExists,
		&project.BranchProtectionEnabled,
		&project.CodeownerApprovalRequired,
		&project.PushMergeRestricted,
		&project.ForcePushDisabled,
		&project.PushRulesEnabled,
		&project.MinApprovalsRequired,
		&project.AuthorApprovalPrevented,
		&project.CommitterApprovalPrevented,
		&project.ApprovalsRemovedOnCommit,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

func (r *projectRepo) Update(ctx context.Context, project *models.Project) error {
	query := `
		UPDATE gitlab_projects SET
			project_present = $2,
			app_name_set = $3,
			moab_id_set = $4,
			codeowners_exists = $5,
			branch_protection_enabled = $6,
			codeowner_approval_required = $7,
			push_merge_restricted = $8,
			force_push_disabled = $9,
			push_rules_enabled = $10,
			min_approvals_required = $11,
			author_approval_prevented = $12,
			committer_approval_prevented = $13,
			approvals_removed_on_commit = $14,
			updated_at = $15
		WHERE project_id = $1
	`

	project.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		project.ProjectID,
		project.ProjectPresent,
		project.AppNameSet,
		project.MoabIDSet,
		project.CodeownersExists,
		project.BranchProtectionEnabled,
		project.CodeownerApprovalRequired,
		project.PushMergeRestricted,
		project.ForcePushDisabled,
		project.PushRulesEnabled,
		project.MinApprovalsRequired,
		project.AuthorApprovalPrevented,
		project.CommitterApprovalPrevented,
		project.ApprovalsRemovedOnCommit,
		project.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

func (r *projectRepo) Delete(ctx context.Context, projectID string) error {
	query := `DELETE FROM gitlab_projects WHERE project_id = $1`

	result, err := r.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

func (r *projectRepo) List(ctx context.Context, limit, offset int) ([]*models.Project, error) {
	query := `
		SELECT 
			project_id, project_present, app_name_set, moab_id_set,
			codeowners_exists, branch_protection_enabled, codeowner_approval_required,
			push_merge_restricted, force_push_disabled, push_rules_enabled,
			min_approvals_required, author_approval_prevented, committer_approval_prevented,
			approvals_removed_on_commit, created_at, updated_at
		FROM gitlab_projects
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ProjectID,
			&project.ProjectPresent,
			&project.AppNameSet,
			&project.MoabIDSet,
			&project.CodeownersExists,
			&project.BranchProtectionEnabled,
			&project.CodeownerApprovalRequired,
			&project.PushMergeRestricted,
			&project.ForcePushDisabled,
			&project.PushRulesEnabled,
			&project.MinApprovalsRequired,
			&project.AuthorApprovalPrevented,
			&project.CommitterApprovalPrevented,
			&project.ApprovalsRemovedOnCommit,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return projects, nil
}

func (r *projectRepo) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM gitlab_projects`

	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count projects: %w", err)
	}

	return count, nil
}
