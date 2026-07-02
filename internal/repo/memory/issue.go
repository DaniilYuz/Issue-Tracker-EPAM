package memory

import (
	"context"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Store) CreateIssue(ctx context.Context, issue *domain.Issue) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("issues", issue); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *Store) GetIssueByID(ctx context.Context, id string) (*domain.Issue, error) {
	txn := s.db.Txn(true)
	defer txn.Abort()

	raw, err := txn.First("issues", "id", id)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Error(codes.NotFound, "issue with this id did not found")
	}

	return raw.(*domain.Issue), nil
}

func (s *Store) UpdateIssue(ctx context.Context, newIssue *domain.Issue) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	existing, err := txn.First("issues", "id", newIssue.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return status.Error(codes.NotFound, "issue not found")
	}

	if err := txn.Insert("issues", newIssue); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *Store) DeleteIssue(ctx context.Context, id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	if err := txn.Delete("issues", &domain.Issue{ID: id}); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *Store) ListIssues(ctx context.Context) ([]*domain.Issue, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("issues", "id")
	if err != nil {
		return nil, err
	}

	var issues []*domain.Issue
	for obj := it.Next(); obj != nil; obj = it.Next() {
		issues = append(issues, obj.(*domain.Issue))
	}

	return issues, nil
}
