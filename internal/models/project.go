// Package models contains the domain models for the application.
// These models represent the core business entities and their rules.
package models

import (
	"time"
)

// Project represents a project and its production readiness checks.
// Each field corresponds to a specific check that validates if the project
// follows the required standards for production deployment.
type Project struct {
	// ProjectID is the unique identifier from GitLab
	ProjectID string `json:"project_id" db:"project_id"`

	// GitLab presence checks
	ProjectPresent    bool `json:"project_present" db:"project_present"`
	AppNameSet        bool `json:"app_name_set" db:"app_name_set"`
	MoabIDSet         bool `json:"moab_id_set" db:"moab_id_set"`
	CodeownersExists  bool `json:"codeowners_exists" db:"codeowners_exists"`

	// Branch protection checks
	BranchProtectionEnabled    bool `json:"branch_protection_enabled" db:"branch_protection_enabled"`
	CodeownerApprovalRequired  bool `json:"codeowner_approval_required" db:"codeowner_approval_required"`
	PushMergeRestricted        bool `json:"push_merge_restricted" db:"push_merge_restricted"`
	ForcePushDisabled          bool `json:"force_push_disabled" db:"force_push_disabled"`

	// Merge request checks
	PushRulesEnabled           bool `json:"push_rules_enabled" db:"push_rules_enabled"`
	MinApprovalsRequired       bool `json:"min_approvals_required" db:"min_approvals_required"`
	AuthorApprovalPrevented    bool `json:"author_approval_prevented" db:"author_approval_prevented"`
	CommitterApprovalPrevented bool `json:"committer_approval_prevented" db:"committer_approval_prevented"`
	ApprovalsRemovedOnCommit   bool `json:"approvals_removed_on_commit" db:"approvals_removed_on_commit"`

	// Metadata
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IsProductionReady checks if all required checks pass.
// A project is considered production-ready when all checks are true.
func (p *Project) IsProductionReady() bool {
	return p.ProjectPresent &&
		p.AppNameSet &&
		p.MoabIDSet &&
		p.CodeownersExists &&
		p.BranchProtectionEnabled &&
		p.CodeownerApprovalRequired &&
		p.PushMergeRestricted &&
		p.ForcePushDisabled &&
		p.PushRulesEnabled &&
		p.MinApprovalsRequired &&
		p.AuthorApprovalPrevented &&
		p.CommitterApprovalPrevented &&
		p.ApprovalsRemovedOnCommit
}

// FailedChecks returns a list of check names that are currently failing.
// This is useful for providing feedback about what needs to be fixed.
func (p *Project) FailedChecks() []string {
	var failed []string

	if !p.ProjectPresent {
		failed = append(failed, "project_present")
	}
	if !p.AppNameSet {
		failed = append(failed, "app_name_set")
	}
	if !p.MoabIDSet {
		failed = append(failed, "moab_id_set")
	}
	if !p.CodeownersExists {
		failed = append(failed, "codeowners_exists")
	}
	if !p.BranchProtectionEnabled {
		failed = append(failed, "branch_protection_enabled")
	}
	if !p.CodeownerApprovalRequired {
		failed = append(failed, "codeowner_approval_required")
	}
	if !p.PushMergeRestricted {
		failed = append(failed, "push_merge_restricted")
	}
	if !p.ForcePushDisabled {
		failed = append(failed, "force_push_disabled")
	}
	if !p.PushRulesEnabled {
		failed = append(failed, "push_rules_enabled")
	}
	if !p.MinApprovalsRequired {
		failed = append(failed, "min_approvals_required")
	}
	if !p.AuthorApprovalPrevented {
		failed = append(failed, "author_approval_prevented")
	}
	if !p.CommitterApprovalPrevented {
		failed = append(failed, "committer_approval_prevented")
	}
	if !p.ApprovalsRemovedOnCommit {
		failed = append(failed, "approvals_removed_on_commit")
	}

	return failed
}