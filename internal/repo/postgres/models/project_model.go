package models

import "git.epam.com/go-language-global-mentoring-program/internal/domain"

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
