package issue

import (
	"context"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/grpc/common/utils"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/repo"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	gen.UnimplementedIssueServiceServer
	validator IssueValidator
	repo      repo.IssueRepository
}

func NewServer(v IssueValidator, r repo.IssueRepository) *Server {
	return &Server{
		validator: v,
		repo:      r,
	}
}

func (s *Server) CreateIssue(ctx context.Context, req *gen.CreateIssueRequest) (*gen.CreateIssueResponse, error) {
	//use validation func
	if err := s.validator.ValidateCreate(req); err != nil {
		return nil, err
	}

	domainIssue := &domain.Issue{
		ID:          utils.GenerateULID(),
		CreateDate:  time.Now(),
		ModifyDate:  time.Now(),
		Summary:     req.GetSummary(),
		Description: req.GetDescription(),
		Status:      utils.IssueProtoStatusToDomain(req.GetStatus()),
		Type:        utils.IssueProtoTypeToDomain(req.GetType()),
		Priority:    utils.IssueProtoPriorityToDomain(req.GetPriority()),
		ProjectID:   req.GetProjectId(),
		AssigneeID:  req.GetAssigneeId(),
	}

	if err := s.repo.CreateIssue(ctx, domainIssue); err != nil {
		return nil, err
	}

	return &gen.CreateIssueResponse{
		Issue: utils.IssueDomainToProto(*domainIssue),
	}, nil
}

func (s *Server) ReadIssue(ctx context.Context, req *gen.ReadIssueRequest) (*gen.ReadIssueResponse, error) {
	//use validate-func
	if err := s.validator.ValidateGet(req); err != nil {
		return nil, err
	}

	domainIssue, err := s.repo.GetIssueByID(ctx, req.GetIssueId())
	if err != nil {
		return nil, err
	}

	return &gen.ReadIssueResponse{
		Issue: utils.IssueDomainToProto(*domainIssue),
	}, nil
}

func (s *Server) UpdateIssue(ctx context.Context, req *gen.UpdateIssueRequest) (*gen.UpdateIssueResponse, error) {
	//just protection against panic
	if req.GetIssue() == nil {
		return nil, status.Error(codes.InvalidArgument, "issue entity is required")
	}
	newValues := req.GetIssue()

	issue, err := s.repo.GetIssueByID(ctx, newValues.IssueId)
	if err != nil {
		return nil, err
	}
	protoIssue := utils.IssueDomainToProto(*issue)

	if err := s.validator.ValidateUpdate(newValues, protoIssue); err != nil {
		return nil, err
	}

	applyUpdate(protoIssue, newValues)

	updatedDomainIssue := utils.IssueProtoToDomain(protoIssue)
	updatedDomainIssue.ModifyDate = time.Now()

	if err := s.repo.UpdateIssue(ctx, updatedDomainIssue); err != nil {
		return nil, err
	}

	return &gen.UpdateIssueResponse{
		Issue: utils.IssueDomainToProto(*updatedDomainIssue),
	}, nil
}

func (s *Server) ListIssues(ctx context.Context, req *gen.ListIssuesRequest) (*gen.ListIssuesResponse, error) {
	//take all issue from map and put inside of slice
	domainIssues, err := s.repo.ListIssues(ctx)
	if err != nil {
		return nil, err
	}

	protoIssues := make([]*gen.Issue, 0, len(domainIssues))
	for _, domainIssue := range domainIssues {
		protoIssues = append(protoIssues, utils.IssueDomainToProto(*domainIssue))
	}

	return &gen.ListIssuesResponse{
		Issues: protoIssues,
	}, nil
}

func (s *Server) DeleteIssue(ctx context.Context, req *gen.DeleteIssueRequest) (*emptypb.Empty, error) {
	//use validate-func
	if err := s.validator.ValidateDelete(req); err != nil {
		return nil, err
	}

	if err := s.repo.DeleteIssue(ctx, req.GetIssueId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
