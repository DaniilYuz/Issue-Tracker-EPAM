package repository

import (
	"context"
	"encoding/json"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/cache"
	"git.epam.com/go-language-global-mentoring-program/internal/cache/utils"
	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"github.com/redis/go-redis/v9"
)

type issueRepository struct {
	redis redis.Cmdable
}

func NewIssueRepository(redis *redis.Client) cache.IssueCacheRepository {
	return &issueRepository{redis: redis}
}

func (r *issueRepository) SetIssue(ctx context.Context, ttl time.Duration, issue *domain.Issue) error {
	data, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	pipe := r.redis.TxPipeline()
	pipe.Set(ctx, utils.IssueKey(issue.ID), data, ttl)
	pipe.SAdd(ctx, "issues", issue.ID)

	_, err = pipe.Exec(ctx)

	return err
}

func (r *issueRepository) SetIssues(ctx context.Context, ttl time.Duration, issues []*domain.Issue) error {
	pipe := r.redis.TxPipeline()
	for _, issue := range issues {
		data, err := json.Marshal(issue)
		if err != nil {
			return err
		}

		pipe.Set(ctx, utils.IssueKey(issue.ID), data, ttl)
		pipe.SAdd(ctx, "issues", issue.ID)
	}

	_, err := pipe.Exec(ctx)

	return err
}

func (r *issueRepository) GetIssueByID(ctx context.Context, id string) (*domain.Issue, error) {
	var issue domain.Issue

	val, err := r.redis.Get(ctx, utils.IssueKey(id)).Result()
	if err == redis.Nil {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(val), &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

func (r *issueRepository) DeleteIssue(ctx context.Context, id string) error {
	pipe := r.redis.TxPipeline()
	pipe.Del(ctx, utils.IssueKey(id))
	pipe.SRem(ctx, "issues", id)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *issueRepository) DeleteAllIssues(ctx context.Context) error {
	ids, err := r.redis.SMembers(ctx, "issues").Result()
	if err != nil {
		return err
	}

	pipe := r.redis.TxPipeline()
	for _, id := range ids {
		pipe.Del(ctx, utils.IssueKey(id))
	}

	pipe.Del(ctx, "issues")
	_, err = pipe.Exec(ctx)

	return err
}
