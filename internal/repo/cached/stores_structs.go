package cached

import (
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/cache"
	"git.epam.com/go-language-global-mentoring-program/internal/repo"
)

type userStore struct {
	next  repo.UserRepository
	cache cache.UserCacheRepository
	ttl   time.Duration
}

type issueStore struct {
	next  repo.IssueRepository
	cache cache.IssueCacheRepository
	ttl   time.Duration
}

type projectStore struct {
	next  repo.ProjectRepository
	cache cache.ProjectCacheRepository
	ttl   time.Duration
}
