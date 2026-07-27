package issue

import (
	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateReadIssueRequest(req *gen.ReadIssueRequest) error {
	if req.GetIssueId() == "" {
		return status.Error(
			codes.InvalidArgument,
			"issue id is required",
		)
	}

	return nil
}
