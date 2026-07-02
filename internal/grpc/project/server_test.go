package project

import (
	"context"
	"testing"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockProjectValidator struct {
	mock.Mock
}

func (m *MockProjectValidator) ValidateCreate(req *gen.CreateProjectRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockProjectValidator) ValidateUpdate(req *gen.UpdateProjectRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockProjectValidator) ValidateGet(req *gen.ReadProjectRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockProjectValidator) ValidateDelete(req *gen.DeleteProjectRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateProject(ctx context.Context, project *domain.Project) error {
	return m.Called(ctx, project).Error(0)
}

func (m *MockStore) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockStore) UpdateProject(ctx context.Context, project *domain.Project) error {
	return m.Called(ctx, project).Error(0)
}

func (m *MockStore) DeleteProject(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockStore) ListProjects(ctx context.Context) ([]*domain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func TestCreateProject(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	req := &gen.CreateProjectRequest{
		Name:        "New Platform",
		Description: "Internal tooling",
	}

	mockValidator.On("ValidateCreate", req).Return(nil)
	mockStore.On("CreateProject", mock.Anything, mock.AnythingOfType("*domain.Project")).Return(nil)

	server := NewServer(mockValidator, mockStore)

	res, err := server.CreateProject(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Project)
	assert.NotEmpty(t, res.Project.ProjectId, "Server must generate ULID for ProjectId")
	assert.Equal(t, req.Name, res.Project.Name)
	assert.Equal(t, req.Description, res.Project.Description)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestCreateProjectValidationFails(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	req := &gen.CreateProjectRequest{}
	validationErr := status.Error(codes.InvalidArgument, "project name is required")

	mockValidator.On("ValidateCreate", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	res, err := server.CreateProject(context.Background(), req)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
}

func TestReadProject(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	projectID := "project-ulid-123"
	req := &gen.ReadProjectRequest{ProjectId: projectID}

	domainProject := &domain.Project{
		ID:          projectID,
		Name:        "Test Project",
		Description: "Some desc",
	}

	mockValidator.On("ValidateGet", req).Return(nil)
	mockStore.On("GetProjectByID", mock.Anything, projectID).Return(domainProject, nil)

	server := NewServer(mockValidator, mockStore)

	res, err := server.ReadProject(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Project)
	assert.Equal(t, projectID, res.Project.ProjectId)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestReadProjectNotFound(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	projectID := "missing-id"
	req := &gen.ReadProjectRequest{ProjectId: projectID}

	mockValidator.On("ValidateGet", req).Return(nil)
	mockStore.On("GetProjectByID", mock.Anything, projectID).Return(nil, status.Error(codes.NotFound, "not found"))

	server := NewServer(mockValidator, mockStore)

	res, err := server.ReadProject(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestReadProjectValidationFails(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	projectID := "bad-id"
	req := &gen.ReadProjectRequest{ProjectId: projectID}
	validationErr := status.Error(codes.InvalidArgument, "invalid id format")

	mockValidator.On("ValidateGet", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	res, err := server.ReadProject(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
}

func TestUpdateProject(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	targetID := "update-project-id"
	reqProject := &gen.Project{
		ProjectId:   targetID,
		Name:        "Updated Name",
		Description: "Updated Desc",
	}

	req := &gen.UpdateProjectRequest{Project: reqProject}

	existingProject := &domain.Project{
		ID:          targetID,
		Name:        "Old Name",
		Description: "Old Desc",
	}

	mockValidator.On("ValidateUpdate", req).Return(nil)
	mockStore.On("GetProjectByID", mock.Anything, targetID).Return(existingProject, nil)
	mockStore.On("UpdateProject", mock.Anything, mock.AnythingOfType("*domain.Project")).Return(nil)

	server := NewServer(mockValidator, mockStore)

	res, err := server.UpdateProject(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, reqProject.Name, res.Project.Name)
	assert.Equal(t, reqProject.Description, res.Project.Description)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestUpdateProjectNotFound(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	reqProject := &gen.Project{ProjectId: "ghost-id", Name: "Ghost"}
	req := &gen.UpdateProjectRequest{Project: reqProject}

	mockValidator.On("ValidateUpdate", req).Return(nil)
	mockStore.On("GetProjectByID", mock.Anything, reqProject.ProjectId).Return(nil, status.Error(codes.NotFound, "not found"))

	server := NewServer(mockValidator, mockStore)

	res, err := server.UpdateProject(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestUpdateProjectValidationFails(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	targetID := "update-project-id"
	reqProject := &gen.Project{ProjectId: targetID}
	req := &gen.UpdateProjectRequest{Project: reqProject}
	validationErr := status.Error(codes.InvalidArgument, "empty project name")

	mockValidator.On("ValidateUpdate", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	res, err := server.UpdateProject(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
}

func TestDeleteProject(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	targetID := "delete-project-id"
	req := &gen.DeleteProjectRequest{ProjectId: targetID}

	mockValidator.On("ValidateDelete", req).Return(nil)
	mockStore.On("DeleteProject", mock.Anything, targetID).Return(nil)

	server := NewServer(mockValidator, mockStore)

	_, err := server.DeleteProject(context.Background(), req)

	require.NoError(t, err)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestDeleteProjectNotFound(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	projectID := "non-existent-id"
	req := &gen.DeleteProjectRequest{ProjectId: projectID}

	mockValidator.On("ValidateDelete", req).Return(nil)
	mockStore.On("DeleteProject", mock.Anything, projectID).Return(status.Error(codes.NotFound, "not found"))

	server := NewServer(mockValidator, mockStore)

	_, err := server.DeleteProject(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestDeleteProjectValidationFails(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	targetID := "delete-project-id"
	req := &gen.DeleteProjectRequest{ProjectId: targetID}
	validationErr := status.Error(codes.InvalidArgument, "invalid format")

	mockValidator.On("ValidateDelete", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	_, err := server.DeleteProject(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	mockValidator.AssertExpectations(t)
}

func TestListProjects(t *testing.T) {
	mockValidator := new(MockProjectValidator)
	mockStore := new(MockStore)

	domainProjects := []*domain.Project{
		{ID: "id-1", Name: "Project 1", Description: "Desc 1"},
		{ID: "id-2", Name: "Project 2", Description: "Desc 2"},
	}
	mockStore.On("ListProjects", mock.Anything).Return(domainProjects, nil)

	server := NewServer(mockValidator, mockStore)

	req := &gen.ListProjectsRequest{}
	res, err := server.ListProjects(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res.Projects, 2)

	projectIDs := []string{res.Projects[0].ProjectId, res.Projects[1].ProjectId}
	assert.Contains(t, projectIDs, "id-1")
	assert.Contains(t, projectIDs, "id-2")

	mockStore.AssertExpectations(t)
}
