package cached

import (
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/cache"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/repo"
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
