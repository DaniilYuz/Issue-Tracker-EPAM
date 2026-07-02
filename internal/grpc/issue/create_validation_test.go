package issue

import (
	"testing"

	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func cloneIssue(req *gen.CreateIssueRequest, mut func(*gen.CreateIssueRequest)) *gen.CreateIssueRequest {
	cp := proto.Clone(req).(*gen.CreateIssueRequest)
	mut(cp)
	return cp
}

func runIssueTests(t *testing.T, tests []struct {
	name    string
	req     *gen.CreateIssueRequest
	wantErr bool
	errCode codes.Code
	errMsg  string
}, fn func(*gen.CreateIssueRequest) error) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fn(tt.req)

			if !tt.wantErr {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)

			assert.Equal(t, tt.errCode, st.Code())
			assert.Equal(t, tt.errMsg, st.Message())
		})
	}
}

func TestValidateCreateRequiredFields(t *testing.T) {
	base := &gen.CreateIssueRequest{
		Summary:     "Test summary",
		Description: "Test description",
		ProjectId:   "project-123",
		Status:      gen.Status_STATUS_NEW,
		Type:        gen.IssueType_ISSUE_TYPE_BUG,
		Priority:    gen.Priority_PRIORITY_MAJOR,
	}

	tests := []struct {
		name    string
		req     *gen.CreateIssueRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{"valid", base, false, 0, ""},

		{"summary empty", cloneIssue(base, func(r *gen.CreateIssueRequest) { r.Summary = "" }), true, codes.InvalidArgument,
			"summary, description, project_id, status, type and priority are required"},

		{"description empty", cloneIssue(base, func(r *gen.CreateIssueRequest) { r.Description = "" }), true, codes.InvalidArgument,
			"summary, description, project_id, status, type and priority are required"},

		{"project empty", cloneIssue(base, func(r *gen.CreateIssueRequest) { r.ProjectId = "" }), true, codes.InvalidArgument,
			"summary, description, project_id, status, type and priority are required"},

		{"status unspecified", cloneIssue(base, func(r *gen.CreateIssueRequest) {
			r.Status = gen.Status_STATUS_UNSPECIFIED
		}), true, codes.InvalidArgument,
			"summary, description, project_id, status, type and priority are required"},

		{"type unspecified", cloneIssue(base, func(r *gen.CreateIssueRequest) {
			r.Type = gen.IssueType_ISSUE_TYPE_UNSPECIFIED
		}), true, codes.InvalidArgument,
			"summary, description, project_id, status, type and priority are required"},

		{"priority unspecified", cloneIssue(base, func(r *gen.CreateIssueRequest) {
			r.Priority = gen.Priority_PRIORITY_UNSPECIFIED
		}), true, codes.InvalidArgument,
			"summary, description, project_id, status, type and priority are required"},
	}

	runIssueTests(t, tests, validateCreateRequiredFields)
}

func TestValidateCreateAllowedStatus(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.CreateIssueRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{"new", &gen.CreateIssueRequest{Status: gen.Status_STATUS_NEW}, false, 0, ""},
		{"assigned", &gen.CreateIssueRequest{Status: gen.Status_STATUS_ASSIGNED}, false, 0, ""},

		{"in progress", &gen.CreateIssueRequest{Status: gen.Status_STATUS_IN_PROGRESS}, true,
			codes.InvalidArgument, "issue can only be created with NEW or ASSIGNED status"},

		{"resolved", &gen.CreateIssueRequest{Status: gen.Status_STATUS_RESOLVED}, true,
			codes.InvalidArgument, "issue can only be created with NEW or ASSIGNED status"},

		{"closed", &gen.CreateIssueRequest{Status: gen.Status_STATUS_CLOSED}, true,
			codes.InvalidArgument, "issue can only be created with NEW or ASSIGNED status"},

		{"unspecified", &gen.CreateIssueRequest{Status: gen.Status_STATUS_UNSPECIFIED}, true,
			codes.InvalidArgument, "issue can only be created with NEW or ASSIGNED status"},
	}

	runIssueTests(t, tests, validateCreateAllowedStatus)
}

func TestValidateCreateAssigneeRules(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.CreateIssueRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{"new no assignee", &gen.CreateIssueRequest{Status: gen.Status_STATUS_NEW, AssigneeId: ""}, false, 0, ""},
		{"assigned with assignee", &gen.CreateIssueRequest{Status: gen.Status_STATUS_ASSIGNED, AssigneeId: "user-1"}, false, 0, ""},

		{"new with assignee", &gen.CreateIssueRequest{Status: gen.Status_STATUS_NEW, AssigneeId: "user-1"}, true,
			codes.InvalidArgument, "NEW issue cannot have assignee"},

		{"assigned without assignee", &gen.CreateIssueRequest{Status: gen.Status_STATUS_ASSIGNED, AssigneeId: ""}, true,
			codes.InvalidArgument, "ASSIGNED issue requires assignee_id"},

		{"other status ok 1", &gen.CreateIssueRequest{Status: gen.Status_STATUS_IN_PROGRESS, AssigneeId: "user-1"}, false, 0, ""},
		{"other status ok 2", &gen.CreateIssueRequest{Status: gen.Status_STATUS_RESOLVED, AssigneeId: ""}, false, 0, ""},
	}

	runIssueTests(t, tests, validateCreateAssigneeRules)
}

func TestValidateCreateResolution(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.CreateIssueRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{"unspecified ok", &gen.CreateIssueRequest{
			Resolution: gen.Resolution_RESOLUTION_RESOLUTION_UNSPECIFIED,
		}, false, 0, ""},

		{"fixed", &gen.CreateIssueRequest{
			Resolution: gen.Resolution_RESOLUTION_FIXED,
		}, true, codes.InvalidArgument, "resolution cannot be set during issue creation"},

		{"wontfix", &gen.CreateIssueRequest{
			Resolution: gen.Resolution_RESOLUTION_WONTFIX,
		}, true, codes.InvalidArgument, "resolution cannot be set during issue creation"},

		{"invalid", &gen.CreateIssueRequest{
			Resolution: gen.Resolution_RESOLUTION_INVALID,
		}, true, codes.InvalidArgument, "resolution cannot be set during issue creation"},

		{"worksforme", &gen.CreateIssueRequest{
			Resolution: gen.Resolution_RESOLUTION_WORKSFORME,
		}, true, codes.InvalidArgument, "resolution cannot be set during issue creation"},
	}

	runIssueTests(t, tests, validateCreateResolution)
}
