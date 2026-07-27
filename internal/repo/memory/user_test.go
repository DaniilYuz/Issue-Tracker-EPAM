package memory

import (
	"context"
	"testing"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	user := &domain.User{
		ID:           "user-123",
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
	}

	err = store.CreateUser(context.Background(), user)

	require.NoError(t, err)

	// Verify user was created
	retrieved, err := store.GetUserByID(context.Background(), "user-123")
	require.NoError(t, err)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.FirstName, retrieved.FirstName)
	assert.Equal(t, user.LastName, retrieved.LastName)
	assert.Equal(t, user.EmailAddress, retrieved.EmailAddress)
}

func TestGetUserByID(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	user := &domain.User{
		ID:           "user-456",
		FirstName:    "Jane",
		LastName:     "Smith",
		EmailAddress: "jane.smith@example.com",
	}

	// Create user first
	err = store.CreateUser(context.Background(), user)
	require.NoError(t, err)

	retrieved, err := store.GetUserByID(context.Background(), "user-456")

	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.FirstName, retrieved.FirstName)
	assert.Equal(t, user.LastName, retrieved.LastName)
	assert.Equal(t, user.EmailAddress, retrieved.EmailAddress)
}

func TestGetUserByIDNotFound(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	retrieved, err := store.GetUserByID(context.Background(), "non-existent-id")

	require.Error(t, err)
	assert.Nil(t, retrieved)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "user with this id did not found")
}

func TestUpdateUser(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	// Create original user
	originalUser := &domain.User{
		ID:           "user-789",
		FirstName:    "Bob",
		LastName:     "Johnson",
		EmailAddress: "bob.johnson@example.com",
	}
	err = store.CreateUser(context.Background(), originalUser)
	require.NoError(t, err)

	// Update user
	updatedUser := &domain.User{
		ID:           "user-789",
		FirstName:    "Robert",
		LastName:     "Johnson",
		EmailAddress: "robert.johnson@example.com",
	}

	err = store.UpdateUser(context.Background(), updatedUser)

	require.NoError(t, err)

	// Verify update
	retrieved, err := store.GetUserByID(context.Background(), "user-789")
	require.NoError(t, err)
	assert.Equal(t, "Robert", retrieved.FirstName)
	assert.Equal(t, "robert.johnson@example.com", retrieved.EmailAddress)
}

func TestUpdateUserNotFound(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	user := &domain.User{
		ID:           "non-existent-user",
		FirstName:    "Ghost",
		LastName:     "User",
		EmailAddress: "ghost@example.com",
	}

	err = store.UpdateUser(context.Background(), user)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "user not found")
}

func TestDeleteUser(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	// Create user first
	user := &domain.User{
		ID:           "user-delete",
		FirstName:    "Delete",
		LastName:     "Me",
		EmailAddress: "delete.me@example.com",
	}
	err = store.CreateUser(context.Background(), user)
	require.NoError(t, err)

	err = store.DeleteUser(context.Background(), "user-delete")

	require.NoError(t, err)

	// Verify user was deleted
	retrieved, err := store.GetUserByID(context.Background(), "user-delete")
	require.Error(t, err)
	assert.Nil(t, retrieved)
}

func TestListUsers(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	// Create multiple users
	users := []*domain.User{
		{ID: "user-1", FirstName: "User", LastName: "One", EmailAddress: "user1@example.com"},
		{ID: "user-2", FirstName: "User", LastName: "Two", EmailAddress: "user2@example.com"},
		{ID: "user-3", FirstName: "User", LastName: "Three", EmailAddress: "user3@example.com"},
	}

	for _, user := range users {
		err = store.CreateUser(context.Background(), user)
		require.NoError(t, err)
	}

	retrieved, err := store.ListUsers(context.Background())

	require.NoError(t, err)
	assert.Len(t, retrieved, 3)

	// Verify all users are present
	userIDs := make([]string, len(retrieved))
	for i, user := range retrieved {
		userIDs[i] = user.ID
	}
	assert.Contains(t, userIDs, "user-1")
	assert.Contains(t, userIDs, "user-2")
	assert.Contains(t, userIDs, "user-3")
}
