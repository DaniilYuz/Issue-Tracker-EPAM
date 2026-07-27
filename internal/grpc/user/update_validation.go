package user

import (
	"net/mail"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateUpdateUserRequest(req *gen.UpdateUserRequest) error {
	user := req.GetUser()

	if user == nil {
		return status.Error(
			codes.InvalidArgument,
			"user entity is required",
		)
	}

	if user.GetUserId() == "" ||
		user.GetFirstName() == "" ||
		user.GetLastName() == "" ||
		user.GetEmailAddress() == "" {

		return status.Error(
			codes.InvalidArgument,
			"user id, first name, last name and email address are required",
		)
	}

	if _, err := mail.ParseAddress(user.GetEmailAddress()); err != nil {
		return status.Error(
			codes.InvalidArgument,
			"invalid email address format",
		)
	}

	return nil
}
