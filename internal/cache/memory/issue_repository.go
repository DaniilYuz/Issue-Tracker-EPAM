package memory

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/cache/utils"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
)

type IssueCache struct {
	store *Store
}

func NewIssueCache(store *Store) *IssueCache {
	return &IssueCache{store: store}
}

func (c *IssueCache) SetIssue(ctx context.Context, ttl time.Duration, issue *domain.Issue) error {
	data, err := json.Marshal(issue)
	if err != nil {
		return err
	}
	c.store.set(utils.IssueKey(issue.ID), data, ttl)
	return nil
}

func (c *IssueCache) SetIssues(ctx context.Context, ttl time.Duration, issues []*domain.Issue) error {
	for _, issue := range issues {
		if err := c.SetIssue(ctx, ttl, issue); err != nil {
			return err
		}
	}
	return nil
}

func (c *IssueCache) GetIssueByID(ctx context.Context, id string) (*domain.Issue, error) {
	data, err := c.store.get(utils.IssueKey(id))
	if err != nil {
		return nil, err
	}
	var issue domain.Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

func (c *IssueCache) DeleteIssue(ctx context.Context, id string) error {
	c.store.delete(utils.IssueKey(id))
	return nil
}

func (c *IssueCache) DeleteAllIssues(ctx context.Context) error {
	c.store.deleteByPrefix("issues:")
	return nil
}
