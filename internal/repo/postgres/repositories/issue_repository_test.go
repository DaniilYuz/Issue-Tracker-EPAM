package repositories

import (
	"context"
	"testing"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupIssueTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	apiDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialer := postgres.New(postgres.Config{
		Conn: apiDB,
	})

	gormDB, err := gorm.Open(dialer, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestIssueRepository_CreateIssue(t *testing.T) {
	db, mock := setupIssueTestDB(t)
	repo := NewIssueRepository(db)

	now := time.Now()
	assignee := "user-123"
	issue := &domain.Issue{
		ID:          "issue-123",
		CreateDate:  now,
		ModifyDate:  now,
		Summary:     "Test Issue",
		Description: "Description",
		Status:      domain.IssueStatusNew,
		Resolution:  domain.ResolutionUnspecified,
		Type:        domain.IssueTypeBug,
		Priority:    domain.PriorityMinor,
		ProjectID:   "project-123",
		AssigneeID:  assignee,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "issues"`).
		WithArgs(
			issue.ID,
			issue.CreateDate,
			issue.ModifyDate,
			issue.Summary,
			issue.Description,
			issue.Status,
			issue.Resolution,
			issue.Type,
			issue.Priority,
			issue.ProjectID,
			assignee,
		).
		WillReturnRows(sqlmock.NewRows([]string{"assignee_id"}).AddRow(assignee))
	mock.ExpectCommit()

	err := repo.CreateIssue(context.Background(), issue)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIssueRepository_GetIssueByID_Success(t *testing.T) {
	db, mock := setupIssueTestDB(t)
	repo := NewIssueRepository(db)

	issueID := "issue-123"
	now := time.Now()
	assignee := "user-123"

	rows := sqlmock.NewRows([]string{
		"id", "create_date", "modify_date", "summary", "description",
		"status", "resolution", "type", "priority", "project_id", "assignee_id",
	}).AddRow(
		issueID, now, now, "Test Issue", "Description",
		domain.IssueStatusNew, domain.ResolutionUnspecified, domain.IssueTypeBug, domain.PriorityMinor, "project-123", assignee,
	)

	mock.ExpectQuery(`SELECT \* FROM "issues" WHERE id = \$1 ORDER BY "issues"\."id" LIMIT \$2`).
		WithArgs(issueID, 1).
		WillReturnRows(rows)

	retrieved, err := repo.GetIssueByID(context.Background(), issueID)

	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, issueID, retrieved.ID)
	assert.Equal(t, "Test Issue", retrieved.Summary)
	assert.Equal(t, assignee, retrieved.AssigneeID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIssueRepository_GetIssueByID_NotFound(t *testing.T) {
	db, mock := setupIssueTestDB(t)
	repo := NewIssueRepository(db)

	issueID := "missing-issue"
	mock.ExpectQuery(`SELECT \* FROM "issues" WHERE id = \$1 ORDER BY "issues"\."id" LIMIT \$2`).
		WithArgs(issueID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	retrieved, err := repo.GetIssueByID(context.Background(), issueID)

	require.Error(t, err)
	assert.Nil(t, retrieved)
	assert.Contains(t, err.Error(), `issue "missing-issue" not found`)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIssueRepository_UpdateIssue_Success(t *testing.T) {
	db, mock := setupIssueTestDB(t)
	repo := NewIssueRepository(db)

	now := time.Now()
	issue := &domain.Issue{
		ID:          "issue-123",
		CreateDate:  now,
		ModifyDate:  now,
		Summary:     "Updated Summary",
		Description: "Updated Description",
		Status:      domain.IssueStatusInProgress,
		Resolution:  domain.ResolutionUnspecified,
		Type:        domain.IssueTypeFeature,
		Priority:    domain.PriorityMajor,
		ProjectID:   "project-123",
		AssigneeID:  "",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "issues" SET .+ WHERE id = \$\d+ AND "id" = \$\d+`).
		WithArgs(
			issue.CreateDate,
			issue.ModifyDate,
			issue.Summary,
			issue.Description,
			issue.Status,
			issue.Resolution,
			issue.Type,
			issue.Priority,
			issue.ProjectID,
			issue.ID,
			issue.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateIssue(context.Background(), issue)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIssueRepository_UpdateIssue_NotFound(t *testing.T) {
	db, mock := setupIssueTestDB(t)
	repo := NewIssueRepository(db)

	issue := &domain.Issue{
		ID: "ghost-issue",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "issues"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.UpdateIssue(context.Background(), issue)

	require.Error(t, err)
	assert.Contains(t, err.Error(), `issue "ghost-issue" not found`)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIssueRepository_DeleteIssue_Success(t *testing.T) {
	db, mock := setupIssueTestDB(t)
	repo := NewIssueRepository(db)

	issueID := "issue-delete"

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "issues"`).
		WithArgs(issueID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteIssue(context.Background(), issueID)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIssueRepository_ListIssues(t *testing.T) {
	db, mock := setupIssueTestDB(t)
	repo := NewIssueRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "create_date", "modify_date", "summary", "description",
		"status", "resolution", "type", "priority", "project_id", "assignee_id",
	}).AddRow(
		"id-1", now, now, "Issue 1", "Desc 1",
		domain.IssueStatusNew, domain.ResolutionUnspecified, domain.IssueTypeBug, domain.PriorityMinor, "proj-1", nil,
	)

	mock.ExpectQuery(`SELECT \* FROM "issues"`).
		WillReturnRows(rows)

	issues, err := repo.ListIssues(context.Background())

	require.NoError(t, err)
	require.Len(t, issues, 1)
	assert.Equal(t, "id-1", issues[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
