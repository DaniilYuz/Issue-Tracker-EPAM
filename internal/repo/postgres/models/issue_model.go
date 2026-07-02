package models

import (
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
)

type IssueModel struct {
	ID          string                 `gorm:"primaryKey"`
	CreatedAt   time.Time              `gorm:"column:create_date;not null"`
	ModifiedAt  time.Time              `gorm:"column:modify_date;not null"`
	Summary     string                 `gorm:"not null"`
	Description string                 `gorm:"not null"`
	Status      domain.IssueStatus     `gorm:"type:issue_status;not null"`
	Resolution  domain.IssueResolution `gorm:"type:issue_resolution;not null;default:Unspecified"`
	Type        domain.IssueType       `gorm:"type:issue_type;not null"`
	Priority    domain.Priority        `gorm:"type:issue_priority;not null"`
	ProjectID   string                 `gorm:"not null"`
	AssigneeID  *string                `gorm:"default:null"`
}

func (IssueModel) TableName() string { return "issues" }

func IssueModelFromDomain(i *domain.Issue) *IssueModel {
	resolution := i.Resolution
	if resolution == "" {
		resolution = domain.ResolutionUnspecified
	}

	status := i.Status
	if status == "" {
		status = domain.IssueStatusNew
	}
	return &IssueModel{
		ID:          i.ID,
		CreatedAt:   i.CreateDate,
		ModifiedAt:  i.ModifyDate,
		Summary:     i.Summary,
		Description: i.Description,
		Status:      status,
		Resolution:  resolution,
		Type:        i.Type,
		Priority:    i.Priority,
		ProjectID:   i.ProjectID,
		AssigneeID:  toNullableString(i.AssigneeID),
	}
}

func IssueModelToDomain(i *IssueModel) *domain.Issue {
	assigneeID := ""
	if i.AssigneeID != nil {
		assigneeID = *i.AssigneeID
	}
	return &domain.Issue{
		ID:          i.ID,
		CreateDate:  i.CreatedAt,
		ModifyDate:  i.ModifiedAt,
		Summary:     i.Summary,
		Description: i.Description,
		Status:      i.Status,
		Resolution:  i.Resolution,
		Type:        i.Type,
		Priority:    i.Priority,
		ProjectID:   i.ProjectID,
		AssigneeID:  assigneeID,
	}
}

func toNullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
