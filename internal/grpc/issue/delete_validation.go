package issue

import (
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateDeleteIssueRequest(req *gen.DeleteIssueRequest) error {
	if req.GetIssueId() == "" {
		return status.Error(
			codes.InvalidArgument,
			"issue id is required",
		)
	}

	return nil
}
