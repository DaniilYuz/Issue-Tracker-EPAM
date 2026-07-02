package issue

import (
	"testing"

	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestValidateUpdateRequiredFields(t *testing.T) {
	base := &gen.Issue{
		IssueId:     "id",
		Summary:     "s",
		Description: "d",
		ProjectId:   "p",
		Status:      gen.Status_STATUS_NEW,
		Type:        gen.IssueType_ISSUE_TYPE_BUG,
		Priority:    gen.Priority_PRIORITY_MAJOR,
	}

	tests := []struct {
		name   string
		modify func(*gen.Issue)
		err    bool
	}{
		{"ok", func(i *gen.Issue) {}, false},
		{"no id", func(i *gen.Issue) { i.IssueId = "" }, true},
		{"no summary", func(i *gen.Issue) { i.Summary = "" }, true},
		{"no desc", func(i *gen.Issue) { i.Description = "" }, true},
		{"no project", func(i *gen.Issue) { i.ProjectId = "" }, true},
		{"bad status", func(i *gen.Issue) { i.Status = gen.Status_STATUS_UNSPECIFIED }, true},
		{"bad type", func(i *gen.Issue) { i.Type = gen.IssueType_ISSUE_TYPE_UNSPECIFIED }, true},
		{"bad priority", func(i *gen.Issue) { i.Priority = gen.Priority_PRIORITY_UNSPECIFIED }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &gen.Issue{
				IssueId:     base.IssueId,
				Summary:     base.Summary,
				Description: base.Description,
				ProjectId:   base.ProjectId,
				Status:      base.Status,
				Type:        base.Type,
				Priority:    base.Priority,
			}

			tt.modify(i)

			err := validateUpdateRequiredFields(i)

			checkGRPCError(t, err, tt.err, codes.InvalidArgument,
				"issue id, summary, description, project_id, status, type and priority are required")
		})
	}
}

// ---------------- STATUS TRANSITION ----------------

func TestValidateUpdateStatusTransition(t *testing.T) {
	tests := []struct {
		name string
		cur  gen.Status
		next gen.Status
		err  bool
		code codes.Code
		msg  string
	}{
		{"ok same", gen.Status_STATUS_NEW, gen.Status_STATUS_NEW, false, 0, ""},
		{"ok forward", gen.Status_STATUS_NEW, gen.Status_STATUS_ASSIGNED, false, 0, ""},
		{"bad backward", gen.Status_STATUS_ASSIGNED, gen.Status_STATUS_NEW, true, codes.FailedPrecondition, "invalid status transition"},
		{"terminal", gen.Status_STATUS_RESOLVED, gen.Status_STATUS_NEW, true, codes.FailedPrecondition, "issue is in terminal status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateStatusTransition(tt.cur, tt.next)
			checkGRPCError(t, err, tt.err, tt.code, tt.msg)
		})
	}
}

// ---------------- ASSIGNEE ----------------

func TestValidateUpdateAssigneeRules(t *testing.T) {
	tests := []struct {
		name string
		i    *gen.Issue
		err  bool
	}{
		{"ok new empty", &gen.Issue{Status: gen.Status_STATUS_NEW}, false},
		{"ok assigned", &gen.Issue{Status: gen.Status_STATUS_ASSIGNED, AssigneeId: "u"}, false},
		{"bad assigned empty", &gen.Issue{Status: gen.Status_STATUS_ASSIGNED}, true},
		{"bad resolved empty", &gen.Issue{Status: gen.Status_STATUS_RESOLVED}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateAssigneeRules(tt.i)
			checkGRPCError(t, err, tt.err, codes.InvalidArgument,
				"assignee_id is required for ASSIGNED, IN_PROGRESS, RESOLVED and CLOSED statuses")
		})
	}
}

// ---------------- RESOLUTION ----------------

func TestValidateUpdateResolutionRules(t *testing.T) {
	tests := []struct {
		name string
		i    *gen.Issue
		err  bool
		msg  string
	}{
		{"ok new", &gen.Issue{Status: gen.Status_STATUS_NEW}, false, ""},
		{"ok resolved", &gen.Issue{
			Status:     gen.Status_STATUS_RESOLVED,
			Resolution: gen.Resolution_RESOLUTION_FIXED,
		}, false, ""},
		{"bad missing", &gen.Issue{
			Status: gen.Status_STATUS_RESOLVED,
		}, true, "resolution is required"},
		{"bad set early", &gen.Issue{
			Status:     gen.Status_STATUS_NEW,
			Resolution: gen.Resolution_RESOLUTION_FIXED,
		}, true, "resolution must not be set"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateResolutionRules(tt.i)
			checkGRPCError(t, err, tt.err, codes.InvalidArgument, tt.msg)
		})
	}
}

// ----------------APPLY UPDATE ----------------

func TestApplyUpdate(t *testing.T) {
	base := &gen.Issue{
		IssueId:    "id",
		Summary:    "old",
		Status:     gen.Status_STATUS_NEW,
		Type:       gen.IssueType_ISSUE_TYPE_BUG,
		Priority:   gen.Priority_PRIORITY_MAJOR,
		ProjectId:  "p",
		CreateDate: timestamppb.Now(),
	}

	update := &gen.Issue{
		Summary:    "new",
		Status:     gen.Status_STATUS_ASSIGNED,
		AssigneeId: "u",
	}

	i := &gen.Issue{
		IssueId:    base.IssueId,
		Summary:    base.Summary,
		Status:     base.Status,
		Type:       base.Type,
		Priority:   base.Priority,
		ProjectId:  base.ProjectId,
		CreateDate: base.CreateDate,
	}

	applyUpdate(i, update)

	assert.Equal(t, "new", i.Summary)
	assert.Equal(t, "u", i.AssigneeId)
	assert.Equal(t, base.IssueId, i.IssueId)
}

// ---------------- BENCHMARKS ----------------

func BenchmarkValidateUpdateRequiredFields(b *testing.B) {
	i := &gen.Issue{
		IssueId: "id", Summary: "s", Description: "d",
		ProjectId: "p", Status: gen.Status_STATUS_NEW,
		Type: gen.IssueType_ISSUE_TYPE_BUG, Priority: gen.Priority_PRIORITY_MAJOR,
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		_ = validateUpdateRequiredFields(i)
	}
}
