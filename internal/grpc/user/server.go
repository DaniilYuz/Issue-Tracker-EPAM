package user

import (
	"context"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"git.epam.com/go-language-global-mentoring-program/internal/grpc/common/utils"
	"git.epam.com/go-language-global-mentoring-program/internal/repo"
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	gen.UnimplementedUserServiceServer
	validator UserValidator
	repo      repo.UserRepository
}

func NewServer(v UserValidator, r repo.UserRepository) *Server {
	return &Server{
		validator: v,
		repo:      r,
	}
}

func (s *Server) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	//use validate-func
	if err := s.validator.ValidateCreate(req); err != nil {
		return nil, err
	}

	domainUser := &domain.User{
		ID:           utils.GenerateULID(),
		FirstName:    req.GetFirstName(),
		LastName:     req.GetLastName(),
		EmailAddress: req.GetEmailAddress(),
	}

	if err := s.repo.CreateUser(ctx, domainUser); err != nil {
		return nil, err
	}

	return &gen.CreateUserResponse{
		User: utils.UserDomainToProto(*domainUser),
	}, nil
}

func (s *Server) ReadUser(ctx context.Context, req *gen.ReadUserRequest) (*gen.ReadUserResponse, error) {
	//use validate-func
	if err := s.validator.ValidateGet(req); err != nil {
		return nil, err
	}

	domainUser, err := s.repo.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &gen.ReadUserResponse{
		User: utils.UserDomainToProto(*domainUser),
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *gen.UpdateUserRequest) (*gen.UpdateUserResponse, error) {
	//use validate-func
	if err := s.validator.ValidateUpdate(req); err != nil {
		return nil, err
	}

	//try to find our user in our "db"
	if _, err := s.repo.GetUserByID(ctx, req.User.GetUserId()); err != nil {
		return nil, err
	}

	domainUser := &domain.User{
		ID:           req.GetUser().GetUserId(),
		FirstName:    req.GetUser().GetFirstName(),
		LastName:     req.GetUser().GetLastName(),
		EmailAddress: req.GetUser().GetEmailAddress(),
	}

	if err := s.repo.UpdateUser(ctx, domainUser); err != nil {
		return nil, err
	}

	return &gen.UpdateUserResponse{
		User: utils.UserDomainToProto(*domainUser),
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *gen.DeleteUserRequest) (*emptypb.Empty, error) {
	//use validate-func
	if err := s.validator.ValidateDelete(req); err != nil {
		return nil, err
	}

	if err := s.repo.DeleteUser(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	//return empty gRPC object
	return &emptypb.Empty{}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *gen.ListUsersRequest) (*gen.ListUsersResponse, error) {
	//take all users from map and put inside of slice
	domainUsers, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*gen.User, 0, len(domainUsers))
	for _, domainUser := range domainUsers {
		protoUsers = append(protoUsers, utils.UserDomainToProto(*domainUser))
	}

	return &gen.ListUsersResponse{
		Users: protoUsers,
	}, nil
}
