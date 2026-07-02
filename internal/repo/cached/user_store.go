package cached

import (
	"context"
	"log"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
)

func (s *userStore) CreateUser(ctx context.Context, user *domain.User) error {
	if err := s.next.CreateUser(ctx, user); err != nil {
		return err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.SetUser(gCtx, s.ttl, user); err != nil {
		log.Printf("cache: failed to set user %s: %v", user.ID, err)
	}

	return nil
}

func (s *userStore) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if user, err := s.cache.GetUserByID(ctx, id); err == nil {
		return user, nil
	}

	user, err := s.next.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.SetUser(gCtx, s.ttl, user); err != nil {
		log.Printf("cache: failed to set user %s: %v", user.ID, err)
	}

	return user, nil
}

func (s *userStore) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := s.next.UpdateUser(ctx, user); err != nil {
		return err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.SetUser(gCtx, s.ttl, user); err != nil {
		log.Printf("cache: failed to set user %s: %v", user.ID, err)
	}

	return nil
}
func (s *userStore) DeleteUser(ctx context.Context, id string) error {
	if err := s.next.DeleteUser(ctx, id); err != nil {
		return err
	}

	gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.cache.DeleteUser(gCtx, id); err != nil {
		log.Printf("cache: failed to delete user %s: %v", id, err)
	}

	return nil
}
func (s *userStore) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return s.next.ListUsers(ctx)
}
