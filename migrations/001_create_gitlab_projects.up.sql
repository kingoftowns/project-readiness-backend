-- Create the gitlab_projects table
-- This table stores GitLab project information and their production readiness checks
CREATE TABLE IF NOT EXISTS gitlab_projects (
    project_id TEXT PRIMARY KEY,
    
    -- GitLab presence checks
    project_present BOOLEAN NOT NULL DEFAULT FALSE,
    app_name_set BOOLEAN NOT NULL DEFAULT FALSE,
    moab_id_set BOOLEAN NOT NULL DEFAULT FALSE,
    codeowners_exists BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Branch protection checks
    branch_protection_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    codeowner_approval_required BOOLEAN NOT NULL DEFAULT FALSE,
    push_merge_restricted BOOLEAN NOT NULL DEFAULT FALSE,
    force_push_disabled BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Merge request checks
    push_rules_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    min_approvals_required BOOLEAN NOT NULL DEFAULT FALSE,
    author_approval_prevented BOOLEAN NOT NULL DEFAULT FALSE,
    committer_approval_prevented BOOLEAN NOT NULL DEFAULT FALSE,
    approvals_removed_on_commit BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_gitlab_projects_created_at ON gitlab_projects(created_at DESC);

-- Create a trigger to automatically update the updated_at timestamp