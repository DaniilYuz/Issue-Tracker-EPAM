package user

import "github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"

type UserValidator interface {
	ValidateCreate(req *gen.CreateUserRequest) error
	ValidateUpdate(req *gen.UpdateUserRequest) error
	ValidateGet(req *gen.ReadUserRequest) error
	ValidateDelete(req *gen.DeleteUserRequest) error
}

type AppValidator struct{}

func (v *AppValidator) ValidateCreate(req *gen.CreateUserRequest) error {
	return validateCreateUserRequest(req)
}

func (v *AppValidator) ValidateUpdate(req *gen.UpdateUserRequest) error {
	return validateUpdateUserRequest(req)
}

func (v *AppValidator) ValidateGet(req *gen.ReadUserRequest) error {
	return validateReadUserRequest(req)
}

func (v *AppValidator) ValidateDelete(req *gen.DeleteUserRequest) error {
	return validateDeleteUserRequest(req)
}
