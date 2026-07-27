package user

import (
	"testing"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateUpdateUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *gen.UpdateUserRequest
		wantErr bool
		errCode codes.Code
		errMsg  string
	}{
		{
			name: "valid request",
			req: &gen.UpdateUserRequest{
				User: &gen.User{
					UserId:       "user-123",
					FirstName:    "John",
					LastName:     "Doe",
					EmailAddress: "john.doe@example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "nil user",
			req: &gen.UpdateUserRequest{
				User: nil,
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "user entity is required",
		},
		{
			name: "empty required fields",
			req: &gen.UpdateUserRequest{
				User: &gen.User{
					UserId:       "",
					FirstName:    "",
					LastName:     "",
					EmailAddress: "",
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "user id, first name, last name and email address are required",
		},
		{
			name: "invalid email",
			req: &gen.UpdateUserRequest{
				User: &gen.User{
					UserId:       "user-123",
					FirstName:    "John",
					LastName:     "Doe",
					EmailAddress: "invalid-email",
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
			errMsg:  "invalid email address format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateUserRequest(tt.req)

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

func TestValidateUpdateUserRequestValidEmails(t *testing.T) {
	emails := []string{
		"user@example.com",
		"user.name+tag@sub.example.co.uk",
		"john.doe@company.com",
	}

	for _, email := range emails {
		t.Run(email, func(t *testing.T) {
			req := &gen.UpdateUserRequest{
				User: &gen.User{
					UserId:       "user-123",
					FirstName:    "John",
					LastName:     "Doe",
					EmailAddress: email,
				},
			}

			err := validateUpdateUserRequest(req)
			assert.NoError(t, err)
		})
	}
}

func TestValidateUpdateUserRequestInvalidEmails(t *testing.T) {
	emails := []string{
		"plainaddress",
		"@missing.com",
		"no-at-symbol.com",
	}

	for _, email := range emails {
		t.Run(email, func(t *testing.T) {
			req := &gen.UpdateUserRequest{
				User: &gen.User{
					UserId:       "user-123",
					FirstName:    "John",
					LastName:     "Doe",
					EmailAddress: email,
				},
			}

			err := validateUpdateUserRequest(req)
			require.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
		})
	}
}

func TestValidateUpdateUserRequestTable(t *testing.T) {
	base := &gen.User{
		UserId:       "user-123",
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
	}

	tests := []struct {
		name   string
		modify func(*gen.User)
		err    bool
	}{
		{"ok", func(u *gen.User) {}, false},
		{"no id", func(u *gen.User) { u.UserId = "" }, true},
		{"no first name", func(u *gen.User) { u.FirstName = "" }, true},
		{"no last name", func(u *gen.User) { u.LastName = "" }, true},
		{"no email", func(u *gen.User) { u.EmailAddress = "" }, true},
		{"invalid email", func(u *gen.User) { u.EmailAddress = "invalid" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Клонируем базовый объект
			u := &gen.User{
				UserId:       base.UserId,
				FirstName:    base.FirstName,
				LastName:     base.LastName,
				EmailAddress: base.EmailAddress,
			}

			// Применяем модификацию
			tt.modify(u)

			req := &gen.UpdateUserRequest{User: u}
			err := validateUpdateUserRequest(req)

			if !tt.err {
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

func BenchmarkValidateUpdateUserRequest(b *testing.B) {
	req := &gen.UpdateUserRequest{
		User: &gen.User{
			UserId:       "user-123",
			FirstName:    "John",
			LastName:     "Doe",
			EmailAddress: "john.doe@example.com",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateUpdateUserRequest(req)
	}
}
