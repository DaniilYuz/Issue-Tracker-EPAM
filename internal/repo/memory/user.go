package memory

import (
	"context"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Store) CreateUser(ctx context.Context, user *domain.User) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("users", user); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	txn := s.db.Txn(true)
	defer txn.Abort()

	raw, err := txn.First("users", "id", id)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Error(codes.NotFound, "user with this id did not found")
	}

	return raw.(*domain.User), nil
}

func (s *Store) UpdateUser(ctx context.Context, newUser *domain.User) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	existing, err := txn.First("users", "id", newUser.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return status.Error(codes.NotFound, "user not found")
	}

	if err := txn.Insert("users", newUser); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	if err := txn.Delete("users", &domain.User{ID: id}); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *Store) ListUsers(ctx context.Context) ([]*domain.User, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("users", "id")
	if err != nil {
		return nil, err
	}

	var users []*domain.User
	for obj := it.Next(); obj != nil; obj = it.Next() {
		users = append(users, obj.(*domain.User))
	}

	return users, nil
}
