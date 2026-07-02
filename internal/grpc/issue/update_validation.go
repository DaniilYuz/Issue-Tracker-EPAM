package issue

import (
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// map for check our workflow transitions
var statusOrder = map[gen.Status]int{
	gen.Status_STATUS_NEW:         0,
	gen.Status_STATUS_ASSIGNED:    1,
	gen.Status_STATUS_IN_PROGRESS: 2,
	gen.Status_STATUS_RESOLVED:    3,
	gen.Status_STATUS_CLOSED:      3,
}

func validateUpdateRequiredFields(newValues *gen.Issue) error {
	if newValues.GetIssueId() == "" ||
		newValues.GetSummary() == "" ||
		newValues.GetDescription() == "" ||
		newValues.GetProjectId() == "" ||
		newValues.GetStatus() == gen.Status_STATUS_UNSPECIFIED ||
		newValues.GetType() == gen.IssueType_ISSUE_TYPE_UNSPECIFIED ||
		newValues.GetPriority() == gen.Priority_PRIORITY_UNSPECIFIED {

		return status.Error(codes.InvalidArgument,
			"issue id, summary, description, project_id, status, type and priority are required")
	}
	return nil
}

func validateUpdateStatusTransition(current, next gen.Status) error {
	if current == next {
		return nil
	}

	currentWeight, ok1 := statusOrder[current]
	nextWeight, ok2 := statusOrder[next]

	if !ok1 || !ok2 {
		return status.Error(codes.FailedPrecondition, "invalid status transition")
	}

	if currentWeight == 3 {
		return status.Error(codes.FailedPrecondition, "issue is in terminal status")
	}

	if nextWeight <= currentWeight {
		return status.Error(codes.FailedPrecondition, "invalid status transition")
	}

	return nil
}

func validateUpdateAssigneeRules(newValues *gen.Issue) error {
	st := newValues.GetStatus()

	if (st == gen.Status_STATUS_ASSIGNED ||
		st == gen.Status_STATUS_IN_PROGRESS ||
		st == gen.Status_STATUS_RESOLVED ||
		st == gen.Status_STATUS_CLOSED) &&
		newValues.GetAssigneeId() == "" {

		return status.Error(
			codes.InvalidArgument,
			"assignee_id is required for ASSIGNED, IN_PROGRESS, RESOLVED and CLOSED statuses",
		)
	}

	return nil
}

func validateUpdateResolutionRules(newValues *gen.Issue) error {
	st := newValues.GetStatus()

	needsResolution :=
		st == gen.Status_STATUS_RESOLVED ||
			st == gen.Status_STATUS_CLOSED

	if needsResolution && newValues.GetResolution() == gen.Resolution_RESOLUTION_RESOLUTION_UNSPECIFIED {
		return status.Error(codes.InvalidArgument, "resolution is required")
	}

	if !needsResolution && newValues.GetResolution() != gen.Resolution_RESOLUTION_RESOLUTION_UNSPECIFIED {
		return status.Error(codes.InvalidArgument, "resolution must not be set")
	}

	return nil
}

func applyUpdate(issue *gen.Issue, newValues *gen.Issue) {
	issue.Summary = newValues.Summary
	issue.Description = newValues.Description
	issue.Status = newValues.Status
	issue.Type = newValues.Type
	issue.Priority = newValues.Priority
	issue.ProjectId = newValues.ProjectId
	issue.ModifyDate = timestamppb.Now()

	if newValues.AssigneeId != "" {
		issue.AssigneeId = newValues.AssigneeId
	}

	if newValues.Status == gen.Status_STATUS_RESOLVED ||
		newValues.Status == gen.Status_STATUS_CLOSED {
		issue.Resolution = newValues.Resolution
	}
}
