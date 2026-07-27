package utils

import (
	"math/rand"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GenerateULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func UserDomainToProto(domainUser domain.User) *gen.User {
	return &gen.User{
		UserId:       domainUser.ID,
		FirstName:    domainUser.FirstName,
		LastName:     domainUser.LastName,
		EmailAddress: domainUser.EmailAddress,
	}
}

func ProjectDomainToProto(domainProject domain.Project) *gen.Project {
	return &gen.Project{
		ProjectId:   domainProject.ID,
		Name:        domainProject.Name,
		Description: domainProject.Description,
	}
}

func IssueDomainToProto(domainIssue domain.Issue) *gen.Issue {
	return &gen.Issue{
		IssueId:     domainIssue.ID,
		CreateDate:  TimeToTimestamp(domainIssue.CreateDate),
		ModifyDate:  TimeToTimestamp(domainIssue.ModifyDate),
		Summary:     domainIssue.Summary,
		Description: domainIssue.Description,
		Status:      IssueDomainStatusToProto(domainIssue.Status),
		Type:        IssueDomainTypeToProto(domainIssue.Type),
		Priority:    IssueDomainPriorityToProto(domainIssue.Priority),
		Resolution:  IssueDomainResolutionToProto(domainIssue.Resolution),
		ProjectId:   domainIssue.ProjectID,
		AssigneeId:  domainIssue.AssigneeID,
	}
}

func IssueProtoToDomain(protoIssue *gen.Issue) *domain.Issue {
	return &domain.Issue{
		ID:          protoIssue.IssueId,
		CreateDate:  TimestampToTime(protoIssue.CreateDate),
		ModifyDate:  TimestampToTime(protoIssue.ModifyDate),
		Summary:     protoIssue.Summary,
		Description: protoIssue.Description,

		Status:     IssueProtoStatusToDomain(protoIssue.Status),
		Resolution: IssueProtoResolutionToDomain(protoIssue.Resolution),
		Type:       IssueProtoTypeToDomain(protoIssue.Type),
		Priority:   IssueProtoPriorityToDomain(protoIssue.Priority),

		ProjectID:  protoIssue.ProjectId,
		AssigneeID: protoIssue.AssigneeId,
	}
}

func TimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
