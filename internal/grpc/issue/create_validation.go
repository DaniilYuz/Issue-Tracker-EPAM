package issue

import (
	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateCreateRequiredFields(req *gen.CreateIssueRequest) error {
	if req.GetSummary() == "" ||
		req.GetDescription() == "" ||
		req.GetProjectId() == "" ||
		req.GetStatus() == gen.Status_STATUS_UNSPECIFIED ||
		req.GetType() == gen.IssueType_ISSUE_TYPE_UNSPECIFIED ||
		req.GetPriority() == gen.Priority_PRIORITY_UNSPECIFIED {

		return status.Error(
			codes.InvalidArgument,
			"summary, description, project_id, status, type and priority are required",
		)
	}
	return nil
}

func validateCreateAllowedStatus(req *gen.CreateIssueRequest) error {
	st := req.GetStatus()

	if st != gen.Status_STATUS_NEW &&
		st != gen.Status_STATUS_ASSIGNED {

		return status.Error(
			codes.InvalidArgument,
			"issue can only be created with NEW or ASSIGNED status",
		)
	}

	return nil
}

func validateCreateAssigneeRules(req *gen.CreateIssueRequest) error {
	st := req.GetStatus()
	assignee := req.GetAssigneeId()

	if st == gen.Status_STATUS_NEW && assignee != "" {
		return status.Error(
			codes.InvalidArgument,
			"NEW issue cannot have assignee",
		)
	}

	if st == gen.Status_STATUS_ASSIGNED && assignee == "" {
		return status.Error(
			codes.InvalidArgument,
			"ASSIGNED issue requires assignee_id",
		)
	}

	return nil
}

func validateCreateResolution(req *gen.CreateIssueRequest) error {
	if req.GetResolution() != gen.Resolution_RESOLUTION_RESOLUTION_UNSPECIFIED {
		return status.Error(
			codes.InvalidArgument,
			"resolution cannot be set during issue creation",
		)
	}

	return nil
}
