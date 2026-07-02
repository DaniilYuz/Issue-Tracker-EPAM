package memory

import (
	"context"
	"encoding/json"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/cache/utils"
	"git.epam.com/go-language-global-mentoring-program/internal/domain"
)

type ProjectCache struct {
	store *Store
}

func NewProjectCache(store *Store) *ProjectCache {
	return &ProjectCache{store: store}
}

func (c *ProjectCache) SetProject(ctx context.Context, ttl time.Duration, project *domain.Project) error {
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}
	c.store.set(utils.ProjectKey(project.ID), data, ttl)
	return nil
}

func (c *ProjectCache) SetProjects(ctx context.Context, ttl time.Duration, projects []*domain.Project) error {
	for _, project := range projects {
		if err := c.SetProject(ctx, ttl, project); err != nil {
			return err
		}
	}
	return nil
}

func (c *ProjectCache) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	data, err := c.store.get(utils.ProjectKey(id))
	if err != nil {
		return nil, err
	}
	var project domain.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *ProjectCache) DeleteProject(ctx context.Context, id string) error {
	c.store.delete(utils.ProjectKey(id))
	return nil
}

func (c *ProjectCache) DeleteAllProjects(ctx context.Context) error {
	c.store.deleteByPrefix("projects:")
	return nil
}
