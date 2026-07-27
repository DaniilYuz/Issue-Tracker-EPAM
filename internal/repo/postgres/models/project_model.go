package models

import "github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"

type ProjectModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
}

func (ProjectModel) TableName() string { return "projects" }

func ProjectModelFromDomain(p *domain.Project) *ProjectModel {
	return &ProjectModel{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}
}

func ProjectModelToDomain(p *ProjectModel) *domain.Project {
	return &domain.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}
}
