package project

import (
	"testing"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateCreateProjectRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.CreateProjectRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{
			name: "valid request",
			req: &gen.CreateProjectRequest{
				Name:        "Test Project",
				Description: "Test description",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: &gen.CreateProjectRequest{
				Name:        "",
				Description: "Test description",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "name and description are required",
		},
		{
			name: "empty description",
			req: &gen.CreateProjectRequest{
				Name:        "Test Project",
				Description: "",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "name and description are required",
		},
		{
			name: "both empty",
			req: &gen.CreateProjectRequest{
				Name:        "",
				Description: "",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "name and description are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateProjectRequest(tt.req)

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

func BenchmarkValidateCreateProjectRequest(b *testing.B) {
	req := &gen.CreateProjectRequest{
		Name:        "Test Project",
		Description: "Test description",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateCreateProjectRequest(req)
	}
}
