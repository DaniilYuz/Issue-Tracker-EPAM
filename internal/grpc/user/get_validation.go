package user

import (
	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
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
