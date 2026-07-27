package issue

import (
	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
)

type IssueValidator interface {
	ValidateCreate(req *gen.CreateIssueRequest) error
	ValidateUpdate(req, issue *gen.Issue) error
	ValidateGet(req *gen.ReadIssueRequest) error
	ValidateDelete(req *gen.DeleteIssueRequest) error
}

type AppValidator struct{}

func (v *AppValidator) ValidateCreate(req *gen.CreateIssueRequest) error {
	if err := validateCreateRequiredFields(req); err != nil {
		return err
	}
	if err := validateCreateAllowedStatus(req); err != nil {
		return err
	}
	if err := validateCreateAssigneeRules(req); err != nil {
		return err
	}
	if err := validateCreateResolution(req); err != nil {
		return err
	}

	return nil
}

func (v *AppValidator) ValidateUpdate(req, issue *gen.Issue) error {
	if err := validateUpdateRequiredFields(req); err != nil {
		return err
	}

	if err := validateUpdateStatusTransition(issue.Status, req.Status); err != nil {
		return err
	}

	if err := validateUpdateAssigneeRules(req); err != nil {
		return err
	}

	if err := validateUpdateResolutionRules(req); err != nil {
		return err
	}

	return nil
}

func (v *AppValidator) ValidateGet(req *gen.ReadIssueRequest) error {
	return validateReadIssueRequest(req)
}

func (v *AppValidator) ValidateDelete(req *gen.DeleteIssueRequest) error {
	return validateDeleteIssueRequest(req)
}
