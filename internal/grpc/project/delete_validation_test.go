package project

import (
	"testing"

	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateDeleteProjectRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.DeleteProjectRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{
			name: "valid request",

			req: &gen.DeleteProjectRequest{
				ProjectId: "project-123",
			},
			wantErr: false,
		},
		{
			name: "empty project id",
			req: &gen.DeleteProjectRequest{
				ProjectId: "",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "project id is required",
		},
		{
			name: "spaces only",
			req: &gen.DeleteProjectRequest{
				ProjectId: "   ",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeleteProjectRequest(tt.req)

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

func TestValidateDeleteProjectRequestTable(t *testing.T) {
	cases := map[string]struct {
		projectId string
		wantErr   bool
	}{
		"valid simple id":    {"project-1", false},
		"valid numeric id":   {"123", false},
		"valid uuid":         {"550e8400-e29b-41d4-a716-446655440000", false},
		"valid complex id":   {"PROJ-123", false},
		"empty id":           {"", true},
		"whitespace only id": {"   ", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := &gen.DeleteProjectRequest{
				ProjectId: tc.projectId,
			}

			err := validateDeleteProjectRequest(req)

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

func BenchmarkValidateDeleteProjectRequest(b *testing.B) {
	req := &gen.DeleteProjectRequest{
		ProjectId: "project-123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateDeleteProjectRequest(req)
	}
}
