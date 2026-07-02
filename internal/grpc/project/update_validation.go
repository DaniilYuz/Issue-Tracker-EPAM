package project

import (
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateUpdateProjectRequest(req *gen.UpdateProjectRequest) error {
	project := req.GetProject()

	if project == nil {
		return status.Error(
			codes.InvalidArgument,
			"project entity is required",
		)
	}

	if project.GetProjectId() == "" ||
		project.GetName() == "" ||
		project.GetDescription() == "" {

		return status.Error(
			codes.InvalidArgument,
			"project id, name and description are required",
		)
	}

	return nil
}
