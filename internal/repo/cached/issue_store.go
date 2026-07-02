package cached

import (
	"context"
	"log"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
)

func (s *issueStore) CreateIssue(ctx context.Context, issue *domain.Issue) error {
	if err := s.next.CreateIssue(ctx, issue); err != nil {
		return err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.SetIssue(gCtx, s.ttl, issue); err != nil {
		log.Printf("cache: failed to set issue %s: %v", issue.ID, err)
	}

	return nil
}

func (s *issueStore) GetIssueByID(ctx context.Context, id string) (*domain.Issue, error) {
	if issue, err := s.cache.GetIssueByID(ctx, id); err == nil {
		return issue, nil
	}

	issue, err := s.next.GetIssueByID(ctx, id)
	if err != nil {
		return nil, err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.SetIssue(gCtx, s.ttl, issue); err != nil {
		log.Printf("cache: failed to set issue %s: %v", issue.ID, err)
	}

	return issue, nil
}

func (s *issueStore) UpdateIssue(ctx context.Context, issue *domain.Issue) error {
	if err := s.next.UpdateIssue(ctx, issue); err != nil {
		return err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.SetIssue(gCtx, s.ttl, issue); err != nil {
		log.Printf("cache: failed to set issue %s: %v", issue.ID, err)
	}

	return nil
}

func (s *issueStore) DeleteIssue(ctx context.Context, id string) error {
	if err := s.next.DeleteIssue(ctx, id); err != nil {
		return err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.DeleteIssue(gCtx, id); err != nil {
		log.Printf("cache: failed to delete issue %s: %v", id, err)
	}

	return nil
}

func (s *issueStore) ListIssues(ctx context.Context) ([]*domain.Issue, error) {
	return s.next.ListIssues(ctx)
}
