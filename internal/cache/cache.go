package cache

import (
	"context"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
)

type UserCacheRepository interface {
	SetUser(ctx context.Context, ttl time.Duration, user *domain.User) error
	SetUsers(ctx context.Context, ttl time.Duration, users []*domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	DeleteAllUsers(ctx context.Context) error
}

type ProjectCacheRepository interface {
	SetProject(ctx context.Context, ttl time.Duration, project *domain.Project) error
	SetProjects(ctx context.Context, ttl time.Duration, projecrs []*domain.Project) error
	GetProjectByID(ctx context.Context, id string) (*domain.Project, error)
	DeleteProject(ctx context.Context, id string) error
	DeleteAllProjects(ctx context.Context) error
}

type IssueCacheRepository interface {
	SetIssue(ctx context.Context, ttl time.Duration, issue *domain.Issue) error
	SetIssues(ctx context.Context, ttl time.Duration, issues []*domain.Issue) error
	GetIssueByID(ctx context.Context, id string) (*domain.Issue, error)
	DeleteIssue(ctx context.Context, id string) error
	DeleteAllIssues(ctx context.Context) error
}

type Cache interface {
	UserCacheRepository
	ProjectCacheRepository
	IssueCacheRepository
}
