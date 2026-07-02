package postgres

import (
	"fmt"

	"git.epam.com/go-language-global-mentoring-program/internal/repo"
	"git.epam.com/go-language-global-mentoring-program/internal/repo/postgres/repositories"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
	repo.UserRepository
	repo.ProjectRepository
	repo.IssueRepository
}

func New(db *gorm.DB) (*Store, error) {
	if db == nil {
		return nil, fmt.Errorf("postgres.New: db is nil")
	}
	return &Store{
		db:                db,
		UserRepository:    repositories.NewUserRepository(db),
		ProjectRepository: repositories.NewProjectRepository(db),
		IssueRepository:   repositories.NewIssueRepository(db),
	}, nil
}

func (s *Store) Ping() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (s *Store) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
