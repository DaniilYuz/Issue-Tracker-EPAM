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

func setupUserTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	apiDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialer := postgres.New(postgres.Config{
		Conn: apiDB,
	})

	gormDB, err := gorm.Open(dialer, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock := setupUserTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{
		ID:           "user-123",
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(user.ID, user.FirstName, user.LastName, user.EmailAddress).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateUser(context.Background(), user)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_Success(t *testing.T) {
	db, mock := setupUserTestDB(t)
	repo := NewUserRepository(db)

	userID := "user-123"
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email_address"}).
		AddRow(userID, "John", "Doe", "john.doe@example.com")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnRows(rows)

	retrieved, err := repo.GetUserByID(context.Background(), userID)

	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, userID, retrieved.ID)
	assert.Equal(t, "John", retrieved.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	db, mock := setupUserTestDB(t)
	repo := NewUserRepository(db)

	userID := "non-existent"
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	retrieved, err := repo.GetUserByID(context.Background(), userID)

	require.Error(t, err)
	assert.Nil(t, retrieved)
	assert.Contains(t, err.Error(), `user "non-existent" not found`)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateUser_Success(t *testing.T) {
	db, mock := setupUserTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{
		ID:           "user-123",
		FirstName:    "Jane",
		LastName:     "Doe",
		EmailAddress: "jane.doe@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET .+ WHERE id = \$\d+ AND "id" = \$\d+`).
		WithArgs(
			user.FirstName,
			user.LastName,
			user.EmailAddress,
			user.ID,
			user.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateUser(context.Background(), user)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateUser_NotFound(t *testing.T) {
	db, mock := setupUserTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{
		ID:           "ghost-id",
		FirstName:    "Ghost",
		LastName:     "User",
		EmailAddress: "ghost@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET .+ WHERE .+`).
		WithArgs(user.FirstName, user.LastName, user.EmailAddress, user.ID, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.UpdateUser(context.Background(), user)

	require.Error(t, err)
	assert.Contains(t, err.Error(), `user "ghost-id" not found`)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_DeleteUser_Success(t *testing.
	T) {
	db, mock := setupUserTestDB(t)
	repo := NewUserRepository(db)

	userID := "user-delete"

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "users"`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteUser(context.Background(), userID)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ListUsers(t *testing.T) {
	db, mock := setupUserTestDB(t)
	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email_address"}).
		AddRow("id-1", "User 1", "Last 1", "user1@example.com").
		AddRow("id-2", "User 2", "Last 2", "user2@example.com")

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(rows)

	users, err := repo.ListUsers(context.Background())

	require.NoError(t, err)
	require.Len(t, users, 2)
	assert.Equal(t, "id-1", users[0].ID)
	assert.Equal(t, "id-2", users[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
