package project

import (
	"testing"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateUpdateProjectRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.UpdateProjectRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{
			name: "valid request",
			req: &gen.UpdateProjectRequest{
				Project: &gen.Project{
					ProjectId:   "project-123",
					Name:        "Test Project",
					Description: "Test Description",
				},
			},
			wantErr: false,
		},
		{
			name: "nil project",
			req: &gen.UpdateProjectRequest{
				Project: nil,
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "project entity is required",
		},
		{
			name: "empty project id",
			req: &gen.UpdateProjectRequest{
				Project: &gen.Project{
					ProjectId:   "",
					Name:        "Test",
					Description: "Description",
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "project id, name and description are required",
		},
		{
			name: "empty name",
			req: &gen.UpdateProjectRequest{
				Project: &gen.Project{
					ProjectId:   "project-123",
					Name:        "",
					Description: "Description",
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "project id, name and description are required",
		},
		{
			name: "empty description",
			req: &gen.UpdateProjectRequest{
				Project: &gen.Project{
					ProjectId:   "project-123",
					Name:        "Test",
					Description: "",
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "project id, name and description are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateProjectRequest(tt.req)

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

func TestValidateUpdateProjectRequestTable(t *testing.T) {
	base := &gen.Project{
		ProjectId:   "project-123",
		Name:        "Test Project",
		Description: "Test Description",
	}

	tests := []struct {
		name   string
		modify func(*gen.Project)
		err    bool
	}{
		{"ok", func(p *gen.Project) {}, false},
		{"no id", func(p *gen.Project) { p.ProjectId = "" }, true},
		{"no name", func(p *gen.Project) { p.Name = "" }, true},
		{"no desc", func(p *gen.Project) { p.Description = "" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Клонируем базовый объект
			p := &gen.Project{
				ProjectId:   base.ProjectId,
				Name:        base.Name,
				Description: base.Description,
			}

			// Применяем модификацию
			tt.modify(p)

			req := &gen.UpdateProjectRequest{Project: p}
			err := validateUpdateProjectRequest(req)

			checkGRPCError(t, err, tt.err, codes.InvalidArgument,
				"project id, name and description are required")
		})
	}
}

func BenchmarkValidateUpdateProjectRequest(b *testing.B) {
	req := &gen.UpdateProjectRequest{
		Project: &gen.Project{
			ProjectId:   "project-123",
			Name:        "Test Project",
			Description: "Test Description",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateUpdateProjectRequest(req)
	}
}
