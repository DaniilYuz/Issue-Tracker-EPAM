package utils

import (
	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
)

// simple mappers. Try to resolve data type dissmiss problem (int in proto and string in domain)
func IssueProtoStatusToDomain(s gen.Status) domain.IssueStatus {
	switch s {
	case gen.Status_STATUS_NEW:
		return domain.IssueStatusNew
	case gen.Status_STATUS_ASSIGNED:
		return domain.IssueStatusAssigned
	case gen.Status_STATUS_IN_PROGRESS:
		return domain.IssueStatusInProgress
	case gen.Status_STATUS_RESOLVED:
		return domain.IssueStatusResolved
	case gen.Status_STATUS_CLOSED:
		return domain.IssueStatusClosed
	case gen.Status_STATUS_REOPENED:
		return domain.IssueStatusReopened
	default:
		return domain.IssueStatusNew
	}
}

func IssueProtoResolutionToDomain(r gen.Resolution) domain.IssueResolution {
	switch r {
	case gen.Resolution_RESOLUTION_FIXED:
		return domain.ResolutionFixed
	case gen.Resolution_RESOLUTION_INVALID:
		return domain.ResolutionInvalid
	case gen.Resolution_RESOLUTION_WONTFIX:
		return domain.ResolutionWontFix
	case gen.Resolution_RESOLUTION_WORKSFORME:
		return domain.ResolutionWorksForMe
	case gen.Resolution_RESOLUTION_RESOLUTION_UNSPECIFIED:
		return domain.ResolutionUnspecified
	default:
		return domain.ResolutionUnspecified
	}
}

func IssueProtoTypeToDomain(t gen.IssueType) domain.IssueType {
	switch t {
	case gen.IssueType_ISSUE_TYPE_COSMETIC:
		return domain.IssueTypeCosmetic
	case gen.IssueType_ISSUE_TYPE_BUG:
		return domain.IssueTypeBug
	case gen.IssueType_ISSUE_TYPE_FEATURE:
		return domain.IssueTypeFeature
	case gen.IssueType_ISSUE_TYPE_PERFORMANCE:
		return domain.IssueTypePerformance
	default:
		return domain.IssueTypeBug
	}
}

func IssueProtoPriorityToDomain(p gen.Priority) domain.Priority {
	switch p {
	case gen.Priority_PRIORITY_CRITICAL:
		return domain.PriorityCritical
	case gen.Priority_PRIORITY_MAJOR:
		return domain.PriorityMajor
	case gen.Priority_PRIORITY_IMPORTANT:
		return domain.PriorityImportant
	case gen.Priority_PRIORITY_MINOR:
		return domain.PriorityMinor
	default:
		return domain.PriorityMajor
	}
}

func IssueDomainStatusToProto(s domain.IssueStatus) gen.Status {
	switch s {
	case domain.IssueStatusNew:
		return gen.Status_STATUS_NEW
	case domain.IssueStatusAssigned:
		return gen.Status_STATUS_ASSIGNED
	case domain.IssueStatusInProgress:
		return gen.Status_STATUS_IN_PROGRESS
	case domain.IssueStatusResolved:
		return gen.Status_STATUS_RESOLVED
	case domain.IssueStatusClosed:
		return gen.Status_STATUS_CLOSED
	case domain.IssueStatusReopened:
		return gen.Status_STATUS_REOPENED
	default:
		return gen.Status_STATUS_UNSPECIFIED
	}
}

func IssueDomainResolutionToProto(r domain.IssueResolution) gen.Resolution {
	switch r {
	case domain.ResolutionFixed:
		return gen.Resolution_RESOLUTION_FIXED
	case domain.ResolutionInvalid:
		return gen.Resolution_RESOLUTION_INVALID
	case domain.ResolutionWontFix:
		return gen.Resolution_RESOLUTION_WONTFIX
	case domain.ResolutionWorksForMe:
		return gen.Resolution_RESOLUTION_WORKSFORME
	default:
		return gen.Resolution_RESOLUTION_RESOLUTION_UNSPECIFIED
	}
}

func IssueDomainTypeToProto(t domain.IssueType) gen.IssueType {
	switch t {
	case domain.IssueTypeCosmetic:
		return gen.IssueType_ISSUE_TYPE_COSMETIC
	case domain.IssueTypeBug:
		return gen.IssueType_ISSUE_TYPE_BUG
	case domain.IssueTypeFeature:
		return gen.IssueType_ISSUE_TYPE_FEATURE
	case domain.IssueTypePerformance:
		return gen.IssueType_ISSUE_TYPE_PERFORMANCE
	default:
		return gen.IssueType_ISSUE_TYPE_UNSPECIFIED
	}
}

func IssueDomainPriorityToProto(p domain.Priority) gen.Priority {
	switch p {
	case domain.PriorityCritical:
		return gen.Priority_PRIORITY_CRITICAL
	case domain.PriorityMajor:
		return gen.Priority_PRIORITY_MAJOR
	case domain.PriorityImportant:
		return gen.Priority_PRIORITY_IMPORTANT
	case domain.PriorityMinor:
		return gen.Priority_PRIORITY_MINOR
	default:
		return gen.Priority_PRIORITY_UNSPECIFIED
	}
}
