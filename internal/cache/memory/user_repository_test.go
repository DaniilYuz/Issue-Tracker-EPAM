package memory

import (
	"context"
	"testing"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCache_SetGetDeleteUser(t *testing.T) {
	store := NewStore()
	cache := NewUserCache(store)
	ctx := context.Background()

	user := &domain.User{
		ID:           "u1",
		FirstName:    "Alice",
		LastName:     "Smith",
		EmailAddress: "alice@example.com",
	}
	err := cache.SetUser(ctx, time.Minute, user)
	require.NoError(t, err)

	got, err := cache.GetUserByID(ctx, "u1")
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.FirstName, got.FirstName)
	assert.Equal(t, user.LastName, got.LastName)
	assert.Equal(t, user.EmailAddress, got.EmailAddress)

	err = cache.DeleteUser(ctx, "u1")
	require.NoError(t, err)

	got, err = cache.GetUserByID(ctx, "u1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestUserCache_SetUsers_DeleteAllUsers(t *testing.T) {
	store := NewStore()
	cache := NewUserCache(store)
	ctx := context.Background()

	users := []*domain.User{
		{
			ID:           "u1",
			FirstName:    "Alice",
			LastName:     "Smith",
			EmailAddress: "alice@example.com",
		},
		{
			ID:           "u2",
			FirstName:    "Bob",
			LastName:     "Johnson",
			EmailAddress: "bob@example.com",
		},
	}
	err := cache.SetUsers(ctx, time.Minute, users)
	require.NoError(t, err)

	got1, err := cache.GetUserByID(ctx, "u1")
	require.NoError(t, err)
	assert.Equal(t, "Alice", got1.FirstName)
	assert.Equal(t, "Smith", got1.LastName)
	assert.Equal(t, "alice@example.com", got1.EmailAddress)

	got2, err := cache.GetUserByID(ctx, "u2")
	require.NoError(t, err)
	assert.Equal(t, "Bob", got2.FirstName)
	assert.Equal(t, "Johnson", got2.LastName)
	assert.Equal(t, "bob@example.com", got2.EmailAddress)

	err = cache.DeleteAllUsers(ctx)
	require.NoError(t, err)

	got1, err = cache.GetUserByID(ctx, "u1")
	assert.Error(t, err)
	assert.Nil(t, got1)

	got2, err = cache.GetUserByID(ctx, "u2")
	assert.Error(t, err)
	assert.Nil(t, got2)
}

func TestUserCache_TTL(t *testing.T) {
	store := NewStore()
	cache := NewUserCache(store)
	ctx := context.Background()

	user := &domain.User{
		ID:           "u1",
		FirstName:    "Alice",
		LastName:     "Smith",
		EmailAddress: "alice@example.com",
	}
	err := cache.SetUser(ctx, 10*time.Millisecond, user)
	require.NoError(t, err)

	time.Sleep(20 * time.Millisecond)
	got, err := cache.GetUserByID(ctx, "u1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestUserCache_NotFound(t *testing.T) {
	store := NewStore()
	cache := NewUserCache(store)
	ctx := context.Background()

	got, err := cache.GetUserByID(ctx, "not-exist")
	assert.Error(t, err)
	assert.Nil(t, got)
}
