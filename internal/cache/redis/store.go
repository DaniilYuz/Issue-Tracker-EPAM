package redis

import (
	"context"
	"fmt"

	"git.epam.com/go-language-global-mentoring-program/internal/cache"
	"git.epam.com/go-language-global-mentoring-program/internal/cache/redis/repository"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	r *redis.Client
	cache.UserCacheRepository
	cache.IssueCacheRepository
	cache.ProjectCacheRepository
}

func NewCache(r *redis.Client) (*Store, error) {
	if r == nil {
		return nil, fmt.Errorf("redis.New: db is nil")
	}

	return &Store{
		r:                      r,
		UserCacheRepository:    repository.NewUserRepository(r),
		IssueCacheRepository:   repository.NewIssueRepository(r),
		ProjectCacheRepository: repository.NewProjectRepository(r),
	}, nil
}

func (s *Store) Ping(ctx context.Context) error {
	return s.r.Ping(ctx).Err()
}

func (s *Store) Close() error {
	return s.r.Close()
}
