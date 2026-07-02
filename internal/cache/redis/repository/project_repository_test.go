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

func setupTestRedisForProject(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return client, mr
}

func TestNewProjectRepository(t *testing.T) {
	client, _ := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)
	assert.NotNil(t, repo)
	assert.IsType(t, &projectRepository{}, repo)
}

func TestSetProject(t *testing.T) {
	client, _ := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	project := &domain.Project{
		ID:          "project-1",
		Name:        "Test Project",
		Description: "Test Description",
	}

	err := repo.SetProject(context.Background(), 15*time.Minute, project)
	assert.NoError(t, err)

	result, err := repo.GetProjectByID(context.Background(), project.ID)
	assert.NoError(t, err)
	assert.Equal(t, project.ID, result.ID)
	assert.Equal(t, project.Name, result.Name)
	assert.Equal(t, project.Description, result.Description)
}

func TestSetProjectError(t *testing.T) {
	client, mr := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	project := &domain.Project{ID: "project-1", Name: "Test Project"}

	mr.Close()

	err := repo.SetProject(context.Background(), time.Minute, project)
	assert.Error(t, err)
}

func TestSetProjects(t *testing.T) {
	client, _ := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	projects := []*domain.Project{
		{ID: "project-1", Name: "Project 1", Description: "Desc 1"},
		{ID: "project-2", Name: "Project 2", Description: "Desc 2"},
	}

	err := repo.SetProjects(context.Background(), 15*time.Minute, projects)
	assert.NoError(t, err)

	for _, project := range projects {
		result, err := repo.GetProjectByID(context.Background(), project.ID)
		assert.NoError(t, err)
		assert.Equal(t, project.ID, result.ID)
		assert.Equal(t, project.Name, result.Name)
	}
}

func TestSetProjectsError(t *testing.T) {
	client, mr := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	projects := []*domain.Project{
		{ID: "project-1", Name: "Project 1"},
	}

	mr.Close()

	err := repo.SetProjects(context.Background(), time.Minute, projects)
	assert.Error(t, err)
}

func TestGetProjectByID(t *testing.T) {
	client, _ := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	project := &domain.Project{
		ID:          "project-1",
		Name:        "Test Project",
		Description: "Test Description",
	}

	err := repo.SetProject(context.Background(), time.Hour, project)
	require.NoError(t, err)

	result, err := repo.GetProjectByID(context.Background(), project.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, project.ID, result.ID)
	assert.Equal(t, project.Name, result.Name)
	assert.Equal(t, project.Description, result.Description)
}

func TestGetProjectByIDNotFound(t *testing.T) {
	client, _ := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	result, err := repo.GetProjectByID(context.Background(), "missing-id")
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, result)
}

func TestGetProjectByIDError(t *testing.T) {
	client, mr := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	mr.Close()

	result, err := repo.GetProjectByID(context.Background(), "project-1")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetProjectByIDUnmarshalError(t *testing.T) {
	client, mr := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	err := mr.Set("projects:project-1", "invalid json data")
	require.NoError(t, err)

	result, err := repo.GetProjectByID(context.Background(), "project-1")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteProject(t *testing.T) {
	client, _ := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	project := &domain.Project{ID: "project-1", Name: "Test Project"}

	err := repo.SetProject(context.Background(), time.Hour, project)
	require.NoError(t, err)

	result, err := repo.GetProjectByID(context.Background(), project.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	err = repo.DeleteProject(context.Background(), project.ID)
	assert.NoError(t, err)

	result, err = repo.GetProjectByID(context.Background(), project.ID)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, result)
}

func TestDeleteProjectError(t *testing.T) {
	client, mr := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	mr.Close()

	err := repo.DeleteProject(context.Background(), "project-1")
	assert.Error(t, err)
}

func TestDeleteAllProjects(t *testing.T) {
	client, mr := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	projects := []*domain.Project{
		{ID: "project-1", Name: "Project 1"},
		{ID: "project-2", Name: "Project 2"},
		{ID: "project-3", Name: "Project 3"},
	}

	for _, project := range projects {
		err := repo.SetProject(context.Background(), time.Hour, project)
		require.NoError(t, err)
	}

	assert.True(t, mr.Exists("projects"))
	members, _ := mr.SMembers("projects")
	assert.Len(t, members, 3)

	err := repo.DeleteAllProjects(context.Background())
	assert.NoError(t, err)

	for _, project := range projects {
		result, err := repo.GetProjectByID(context.Background(), project.ID)
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
		assert.Nil(t, result)
	}

	assert.False(t, mr.Exists("projects"))
}

func TestDeleteAllProjectsError(t *testing.T) {
	client, mr := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	mr.Close()

	err := repo.DeleteAllProjects(context.Background())
	assert.Error(t, err)
}

func TestDeleteAllProjectsEmpty(t *testing.T) {
	client, _ := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	err := repo.DeleteAllProjects(context.Background())
	assert.NoError(t, err)
}

func TestProjectExpiration(t *testing.T) {
	client, mr := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	project := &domain.Project{ID: "project-1", Name: "Test Project"}

	err := repo.SetProject(context.Background(), 100*time.Millisecond, project)
	require.NoError(t, err)

	result, err := repo.GetProjectByID(context.Background(), project.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	mr.FastForward(200 * time.Millisecond)

	result, err = repo.GetProjectByID(context.Background(), project.ID)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, result)
}

func TestProjectRedisState(t *testing.T) {
	client, mr := setupTestRedisForProject(t)
	repo := NewProjectRepository(client)

	project := &domain.Project{ID: "project-1", Name: "Test Project"}

	err := repo.SetProject(context.Background(), time.Hour, project)
	require.NoError(t, err)

	assert.True(t, mr.Exists("projects:project-1"))
	assert.True(t, mr.Exists("projects"))

	members, _ := mr.SMembers("projects")
	assert.Contains(t, members, "project-1")
}
