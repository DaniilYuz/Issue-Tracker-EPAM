package user

import (
	"testing"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateCreateUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.CreateUserRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{
			name: "valid request",
			req: &gen.CreateUserRequest{
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name: "empty first name",
			req: &gen.CreateUserRequest{
				FirstName:    "",
				LastName:     "Doe",
				EmailAddress: "john.doe@example.com",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "first name, last name and email address are required",
		},
		{
			name: "empty last name",
			req: &gen.CreateUserRequest{
				FirstName:    "John",
				LastName:     "",
				EmailAddress: "john.doe@example.com",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "first name, last name and email address are required",
		},
		{
			name: "empty email",
			req: &gen.CreateUserRequest{
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: "",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "first name, last name and email address are required",
		},
		{
			name: "invalid email format",
			req: &gen.CreateUserRequest{
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: "invalid-email",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "invalid email address format",
		},
		{
			name: "email without @",
			req: &gen.CreateUserRequest{
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: "john.doe.example.com",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "invalid email address format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateUserRequest(tt.req)

			if !tt.wantErr {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
			st, ok := status.FromError(err)
			require.True(t, ok)

			assert.Equal(t,
				tt.errCode, st.Code())
			assert.Equal(t, tt.errMsg, st.Message())
		})
	}
}

func TestValidateCreateUserRequestValidEmails(t *testing.T) {
	validEmails := []string{
		"user@example.com",
		"user.name+tag@sub.example.co.uk",
		"john.doe@company.com",
		"test123@domain.org",
	}

	for _, email := range validEmails {
		t.Run(email, func(t *testing.T) {
			req := &gen.CreateUserRequest{
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: email,
			}

			err := validateCreateUserRequest(req)
			assert.NoError(t, err)
		})
	}
}

func BenchmarkValidateCreateUserRequest(b *testing.B) {
	req := &gen.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateCreateUserRequest(req)
	}
}
