package repo

import (
	"context"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]*domain.User, error)
}

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *domain.Project) error
	GetProjectByID(ctx context.Context, id string) (*domain.Project, error)
	UpdateProject(ctx context.Context, project *domain.Project) error
	DeleteProject(ctx context.Context, id string) error
	ListProjects(ctx context.Context) ([]*domain.Project, error)
}

type IssueRepository interface {
	CreateIssue(ctx context.Context, issue *domain.Issue) error
	GetIssueByID(ctx context.Context, id string) (*domain.Issue, error)
	UpdateIssue(ctx context.Context, issue *domain.Issue) error
	DeleteIssue(ctx context.Context, id string) error
	ListIssues(ctx context.Context) ([]*domain.Issue, error)
}

type Store interface {
	UserRepository
	ProjectRepository
	IssueRepository
}
