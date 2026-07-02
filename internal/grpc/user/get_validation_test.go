package user

import (
	"testing"

	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateReadUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		userId  string
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{
			name:    "valid user_id",
			userId:  "user-123",
			wantErr: false,
		},
		{
			name:    "valid uuid",
			userId:  "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "empty user_id",
			userId:  "",
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "user id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &gen.ReadUserRequest{
				UserId: tt.userId,
			}

			err := validateReadUserRequest(req)

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

func TestValidateReadUserRequestTable(t *testing.T) {
	cases := map[string]struct {
		userId  string
		wantErr bool
	}{
		"valid id":   {"user-1", false},
		"uuid":       {"550e8400-e29b-41d4-a716-446655440000", false},
		"numeric":    {"123", false},
		"empty":      {"", true},
		"whitespace": {"   ", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := validateReadUserRequest(&gen.ReadUserRequest{
				UserId: tc.userId,
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

func BenchmarkValidateReadUserRequest(b *testing.B) {
	req := &gen.ReadUserRequest{UserId: "user-123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateReadUserRequest(req)
	}
}
