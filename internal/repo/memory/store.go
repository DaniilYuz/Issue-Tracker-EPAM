package memory

import (
	"git.epam.com/go-language-global-mentoring-program/internal/repo"
	"github.com/hashicorp/go-memdb"
)

type Store struct {
	db *memdb.MemDB
}

func NewStore() (repo.Store, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"users": {
				Name: "users",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"email": {
						Name:    "email",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "EmailAddress"},
					},
				},
			},
			"projects": {
				Name: "projects",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
				},
			},
			"issues": {
				Name: "issues",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"project_id": {
						Name:    "project_id",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "ProjectID"},
					},
					"assignee_id": {
						Name:    "assignee_id",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "AssigneeID"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}
