package models

import (
	"time"
)

type Project struct {
	ProjectID string `json:"project_id" db:"project_id"`

	// GitLab presence checks
	ProjectPresent   bool `json:"project_present" db:"project_present"`
	AppNameSet       bool `json:"app_name_set" db:"app_name_set"`
	MoabIDSet        bool `json:"moab_id_set" db:"moab_id_set"`
	CodeownersExists bool `json:"codeowners_exists" db:"codeowners_exists"`

	// Branch protection checks
	BranchProtectionEnabled   bool `json:"branch_protection_enabled" db:"branch_protection_enabled"`
	CodeownerApprovalRequired bool `json:"codeowner_approval_required" db:"codeowner_approval_required"`
	PushMergeRestricted       bool `json:"push_merge_restricted" db:"push_merge_restricted"`
	ForcePushDisabled         bool `json:"force_push_disabled" db:"force_push_disabled"`

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
