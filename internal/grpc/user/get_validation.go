package user

import (
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateReadUserRequest(req *gen.ReadUserRequest) error {
	if req.GetUserId() == "" {
		return status.Error(
			codes.InvalidArgument,
			"user id is required",
		)
	}

	return nil
}
