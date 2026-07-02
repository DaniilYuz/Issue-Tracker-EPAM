package user

import (
	"net/mail"

	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateCreateUserRequest(req *gen.CreateUserRequest) error {
	if req.GetFirstName() == "" ||
		req.GetLastName() == "" ||
		req.GetEmailAddress() == "" {

		return status.Error(
			codes.InvalidArgument,
			"first name, last name and email address are required",
		)
	}

	if _, err := mail.ParseAddress(req.GetEmailAddress()); err != nil {
		return status.Error(
			codes.InvalidArgument,
			"invalid email address format",
		)
	}

	return nil
}
