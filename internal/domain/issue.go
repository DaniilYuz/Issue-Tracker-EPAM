package domain

import "time"

type Issue struct {
	ID string

	CreateDate time.Time
	ModifyDate time.Time

	Summary     string
	Description string

	Status     IssueStatus
	Resolution IssueResolution
	Type       IssueType
	Priority   Priority

	ProjectID  string
	AssigneeID string
}
