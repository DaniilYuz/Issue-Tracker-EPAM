package cached

import (
	"context"
	"log"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
)

func (s *projectStore) CreateProject(ctx context.Context, project *domain.Project) error {
	if err := s.next.CreateProject(ctx, project); err != nil {
		return err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.SetProject(gCtx, s.ttl, project); err != nil {
		log.Printf("cache: failed to set project %s: %v", project.ID, err)
	}

	return nil
}

func (s *projectStore) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	if project, err := s.cache.GetProjectByID(ctx, id); err == nil {
		return project, nil
	}

	project, err := s.next.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.SetProject(gCtx, s.ttl, project); err != nil {
		log.Printf("cache: failed to set project %s: %v", project.ID, err)
	}

	return project, nil
}

func (s *projectStore) UpdateProject(ctx context.Context, project *domain.Project) error {
	if err := s.next.UpdateProject(ctx, project); err != nil {
		return err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.SetProject(gCtx, s.ttl, project); err != nil {
		log.Printf("cache: failed to set project %s: %v", project.ID, err)
	}

	return nil
}
func (s *projectStore) DeleteProject(ctx context.Context, id string) error {
	if err := s.next.DeleteProject(ctx, id); err != nil {
		return err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.DeleteProject(gCtx, id); err != nil {
		log.Printf("cache: failed to delete project %s: %v", id, err)
	}

	return nil
}

func (s *projectStore) ListProjects(ctx context.Context) ([]*domain.Project, error) {
	return s.next.ListProjects(ctx)
}
