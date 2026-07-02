package memory

import (
	"context"
	"testing"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateProject(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	project := &domain.Project{
		ID:          "project-123",
		Name:        "Test Project",
		Description: "A test project",
	}

	err = store.CreateProject(context.Background(), project)

	require.NoError(t, err)

	// Verify project was created
	retrieved, err := store.GetProjectByID(context.Background(), "project-123")
	require.NoError(t, err)
	assert.Equal(t, project.ID, retrieved.ID)
	assert.Equal(t, project.Name, retrieved.Name)
	assert.Equal(t, project.Description, retrieved.Description)
}

func TestGetProjectByID(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	project := &domain.Project{
		ID:          "project-456",
		Name:        "Another Project",
		Description: "Another test project",
	}

	// Create project first
	err = store.CreateProject(context.Background(), project)
	require.NoError(t, err)

	retrieved, err := store.GetProjectByID(context.Background(), "project-456")

	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, project.ID, retrieved.ID)
	assert.Equal(t, project.Name, retrieved.Name)
	assert.Equal(t, project.Description, retrieved.Description)
}

func TestGetProjectByIDNotFound(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	retrieved, err := store.GetProjectByID(context.Background(), "non-existent-project")

	require.Error(t, err)
	assert.Nil(t, retrieved)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "project with this id did not found")
}

func TestUpdateProject(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	// Create original project
	originalProject := &domain.Project{
		ID:          "project-789",
		Name:        "Original Project",
		Description: "Original description",
	}
	err = store.CreateProject(context.Background(), originalProject)
	require.NoError(t, err)

	// Update project
	updatedProject := &domain.Project{
		ID:          "project-789",
		Name:        "Updated Project",
		Description: "Updated description",
	}

	err = store.UpdateProject(context.Background(), updatedProject)

	require.NoError(t, err)

	// Verify update
	retrieved, err := store.GetProjectByID(context.Background(), "project-789")
	require.NoError(t, err)
	assert.Equal(t, "Updated Project", retrieved.Name)
	assert.Equal(t, "Updated description", retrieved.Description)
}

func TestUpdateProjectNotFound(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	project := &domain.Project{
		ID:          "non-existent-project",
		Name:        "Ghost Project",
		Description: "This project doesn't exist",
	}

	err = store.UpdateProject(context.Background(), project)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "project not found")
}

func TestDeleteProject(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	// Create project first
	project := &domain.Project{
		ID:          "project-delete",
		Name:        "Delete Me",
		Description: "This project will be deleted",
	}
	err = store.CreateProject(context.Background(), project)
	require.NoError(t, err)

	err = store.DeleteProject(context.Background(), "project-delete")

	require.NoError(t, err)

	// Verify project was deleted
	retrieved, err := store.GetProjectByID(context.Background(), "project-delete")
	require.Error(t, err)
	assert.Nil(t, retrieved)
}

func TestListProjects(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	// Create multiple projects
	projects := []*domain.Project{
		{ID: "project-1", Name: "Project One", Description: "First project"},
		{ID: "project-2", Name: "Project Two", Description: "Second project"},
		{ID: "project-3", Name: "Project Three", Description: "Third project"},
	}

	for _, project := range projects {
		err = store.CreateProject(context.Background(), project)
		require.NoError(t, err)
	}

	retrieved, err := store.ListProjects(context.Background())

	require.NoError(t, err)
	assert.Len(t, retrieved, 3)

	// Verify all projects are present
	projectIDs := make([]string, len(retrieved))
	for i, project := range retrieved {
		projectIDs[i] = project.ID
	}
	assert.Contains(t, projectIDs, "project-1")
	assert.Contains(t, projectIDs, "project-2")
	assert.Contains(t, projectIDs, "project-3")
}
