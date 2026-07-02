package project

import (
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
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
