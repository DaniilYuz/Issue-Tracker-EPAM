package cached

import (
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/cache"
	"git.epam.com/go-language-global-mentoring-program/internal/repo"
)

type Store struct {
	repo.UserRepository
	repo.IssueRepository
	repo.ProjectRepository
}

func NewStore(next repo.Store, c cache.Cache, ttl time.Duration) repo.Store {
	return &Store{
		UserRepository:    newUserStore(next, c, ttl),
		ProjectRepository: newProjectStore(next, c, ttl),
		IssueRepository:   newIssueStore(next, c, ttl),
	}
}

func newUserStore(next repo.Store, c cache.Cache, ttl time.Duration) repo.UserRepository {
	return &userStore{
		next:  next,
		cache: c,
		ttl:   ttl,
	}
}

func newIssueStore(next repo.Store, c cache.Cache, ttl time.Duration) repo.IssueRepository {
	return &issueStore{
		next:  next,
		cache: c,
		ttl:   ttl,
	}
}

func newProjectStore(next repo.Store, c cache.Cache, ttl time.Duration) repo.ProjectRepository {
	return &projectStore{
		next:  next,
		cache: c,
		ttl:   ttl,
	}
}
