package project

import (
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateDeleteProjectRequest(req *gen.DeleteProjectRequest) error {
	if req.GetProjectId() == "" {
		return status.Error(
			codes.InvalidArgument,
			"project id is required",
		)
	}

	return nil
}
