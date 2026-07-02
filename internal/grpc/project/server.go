package project

import (
	"context"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"git.epam.com/go-language-global-mentoring-program/internal/grpc/common/utils"
	"git.epam.com/go-language-global-mentoring-program/internal/repo"
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	gen.UnimplementedProjectServiceServer
	validator ProjectValidator
	repo      repo.ProjectRepository
}

func NewServer(v ProjectValidator, r repo.ProjectRepository) *Server {
	return &Server{
		validator: v,
		repo:      r,
	}
}

func (s *Server) CreateProject(ctx context.Context, req *gen.CreateProjectRequest) (*gen.CreateProjectResponse, error) {
	//use validate-func
	if err := s.validator.ValidateCreate(req); err != nil {
		return nil, err
	}

	domainProject := &domain.Project{
		ID:          utils.GenerateULID(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}

	// save our project
	if err := s.repo.CreateProject(ctx, domainProject); err != nil {
		return nil, err
	}

	return &gen.CreateProjectResponse{
		Project: utils.ProjectDomainToProto(*domainProject),
	}, nil
}

func (s *Server) ReadProject(ctx context.Context, req *gen.ReadProjectRequest) (*gen.ReadProjectResponse, error) {
	//use validate-func
	if err := s.validator.ValidateGet(req); err != nil {
		return nil, err
	}

	domainProject, err := s.repo.GetProjectByID(ctx, req.GetProjectId())
	if err != nil {
		return nil, err
	}

	return &gen.ReadProjectResponse{
		Project: utils.ProjectDomainToProto(*domainProject),
	}, nil
}

func (s *Server) UpdateProject(ctx context.Context, req *gen.UpdateProjectRequest) (*gen.UpdateProjectResponse, error) {
	//use validate-func
	if err := s.validator.ValidateUpdate(req); err != nil {
		return nil, err
	}

	if _, err := s.repo.GetProjectByID(ctx, req.GetProject().ProjectId); err != nil {
		return nil, err
	}

	domainProject := &domain.Project{
		ID:          req.Project.GetProjectId(),
		Name:        req.GetProject().GetName(),
		Description: req.GetProject().GetDescription(),
	}

	if err := s.repo.UpdateProject(ctx, domainProject); err != nil {
		return nil, err
	}
	return &gen.UpdateProjectResponse{
		Project: utils.ProjectDomainToProto(*domainProject),
	}, nil
}

func (s *Server) DeleteProject(ctx context.Context, req *gen.DeleteProjectRequest) (*emptypb.Empty, error) {
	//use validate-func
	if err := s.validator.ValidateDelete(req); err != nil {
		return nil, err
	}

	if err := s.repo.DeleteProject(ctx, req.GetProjectId()); err != nil {
		return nil, err
	}

	// return empty gRPC object
	return &emptypb.Empty{}, nil
}

func (s *Server) ListProjects(ctx context.Context, req *gen.ListProjectsRequest) (*gen.ListProjectsResponse, error) {
	//take all users from map and put inside of slice
	domainProjects, err := s.repo.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	protoProjects := make([]*gen.Project, 0, len(domainProjects))
	for _, domainProject := range domainProjects {
		protoProjects = append(protoProjects, utils.ProjectDomainToProto(*domainProject))
	}

	return &gen.ListProjectsResponse{
		Projects: protoProjects,
	}, nil
}
