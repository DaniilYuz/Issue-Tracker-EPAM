package cached

import (
	"context"
	"testing"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockProjectRepo struct {
	mock.Mock
}

func (m *MockProjectRepo) CreateProject(ctx context.Context, project *domain.Project) error {
	return m.Called(ctx, project).Error(0)
}

func (m *MockProjectRepo) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepo) UpdateProject(ctx context.Context, project *domain.Project) error {
	return m.Called(ctx, project).Error(0)
}

func (m *MockProjectRepo) DeleteProject(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockProjectRepo) ListProjects(ctx context.Context) ([]*domain.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

type MockProjectCache struct {
	mock.Mock
}

func (m *MockProjectCache) SetProject(ctx context.Context, ttl time.Duration, project *domain.Project) error {
	return m.Called(ctx, ttl, project).Error(0)
}

func (m *MockProjectCache) SetProjects(ctx context.Context, ttl time.Duration, projects []*domain.Project) error {
	return m.Called(ctx, ttl, projects).Error(0)
}

func (m *MockProjectCache) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectCache) DeleteProject(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockProjectCache) DeleteAllProjects(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func newProjectStoreForTest(next *MockProjectRepo, cache *MockProjectCache) *projectStore {
	return &projectStore{
		next:  next,
		cache: cache,
		ttl:   15 * time.Minute,
	}
}

func TestProjectStore_CreateProject(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	project := &domain.Project{ID: "project-1", Name: "Issue Tracker"}

	mockNext.On("CreateProject", mock.Anything, project).Return(nil)
	mockCache.On("SetProject", mock.Anything, 15*time.Minute, project).Return(nil)

	store := newProjectStoreForTest(mockNext, mockCache)

	err := store.CreateProject(context.Background(), project)

	require.NoError(t, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestProjectStore_CreateProject_NextError(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	project := &domain.Project{ID: "project-1"}
	nextErr := status.Error(codes.Internal, "db error")

	mockNext.On("CreateProject", mock.Anything, project).Return(nextErr)

	store := newProjectStoreForTest(mockNext, mockCache)

	err := store.CreateProject(context.Background(), project)

	require.Error(t, err)
	assert.Equal(t, nextErr, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "SetProject", mock.Anything, mock.Anything, mock.Anything)
}

func TestProjectStore_CreateProject_CacheErrorIsSwallowed(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	project := &domain.Project{ID: "project-1"}

	mockNext.On("CreateProject", mock.Anything, project).Return(nil)
	mockCache.On("SetProject", mock.Anything, 15*time.Minute, project).Return(status.Error(codes.Unavailable, "redis down"))

	store := newProjectStoreForTest(mockNext, mockCache)

	err := store.CreateProject(context.Background(), project)

	require.NoError(t, err, "cache errors must not fail the operation")
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestProjectStore_GetProjectByID_CacheHit(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	project := &domain.Project{ID: "project-1", Name: "Cached"}

	mockCache.On("GetProjectByID", mock.Anything, "project-1").Return(project, nil)

	store := newProjectStoreForTest(mockNext, mockCache)

	result, err := store.GetProjectByID(context.Background(), "project-1")

	require.NoError(t, err)
	assert.Equal(t, project, result)
	mockCache.AssertExpectations(t)
	mockNext.AssertNotCalled(t, "GetProjectByID", mock.Anything, mock.Anything)
}

func TestProjectStore_GetProjectByID_CacheMiss(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	project := &domain.Project{ID: "project-1", Name: "From DB"}

	mockCache.On("GetProjectByID", mock.Anything, "project-1").Return(nil, status.Error(codes.NotFound, "cache miss"))
	mockNext.On("GetProjectByID", mock.Anything, "project-1").Return(project, nil)
	mockCache.On("SetProject", mock.Anything, 15*time.Minute, project).Return(nil)

	store := newProjectStoreForTest(mockNext, mockCache)

	result, err := store.GetProjectByID(context.Background(), "project-1")

	require.NoError(t, err)
	assert.Equal(t, project, result)
	mockCache.AssertExpectations(t)
	mockNext.AssertExpectations(t)
}

func TestProjectStore_GetProjectByID_NextError(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	mockCache.On("GetProjectByID", mock.Anything, "missing").Return(nil, status.Error(codes.NotFound, "cache miss"))
	mockNext.On("GetProjectByID", mock.Anything, "missing").Return(nil, status.Error(codes.NotFound, "project \"missing\" not found"))

	store := newProjectStoreForTest(mockNext, mockCache)

	result, err := store.GetProjectByID(context.Background(), "missing")

	require.Error(t, err)
	assert.Nil(t, result)
	mockCache.AssertExpectations(t)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "SetProject", mock.Anything, mock.Anything, mock.Anything)
}

func TestProjectStore_UpdateProject(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	project := &domain.Project{ID: "project-1", Name: "Updated"}

	mockNext.On("UpdateProject", mock.Anything, project).Return(nil)
	mockCache.On("SetProject", mock.Anything, 15*time.Minute, project).Return(nil)

	store := newProjectStoreForTest(mockNext, mockCache)

	err := store.UpdateProject(context.Background(), project)

	require.NoError(t, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestProjectStore_UpdateProject_NextError(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	project := &domain.Project{ID: "project-1"}
	nextErr := status.Error(codes.NotFound, "project not found")

	mockNext.On("UpdateProject", mock.Anything, project).Return(nextErr)

	store := newProjectStoreForTest(mockNext, mockCache)

	err := store.UpdateProject(context.Background(), project)

	require.Error(t, err)
	assert.Equal(t, nextErr, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "SetProject", mock.Anything, mock.Anything, mock.Anything)
}

func TestProjectStore_DeleteProject(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	mockNext.On("DeleteProject", mock.Anything, "project-1").Return(nil)
	mockCache.On("DeleteProject", mock.Anything, "project-1").Return(nil)

	store := newProjectStoreForTest(mockNext, mockCache)

	err := store.DeleteProject(context.Background(), "project-1")

	require.NoError(t, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestProjectStore_DeleteProject_NextError(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	nextErr := status.Error(codes.NotFound, "project not found")
	mockNext.On("DeleteProject", mock.Anything, "project-1").Return(nextErr)

	store := newProjectStoreForTest(mockNext, mockCache)

	err := store.DeleteProject(context.Background(), "project-1")

	require.Error(t, err)
	assert.Equal(t, nextErr, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "DeleteProject", mock.Anything, mock.Anything)
}

func TestProjectStore_ListProjects_PassthroughToNext(t *testing.T) {
	mockNext := new(MockProjectRepo)
	mockCache := new(MockProjectCache)

	projects := []*domain.Project{
		{ID: "project-1"},
		{ID: "project-2"},
	}

	mockNext.On("ListProjects", mock.Anything).Return(projects, nil)

	store := newProjectStoreForTest(mockNext, mockCache)

	result, err := store.ListProjects(context.Background())

	require.NoError(t, err)
	assert.Equal(t, projects, result)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "GetProjectByID", mock.Anything, mock.Anything)
	mockCache.AssertNotCalled(t, "SetProjects", mock.Anything, mock.Anything, mock.Anything)
}
