package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/repo"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/repo/postgres/models"
	pgmodels "github.com/DaniilYuz/Issue-Tracker-EPAM/internal/repo/postgres/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type issueRepository struct {
	db *gorm.DB
}

func NewIssueRepository(db *gorm.DB) repo.IssueRepository {
	return &issueRepository{db: db}
}

func (r *issueRepository) CreateIssue(ctx context.Context, issue *domain.Issue) error {
	if issue.ID == "" {
		issue.ID = uuid.NewString()
	}

	m := models.IssueModelFromDomain(issue)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}

	return nil
}

func (r *issueRepository) GetIssueByID(ctx context.Context, issueID string) (*domain.Issue, error) {
	var m pgmodels.IssueModel

	if err := r.db.WithContext(ctx).First(&m, "id = ?", issueID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("issue %q not found", issueID)
		}
		return nil, fmt.Errorf("issueRepository.GetIssueByID: %w", err)
	}
	return pgmodels.IssueModelToDomain(&m), nil
}

func (r *issueRepository) UpdateIssue(ctx context.Context, issue *domain.Issue) error {
	m := pgmodels.IssueModelFromDomain(issue)

	res := r.db.WithContext(ctx).Model(m).Where("id = ?", m.ID).Updates(m)
	if res.Error != nil {
		return fmt.Errorf("issueRepository: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("issue %q not found", issue.ID)
	}

	return nil
}

func (r *issueRepository) DeleteIssue(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&pgmodels.IssueModel{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("issueRepository.DeleteIssue: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("issue %q not found", id)
	}
	return nil
}

func (r *issueRepository) ListIssues(ctx context.Context) ([]*domain.Issue, error) {
	var models []pgmodels.IssueModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("issueRepository.ListIssues: %w", err)
	}
	issues := make([]*domain.Issue, len(models))
	for i := range models {
		issues[i] = pgmodels.IssueModelToDomain(&models[i])
	}
	return issues, nil
}
