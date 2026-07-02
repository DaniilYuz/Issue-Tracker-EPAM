package issue

import (
	"testing"

	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateReadIssueRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.ReadIssueRequest
		wantErr bool
	}{
		{
			name: "valid issue id",
			req: &gen.ReadIssueRequest{
				IssueId: "issue-123",
			},
			wantErr: false,
		},
		{
			name: "valid uuid",
			req: &gen.ReadIssueRequest{
				IssueId: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: false,
		},
		{
			name: "empty issue id",
			req: &gen.ReadIssueRequest{
				IssueId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReadIssueRequest(tt.req)

			if !tt.wantErr {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
			st, ok := status.FromError(err)
			require.True(t, ok)

			assert.Equal(t, codes.InvalidArgument, st.Code())
			assert.Equal(t, "issue id is required", st.Message())
		})
	}
}

func TestValidateReadIssueRequestTable(t *testing.T) {
	cases := map[string]struct {
		issueId string
		wantErr bool
	}{
		"valid id":   {"issue-1", false},
		"uuid":       {"550e8400-e29b-41d4-a716-446655440000", false},
		"numeric":    {"123", false},
		"empty":      {"", true},
		"whitespace": {"   ", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := validateReadIssueRequest(&gen.ReadIssueRequest{
				IssueId: tc.issueId,
			})

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

func BenchmarkValidateReadIssueRequest(b *testing.B) {
	req := &gen.ReadIssueRequest{IssueId: "issue-123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateReadIssueRequest(req)
	}
}

func BenchmarkValidateReadIssueRequestEmpty(b *testing.B) {
	req := &gen.ReadIssueRequest{IssueId: ""}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateReadIssueRequest(req)
	}
}
