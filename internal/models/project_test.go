package models

import (
	"testing"
	"time"
)

func TestProject_IsProductionReady(t *testing.T) {
	tests := []struct {
		name    string
		project Project
		want    bool
	}{
		{
			name: "all checks pass",
			project: Project{
				ProjectID:                  "12345",
				ProjectPresent:             true,
				AppNameSet:                 true,
				MoabIDSet:                  true,
				CodeownersExists:           true,
				BranchProtectionEnabled:    true,
				CodeownerApprovalRequired:  true,
				PushMergeRestricted:        true,
				ForcePushDisabled:          true,
				PushRulesEnabled:           true,
				MinApprovalsRequired:       true,
				AuthorApprovalPrevented:    true,
				CommitterApprovalPrevented: true,
				ApprovalsRemovedOnCommit:   true,
			},
			want: true,
		},
		{
			name: "missing app name",
			project: Project{
				ProjectID:                  "12345",
				ProjectPresent:             true,
				AppNameSet:                 false, // This check fails
				MoabIDSet:                  true,
				CodeownersExists:           true,
				BranchProtectionEnabled:    true,
				CodeownerApprovalRequired:  true,
				PushMergeRestricted:        true,
				ForcePushDisabled:          true,
				PushRulesEnabled:           true,
				MinApprovalsRequired:       true,
				AuthorApprovalPrevented:    true,
				CommitterApprovalPrevented: true,
				ApprovalsRemovedOnCommit:   true,
			},
			want: false,
		},
		{
			name: "multiple failures",
			project: Project{
				ProjectID:                  "12345",
				ProjectPresent:             true,
				AppNameSet:                 false, // Fails
				MoabIDSet:                  false, // Fails
				CodeownersExists:           true,
				BranchProtectionEnabled:    false, // Fails
				CodeownerApprovalRequired:  true,
				PushMergeRestricted:        true,
				ForcePushDisabled:          true,
				PushRulesEnabled:           true,
				MinApprovalsRequired:       true,
				AuthorApprovalPrevented:    true,
				CommitterApprovalPrevented: true,
				ApprovalsRemovedOnCommit:   true,
			},
			want: false,
		},
		{
			name:    "empty project",
			project: Project{},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.project.IsProductionReady(); got != tt.want {
				t.Errorf("IsProductionReady() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProject_FailedChecks(t *testing.T) {
	tests := []struct {
		name    string
		project Project
		want    []string
	}{
		{
			name: "no failures",
			project: Project{
				ProjectPresent:             true,
				AppNameSet:                 true,
				MoabIDSet:                  true,
				CodeownersExists:           true,
				BranchProtectionEnabled:    true,
				CodeownerApprovalRequired:  true,
				PushMergeRestricted:        true,
				ForcePushDisabled:          true,
				PushRulesEnabled:           true,
				MinApprovalsRequired:       true,
				AuthorApprovalPrevented:    true,
				CommitterApprovalPrevented: true,
				ApprovalsRemovedOnCommit:   true,
			},
			want: []string{},
		},
		{
			name: "single failure",
			project: Project{
				ProjectPresent:             true,
				AppNameSet:                 false, // This fails
				MoabIDSet:                  true,
				CodeownersExists:           true,
				BranchProtectionEnabled:    true,
				CodeownerApprovalRequired:  true,
				PushMergeRestricted:        true,
				ForcePushDisabled:          true,
				PushRulesEnabled:           true,
				MinApprovalsRequired:       true,
				AuthorApprovalPrevented:    true,
				CommitterApprovalPrevented: true,
				ApprovalsRemovedOnCommit:   true,
			},
			want: []string{"app_name_set"},
		},
		{
			name: "multiple failures",
			project: Project{
				ProjectPresent:             false, // Fails
				AppNameSet:                 false, // Fails
				MoabIDSet:                  false, // Fails
				CodeownersExists:           true,
				BranchProtectionEnabled:    true,
				CodeownerApprovalRequired:  true,
				PushMergeRestricted:        true,
				ForcePushDisabled:          true,
				PushRulesEnabled:           true,
				MinApprovalsRequired:       true,
				AuthorApprovalPrevented:    true,
				CommitterApprovalPrevented: true,
				ApprovalsRemovedOnCommit:   true,
			},
			want: []string{"project_present", "app_name_set", "moab_id_set"},
		},
		{
			name:    "all failures",
			project: Project{}, // All fields are false by default
			want: []string{
				"project_present",
				"app_name_set",
				"moab_id_set",
				"codeowners_exists",
				"branch_protection_enabled",
				"codeowner_approval_required",
				"push_merge_restricted",
				"force_push_disabled",
				"push_rules_enabled",
				"min_approvals_required",
				"author_approval_prevented",
				"committer_approval_prevented",
				"approvals_removed_on_commit",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.project.FailedChecks()
			
			// Check length
			if len(got) != len(tt.want) {
				t.Errorf("FailedChecks() returned %d items, want %d", len(got), len(tt.want))
				t.Errorf("Got: %v", got)
				t.Errorf("Want: %v", tt.want)
				return
			}
			
			// Check each item
			for i, check := range got {
				if check != tt.want[i] {
					t.Errorf("FailedChecks()[%d] = %v, want %v", i, check, tt.want[i])
				}
			}
		})
	}
}

// BenchmarkIsProductionReady benchmarks the IsProductionReady method
func BenchmarkIsProductionReady(b *testing.B) {
	project := Project{
		ProjectPresent:             true,
		AppNameSet:                 true,
		MoabIDSet:                  true,
		CodeownersExists:           true,
		BranchProtectionEnabled:    true,
		CodeownerApprovalRequired:  true,
		PushMergeRestricted:        true,
		ForcePushDisabled:          true,
		PushRulesEnabled:           true,
		MinApprovalsRequired:       true,
		AuthorApprovalPrevented:    true,
		CommitterApprovalPrevented: true,
		ApprovalsRemovedOnCommit:   true,
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}

	for i := 0; i < b.N; i++ {
		_ = project.IsProductionReady()
	}
}