package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupProjectTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	apiDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialer := postgres.New(postgres.Config{
		Conn: apiDB,
	})

	gormDB, err := gorm.Open(dialer, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestProjectRepository_CreateProject(t *testing.T) {
	db, mock := setupProjectTestDB(t)
	repo := NewProjectRepository(db)

	project := &domain.Project{
		ID:          "project-123",
		Name:        "Test Project",
		Description: "A test project description",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "projects"`).
		WithArgs(project.ID, project.Name, project.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateProject(context.Background(), project)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_GetProjectByID_Success(t *testing.T) {
	db, mock := setupProjectTestDB(t)
	repo := NewProjectRepository(db)

	projectID := "project-123"
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(projectID, "Test Project", "Description")

	mock.ExpectQuery(`SELECT \* FROM "projects" WHERE id = \$1 ORDER BY "projects"\."id" LIMIT \$2`).
		WithArgs(projectID, 1).
		WillReturnRows(rows)

	retrieved, err := repo.GetProjectByID(context.Background(), projectID)

	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, projectID, retrieved.ID)
	assert.Equal(t, "Test Project", retrieved.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_GetProjectByID_NotFound(t *testing.T) {
	db, mock := setupProjectTestDB(t)
	repo := NewProjectRepository(db)

	projectID := "missing-project"
	mock.ExpectQuery(`SELECT \* FROM "projects" WHERE id = \$1 ORDER BY "projects"\."id" LIMIT \$2`).
		WithArgs(projectID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	retrieved, err := repo.GetProjectByID(context.Background(), projectID)

	require.Error(t, err)
	assert.Nil(t, retrieved)
	assert.Contains(t, err.Error(), `project "missing-project" not found`)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_UpdateProject_Success(t *testing.T) {
	db, mock := setupProjectTestDB(t)
	repo := NewProjectRepository(db)

	project := &domain.Project{
		ID:          "project-123",
		Name:        "Updated Project",
		Description: "Updated description",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "projects" SET .+ WHERE id = \$\d+ AND "id" = \$\d+`).
		WithArgs(project.Name, project.Description, project.ID, project.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateProject(context.Background(), project)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_UpdateProject_NotFound(t *testing.T) {
	db, mock := setupProjectTestDB(t)
	repo := NewProjectRepository(db)

	project := &domain.Project{
		ID:          "ghost-project",
		Name:        "Ghost Project",
		Description: "Ghost Description",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "projects" SET .+ WHERE .+`).
		WithArgs(project.Name, project.Description, project.ID, project.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.UpdateProject(context.Background(), project)

	require.Error(t, err)
	assert.Contains(t, err.Error(), `project "ghost-project" not found`)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_DeleteProject_Success(t *testing.T) {
	db, mock := setupProjectTestDB(t)
	repo := NewProjectRepository(db)

	projectID := "project-delete"

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "projects"`).
		WithArgs(projectID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteProject(context.Background(), projectID)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_ListProjects(t *testing.T) {
	db, mock := setupProjectTestDB(t)
	repo := NewProjectRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow("id-1", "Project 1", "Desc 1").
		AddRow("id-2", "Project 2", "Desc 2")

	mock.ExpectQuery(`SELECT \* FROM "projects"`).
		WillReturnRows(rows)

	projects, err := repo.ListProjects(context.Background())

	require.NoError(t, err)
	require.Len(t, projects, 2)
	assert.Equal(t, "id-1", projects[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
