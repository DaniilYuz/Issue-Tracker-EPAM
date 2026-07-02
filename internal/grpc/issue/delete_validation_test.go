package issue

import (
	"testing"

	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateDeleteIssueRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.DeleteIssueRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{
			name: "valid issue id",
			req: &gen.DeleteIssueRequest{
				IssueId: "issue-123",
			},
			wantErr: false,
		},
		{
			name: "valid uuid issue id",
			req: &gen.DeleteIssueRequest{
				IssueId: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: false,
		},
		{
			name: "empty issue id",
			req: &gen.DeleteIssueRequest{
				IssueId: "",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "issue id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeleteIssueRequest(tt.req)

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

func TestValidateDeleteIssueRequestTable(t *testing.T) {
	cases := map[string]struct {
		issueId string
		wantErr bool
	}{
		"valid simple id":    {"issue-1", false},
		"valid numeric id":   {"123", false},
		"valid uuid":         {"550e8400-e29b-41d4-a716-446655440000", false},
		"valid complex id":   {"PROJ-123", false},
		"empty id":           {"", true},
		"whitespace only id": {"   ", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := &gen.DeleteIssueRequest{
				IssueId: tc.issueId,
			}

			err := validateDeleteIssueRequest(req)

			if !tc.wantErr {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
		})
	}
}

func BenchmarkValidateDeleteIssueRequest(b *testing.B) {
	req := &gen.DeleteIssueRequest{IssueId: "issue-123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateDeleteIssueRequest(req)
	}
}

func BenchmarkValidateDeleteIssueRequestEmpty(b *testing.B) {
	req := &gen.DeleteIssueRequest{IssueId: ""}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateDeleteIssueRequest(req)
	}
}
