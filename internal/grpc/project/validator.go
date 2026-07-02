package project

import (
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
)

type ProjectValidator interface {
	ValidateCreate(req *gen.CreateProjectRequest) error
	ValidateUpdate(req *gen.UpdateProjectRequest) error
	ValidateGet(req *gen.ReadProjectRequest) error
	ValidateDelete(req *gen.DeleteProjectRequest) error
}

type AppValidator struct{}

func (v *AppValidator) ValidateCreate(req *gen.CreateProjectRequest) error {
	return validateCreateProjectRequest(req)
}

func (v *AppValidator) ValidateUpdate(req *gen.UpdateProjectRequest) error {
	return validateUpdateProjectRequest(req)
}

func (v *AppValidator) ValidateGet(req *gen.ReadProjectRequest) error {
	return validateReadProjectRequest(req)
}

func (v *AppValidator) ValidateDelete(req *gen.DeleteProjectRequest) error {
	return validateDeleteProjectRequest(req)
}
