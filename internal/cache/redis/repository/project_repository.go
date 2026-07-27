package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/cache"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/cache/utils"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/redis/go-redis/v9"
)

type projectRepository struct {
	redis redis.Cmdable
}

func NewProjectRepository(redis *redis.Client) cache.ProjectCacheRepository {
	return &projectRepository{redis: redis}
}

func (r *projectRepository) SetProject(ctx context.Context, ttl time.Duration, project *domain.Project) error {
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	pipe := r.redis.TxPipeline()
	pipe.Set(ctx, utils.ProjectKey(project.ID), data, ttl)
	pipe.SAdd(ctx, "projects", project.ID)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *projectRepository) SetProjects(ctx context.Context, ttl time.Duration, projecrs []*domain.Project) error {
	pipe := r.redis.TxPipeline()
	for _, project := range projecrs {
		data, err := json.Marshal(project)
		if err != nil {
			return err
		}

		pipe.Set(ctx, utils.ProjectKey(project.ID), data, ttl)
		pipe.SAdd(ctx, "projects", project.ID)
	}

	_, err := pipe.Exec(ctx)

	return err
}

func (r *projectRepository) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	var project domain.Project

	val, err := r.redis.Get(ctx, utils.ProjectKey(id)).Result()
	if err == redis.Nil {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(val), &project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) DeleteProject(ctx context.Context, id string) error {
	pipe := r.redis.TxPipeline()
	pipe.Del(ctx, utils.ProjectKey(id))
	pipe.SRem(ctx, "projects", id)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *projectRepository) DeleteAllProjects(ctx context.Context) error {
	ids, err := r.redis.SMembers(ctx, "projects").Result()
	if err != nil {
		return err
	}

	pipe := r.redis.TxPipeline()
	for _, id := range ids {
		pipe.Del(ctx, utils.ProjectKey(id))
	}
	pipe.Del(ctx, "projects")
	_, err = pipe.Exec(ctx)
	return err
}
