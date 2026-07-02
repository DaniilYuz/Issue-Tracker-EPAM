package repository

import (
	"context"
	"testing"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedisForUser(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return client, mr
}

func TestNewUserRepository(t *testing.T) {
	client, _ := setupTestRedisForUser(t)
	repo := NewUserRepository(client)
	assert.NotNil(t, repo)
	assert.IsType(t, &userRepository{}, repo)
}

func TestSetUser(t *testing.T) {
	client, _ := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	user := &domain.User{
		ID:           "user-1",
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
	}

	err := repo.SetUser(context.Background(), 15*time.Minute, user)
	assert.NoError(t, err)

	result, err := repo.GetUserByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.FirstName, result.FirstName)
	assert.Equal(t, user.EmailAddress, result.EmailAddress)
}

func TestSetUserError(t *testing.T) {
	client, mr := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	user := &domain.User{ID: "user-1", FirstName: "John", LastName: "Doe"}

	mr.Close()

	err := repo.SetUser(context.Background(), time.Minute, user)
	assert.Error(t, err)
}

func TestSetUsers(t *testing.T) {
	client, _ := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	users := []*domain.User{
		{ID: "user-1", FirstName: "John", LastName: "Doe", EmailAddress: "john@example.com"},
		{ID: "user-2", FirstName: "Jane", LastName: "Smith", EmailAddress: "jane@example.com"},
	}

	err := repo.SetUsers(context.Background(), 15*time.Minute, users)
	assert.NoError(t, err)

	for _, user := range users {
		result, err := repo.GetUserByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.FirstName, result.FirstName)
	}
}

func TestSetUsersError(t *testing.T) {
	client, mr := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	users := []*domain.User{
		{ID: "user-1", FirstName: "John", LastName: "Doe"},
	}

	mr.Close()

	err := repo.SetUsers(context.Background(), time.Minute, users)
	assert.Error(t, err)
}

func TestGetUserByID(t *testing.T) {
	client, _ := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	user := &domain.User{
		ID:           "user-1",
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
	}

	err := repo.SetUser(context.Background(), time.Hour, user)
	require.NoError(t, err)

	result, err := repo.GetUserByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.FirstName, result.FirstName)
	assert.Equal(t, user.EmailAddress, result.EmailAddress)
}

func TestGetUserByIDNotFound(t *testing.T) {
	client, _ := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	result, err := repo.GetUserByID(context.Background(), "missing-id")
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, result)
}

func TestGetUserByIDError(t *testing.T) {
	client, mr := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	mr.Close()

	result, err := repo.GetUserByID(context.Background(), "user-1")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetUserByIDUnmarshalError(t *testing.T) {
	client, mr := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	err := mr.Set("users:user-1", "invalid json data")
	require.NoError(t, err)

	result, err := repo.GetUserByID(context.Background(), "user-1")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteUser(t *testing.T) {
	client, _ := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	user := &domain.User{ID: "user-1", FirstName: "John", LastName: "Doe"}

	err := repo.SetUser(context.Background(), time.Hour, user)
	require.NoError(t, err)

	result, err := repo.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	err = repo.DeleteUser(context.Background(), user.ID)
	assert.NoError(t, err)

	result, err = repo.GetUserByID(context.Background(), user.ID)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, result)
}

func TestDeleteUserError(t *testing.T) {
	client, mr := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	mr.Close()

	err := repo.DeleteUser(context.Background(), "user-1")
	assert.Error(t, err)
}

func TestDeleteAllUsers(t *testing.T) {
	client, mr := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	users := []*domain.User{
		{ID: "user-1", FirstName: "John", LastName: "Doe"},
		{ID: "user-2", FirstName: "Jane", LastName: "Smith"},
		{ID: "user-3", FirstName: "Bob", LastName: "Johnson"},
	}

	for _, user := range users {
		err := repo.SetUser(context.Background(), time.Hour, user)
		require.NoError(t, err)
	}

	assert.True(t, mr.Exists("users"))
	members, _ := mr.SMembers("users")
	assert.Len(t, members, 3)

	err := repo.DeleteAllUsers(context.Background())
	assert.NoError(t, err)

	for _, user := range users {
		result, err := repo.GetUserByID(context.Background(), user.ID)
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
		assert.Nil(t, result)
	}

	assert.False(t, mr.Exists("users"))
}

func TestDeleteAllUsersError(t *testing.T) {
	client, mr := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	mr.Close()

	err := repo.DeleteAllUsers(context.Background())
	assert.Error(t, err)
}

func TestDeleteAllUsersEmpty(t *testing.T) {
	client, _ := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	err := repo.DeleteAllUsers(context.Background())
	assert.NoError(t, err)
}

func TestUserExpiration(t *testing.T) {
	client, mr := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	user := &domain.User{ID: "user-1", FirstName: "John", LastName: "Doe"}

	err := repo.SetUser(context.Background(), 100*time.Millisecond, user)
	require.NoError(t, err)

	result, err := repo.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	mr.FastForward(200 * time.Millisecond)

	result, err = repo.GetUserByID(context.Background(), user.ID)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, result)
}

func TestUserRedisState(t *testing.T) {
	client, mr := setupTestRedisForUser(t)
	repo := NewUserRepository(client)

	user := &domain.User{ID: "user-1", FirstName: "John", LastName: "Doe"}

	err := repo.SetUser(context.Background(), time.Hour, user)
	require.NoError(t, err)

	assert.True(t, mr.Exists("users:user-1"))
	assert.True(t, mr.Exists("users"))

	members, _ := mr.SMembers("users")
	assert.Contains(t, members, "user-1")
}
