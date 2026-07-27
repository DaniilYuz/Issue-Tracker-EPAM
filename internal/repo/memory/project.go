package memory

import (
	"context"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Store) CreateProject(ctx context.Context, project *domain.Project) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("projects", project); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *Store) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	txn := s.db.Txn(true)
	defer txn.Abort()

	raw, err := txn.First("projects", "id", id)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Error(codes.NotFound, "project with this id did not found")
	}

	return raw.(*domain.Project), nil
}

func (s *Store) UpdateProject(ctx context.Context, newProject *domain.Project) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	existing, err := txn.First("projects", "id", newProject.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return status.Error(codes.NotFound, "project not found")
	}

	if err := txn.Insert("projects", newProject); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *Store) DeleteProject(ctx context.Context, id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	if err := txn.Delete("projects", &domain.Project{ID: id}); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *Store) ListProjects(ctx context.Context) ([]*domain.Project, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("projects", "id")
	if err != nil {
		return nil, err
	}

	var projects []*domain.Project
	for obj := it.Next(); obj != nil; obj = it.Next() {
		projects = append(projects, obj.(*domain.Project))
	}

	return projects, nil
}
