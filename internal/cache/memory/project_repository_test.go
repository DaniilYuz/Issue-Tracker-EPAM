package memory

import (
	"context"
	"testing"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCache_SetGetDeleteProject(t *testing.T) {
	store := NewStore()
	cache := NewProjectCache(store)
	ctx := context.Background()

	project := &domain.Project{
		ID:          "p1",
		Name:        "Project1",
		Description: "Description1",
	}
	err := cache.SetProject(ctx, time.Minute, project)
	require.NoError(t, err)

	got, err := cache.GetProjectByID(ctx, "p1")
	require.NoError(t, err)
	assert.Equal(t, project.ID, got.ID)
	assert.Equal(t, project.Name, got.Name)
	assert.Equal(t, project.Description, got.Description)

	err = cache.DeleteProject(ctx, "p1")
	require.NoError(t, err)

	got, err = cache.GetProjectByID(ctx, "p1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestProjectCache_SetProjects_DeleteAllProjects(t *testing.T) {
	store := NewStore()
	cache := NewProjectCache(store)
	ctx := context.Background()

	projects := []*domain.Project{
		{
			ID:          "p1",
			Name:        "Project1",
			Description: "Description1",
		},
		{
			ID:          "p2",
			Name:        "Project2",
			Description: "Description2",
		},
	}
	err := cache.SetProjects(ctx, time.Minute, projects)
	require.NoError(t, err)

	got1, err := cache.GetProjectByID(ctx, "p1")
	require.NoError(t, err)
	assert.Equal(t, "Project1", got1.Name)
	assert.Equal(t, "Description1", got1.Description)

	got2, err := cache.GetProjectByID(ctx, "p2")
	require.NoError(t, err)
	assert.Equal(t, "Project2", got2.Name)
	assert.Equal(t, "Description2", got2.Description)

	err = cache.DeleteAllProjects(ctx)
	require.NoError(t, err)

	got1, err = cache.GetProjectByID(ctx, "p1")
	assert.Error(t, err)
	assert.Nil(t, got1)

	got2, err = cache.GetProjectByID(ctx, "p2")
	assert.Error(t, err)
	assert.Nil(t, got2)
}

func TestProjectCache_TTL(t *testing.T) {
	store := NewStore()
	cache := NewProjectCache(store)
	ctx := context.Background()

	project := &domain.Project{
		ID:          "p1",
		Name:        "Project1",
		Description: "Description1",
	}
	err := cache.SetProject(ctx, 10*time.Millisecond, project)
	require.NoError(t, err)

	time.Sleep(20 * time.Millisecond)
	got, err := cache.GetProjectByID(ctx, "p1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestProjectCache_NotFound(t *testing.T) {
	store := NewStore()
	cache := NewProjectCache(store)
	ctx := context.Background()

	got, err := cache.GetProjectByID(ctx, "not-exist")
	assert.Error(t, err)
	assert.Nil(t, got)
}
