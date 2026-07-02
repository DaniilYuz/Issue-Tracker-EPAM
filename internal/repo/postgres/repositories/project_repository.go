package repositories

import (
	"context"
	"errors"
	"fmt"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"git.epam.com/go-language-global-mentoring-program/internal/repo"
	pgmodels "git.epam.com/go-language-global-mentoring-program/internal/repo/postgres/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) repo.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProject(ctx context.Context, project *domain.Project) error {
	if project.ID == "" {
		project.ID = uuid.NewString()
	}

	m := pgmodels.ProjectModelFromDomain(project)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}

	return nil
}

func (r *projectRepository) GetProjectByID(ctx context.Context, projectID string) (*domain.Project, error) {
	var m pgmodels.ProjectModel

	if err := r.db.WithContext(ctx).First(&m, "id = ? ", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("project %q not found", projectID)
		}
		return nil, fmt.Errorf("projectRepository.GetByID: %w", err)
	}
	return pgmodels.ProjectModelToDomain(&m), nil
}

func (r *projectRepository) UpdateProject(ctx context.Context, project *domain.Project) error {
	m := pgmodels.ProjectModelFromDomain(project)

	res := r.db.WithContext(ctx).Model(m).Where("id = ?", m.ID).Updates(m)
	if res.Error != nil {
		return fmt.Errorf("projectRepository: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("project %q not found", project.ID)
	}

	return nil
}

func (r *projectRepository) DeleteProject(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&pgmodels.ProjectModel{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("projectRepository.Delete: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("project %q not found", id)
	}
	return nil
}

func (r *projectRepository) ListProjects(ctx context.Context) ([]*domain.Project, error) {
	var models []pgmodels.ProjectModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("projectRepository.List: %w", err)
	}
	projects := make([]*domain.Project, len(models))
	for i := range models {
		projects[i] = pgmodels.ProjectModelToDomain(&models[i])
	}
	return projects, nil
}
