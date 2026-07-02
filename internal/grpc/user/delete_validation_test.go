package user

import (
	"testing"

	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateDeleteUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.DeleteUserRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{
			name: "valid user id",
			req: &gen.DeleteUserRequest{
				UserId: "user-123",
			},
			wantErr: false,
		},
		{
			name: "valid uuid",
			req: &gen.DeleteUserRequest{
				UserId: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: false,
		},
		{
			name: "empty user id",
			req: &gen.DeleteUserRequest{
				UserId: "",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "user id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeleteUserRequest(tt.req)

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

func TestValidateDeleteUserRequestTable(t *testing.T) {
	cases := map[string]struct {
		userId  string
		wantErr bool
	}{
		"valid simple id":    {"user-1", false},
		"valid numeric id":   {"123", false},
		"valid uuid":         {"550e8400-e29b-41d4-a716-446655440000", false},
		"valid complex id":   {"USER-123", false},
		"empty id":           {"", true},
		"whitespace only id": {"   ", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := &gen.DeleteUserRequest{
				UserId: tc.userId,
			}

			err := validateDeleteUserRequest(req)

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

func BenchmarkValidateDeleteUserRequest(b *testing.B) {
	req := &gen.DeleteUserRequest{UserId: "user-123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateDeleteUserRequest(req)
	}
}
