package domain

type IssueStatus string

const (
	IssueStatusNew        IssueStatus = "NEW"
	IssueStatusAssigned   IssueStatus = "ASSIGNED"
	IssueStatusInProgress IssueStatus = "IN_PROGRESS"
	IssueStatusResolved   IssueStatus = "RESOLVED"
	IssueStatusClosed     IssueStatus = "CLOSED"
	IssueStatusReopened   IssueStatus = "REOPENED"
)

type IssueResolution string

const (
	ResolutionFixed       IssueResolution = "FIXED"
	ResolutionInvalid     IssueResolution = "INVALID"
	ResolutionWontFix     IssueResolution = "WONTFIX"
	ResolutionWorksForMe  IssueResolution = "WORKSFORME"
	ResolutionUnspecified IssueResolution = "Unspecified"
)

type IssueType string

const (
	IssueTypeCosmetic    IssueType = "COSMETIC"
	IssueTypeBug         IssueType = "BUG"
	IssueTypeFeature     IssueType = "FEATURE"
	IssueTypePerformance IssueType = "PERFORMANCE"
)

type Priority string

const (
	PriorityCritical  Priority = "CRITICAL"
	PriorityMajor     Priority = "MAJOR"
	PriorityImportant Priority = "IMPORTANT"
	PriorityMinor     Priority = "MINOR"
)
