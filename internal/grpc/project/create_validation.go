package project

import (
	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateCreateProjectRequest(req *gen.CreateProjectRequest) error {
	if req.GetName() == "" || req.GetDescription() == "" {
		return status.Error(
			codes.InvalidArgument,
			"name and description are required",
		)
	}

	return nil
}
