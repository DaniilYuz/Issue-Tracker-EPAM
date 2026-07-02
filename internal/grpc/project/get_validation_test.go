package project

import (
	"testing"

	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateReadProjectRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.ReadProjectRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{
			name: "valid request",
			req: &gen.ReadProjectRequest{
				ProjectId: "project-123",
			},
			wantErr: false,
		},
		{
			name: "empty project id",
			req: &gen.ReadProjectRequest{
				ProjectId: "",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "project id is required",
		},
		{
			name: "spaces only",
			req: &gen.ReadProjectRequest{
				ProjectId: "   ",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReadProjectRequest(tt.req)

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

func TestValidateReadProjectRequestTable(t *testing.T) {
	cases := map[string]struct {
		projectId string
		wantErr   bool
	}{
		"valid id":   {"project-1", false},
		"uuid":       {"550e8400-e29b-41d4-a716-446655440000", false},
		"numeric":    {"123", false},
		"empty":      {"", true},
		"whitespace": {"   ", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := validateReadProjectRequest(&gen.ReadProjectRequest{
				ProjectId: tc.projectId,
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

func BenchmarkValidateReadProjectRequest(b *testing.B) {
	req := &gen.ReadProjectRequest{ProjectId: "project-123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateReadProjectRequest(req)
	}
}
