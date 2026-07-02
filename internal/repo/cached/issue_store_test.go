package cached

import (
	"context"
	"testing"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockIssueRepo struct {
	mock.Mock
}

func (m *MockIssueRepo) CreateIssue(ctx context.Context, issue *domain.Issue) error {
	return m.Called(ctx, issue).Error(0)
}

func (m *MockIssueRepo) GetIssueByID(ctx context.Context, id string) (*domain.Issue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Issue), args.Error(1)
}

func (m *MockIssueRepo) UpdateIssue(ctx context.Context, issue *domain.Issue) error {
	return m.Called(ctx, issue).Error(0)
}

func (m *MockIssueRepo) DeleteIssue(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockIssueRepo) ListIssues(ctx context.Context) ([]*domain.Issue, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Issue), args.Error(1)
}

type MockIssueCache struct {
	mock.Mock
}

func (m *MockIssueCache) SetIssue(ctx context.Context, ttl time.Duration, issue *domain.Issue) error {
	return m.Called(ctx, ttl, issue).Error(0)
}

func (m *MockIssueCache) SetIssues(ctx context.Context, ttl time.Duration, issues []*domain.Issue) error {
	return m.Called(ctx, ttl, issues).Error(0)
}

func (m *MockIssueCache) GetIssueByID(ctx context.Context, id string) (*domain.Issue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Issue), args.Error(1)
}

func (m *MockIssueCache) DeleteIssue(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockIssueCache) DeleteAllIssues(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func newIssueStoreForTest(next *MockIssueRepo, cache *MockIssueCache) *issueStore {
	return &issueStore{
		next:  next,
		cache: cache,
		ttl:   15 * time.Minute,
	}
}

func TestIssueStore_CreateIssue(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	issue := &domain.Issue{ID: "issue-1", Summary: "Fix bug"}

	mockNext.On("CreateIssue", mock.Anything, issue).Return(nil)
	mockCache.On("SetIssue", mock.Anything, 15*time.Minute, issue).Return(nil)

	store := newIssueStoreForTest(mockNext, mockCache)

	err := store.CreateIssue(context.Background(), issue)

	require.NoError(t, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestIssueStore_CreateIssue_NextError(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	issue := &domain.Issue{ID: "issue-1"}
	nextErr := status.Error(codes.Internal, "db error")

	mockNext.On("CreateIssue", mock.Anything, issue).Return(nextErr)

	store := newIssueStoreForTest(mockNext, mockCache)

	err := store.CreateIssue(context.Background(), issue)

	require.Error(t, err)
	assert.Equal(t, nextErr, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "SetIssue", mock.Anything, mock.Anything, mock.Anything)
}

func TestIssueStore_CreateIssue_CacheErrorIsSwallowed(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	issue := &domain.Issue{ID: "issue-1"}

	mockNext.On("CreateIssue", mock.Anything, issue).Return(nil)
	mockCache.On("SetIssue", mock.Anything, 15*time.Minute, issue).Return(status.Error(codes.Unavailable, "redis down"))

	store := newIssueStoreForTest(mockNext, mockCache)

	err := store.CreateIssue(context.Background(), issue)

	require.NoError(t, err, "cache errors must not fail the operation")
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestIssueStore_GetIssueByID_CacheHit(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	issue := &domain.Issue{ID: "issue-1", Summary: "Cached"}

	mockCache.On("GetIssueByID", mock.Anything, "issue-1").Return(issue, nil)

	store := newIssueStoreForTest(mockNext, mockCache)

	result, err := store.GetIssueByID(context.Background(), "issue-1")

	require.NoError(t, err)
	assert.Equal(t, issue, result)
	mockCache.AssertExpectations(t)
	mockNext.AssertNotCalled(t, "GetIssueByID", mock.Anything, mock.Anything)
}

func TestIssueStore_GetIssueByID_CacheMiss(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	issue := &domain.Issue{ID: "issue-1", Summary: "From DB"}

	mockCache.On("GetIssueByID", mock.Anything, "issue-1").Return(nil, status.Error(codes.NotFound, "cache miss"))
	mockNext.On("GetIssueByID", mock.Anything, "issue-1").Return(issue, nil)
	mockCache.On("SetIssue", mock.Anything, 15*time.Minute, issue).Return(nil)

	store := newIssueStoreForTest(mockNext, mockCache)

	result, err := store.GetIssueByID(context.Background(), "issue-1")

	require.NoError(t, err)
	assert.Equal(t, issue, result)
	mockCache.AssertExpectations(t)
	mockNext.AssertExpectations(t)
}

func TestIssueStore_GetIssueByID_NextError(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	mockCache.On("GetIssueByID", mock.Anything, "missing").Return(nil, status.Error(codes.NotFound, "cache miss"))
	mockNext.On("GetIssueByID", mock.Anything, "missing").Return(nil, status.Error(codes.NotFound, "issue \"missing\" not found"))

	store := newIssueStoreForTest(mockNext, mockCache)

	result, err := store.GetIssueByID(context.Background(), "missing")

	require.Error(t, err)
	assert.Nil(t, result)
	mockCache.AssertExpectations(t)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "SetIssue", mock.Anything, mock.Anything, mock.Anything)
}

func TestIssueStore_UpdateIssue(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	issue := &domain.Issue{ID: "issue-1", Summary: "Updated"}

	mockNext.On("UpdateIssue", mock.Anything, issue).Return(nil)
	mockCache.On("SetIssue", mock.Anything, 15*time.Minute, issue).Return(nil)

	store := newIssueStoreForTest(mockNext, mockCache)

	err := store.UpdateIssue(context.Background(), issue)

	require.NoError(t, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestIssueStore_UpdateIssue_NextError(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	issue := &domain.Issue{ID: "issue-1"}
	nextErr := status.Error(codes.NotFound, "issue not found")

	mockNext.On("UpdateIssue", mock.Anything, issue).Return(nextErr)

	store := newIssueStoreForTest(mockNext, mockCache)

	err := store.UpdateIssue(context.Background(), issue)

	require.Error(t, err)
	assert.Equal(t, nextErr, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "SetIssue", mock.Anything, mock.Anything, mock.Anything)
}

func TestIssueStore_DeleteIssue(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	mockNext.On("DeleteIssue", mock.Anything, "issue-1").Return(nil)
	mockCache.On("DeleteIssue", mock.Anything, "issue-1").Return(nil)

	store := newIssueStoreForTest(mockNext, mockCache)

	err := store.DeleteIssue(context.Background(), "issue-1")

	require.NoError(t, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestIssueStore_DeleteIssue_NextError(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	nextErr := status.Error(codes.NotFound, "issue not found")
	mockNext.On("DeleteIssue", mock.Anything, "issue-1").Return(nextErr)

	store := newIssueStoreForTest(mockNext, mockCache)

	err := store.DeleteIssue(context.Background(), "issue-1")

	require.Error(t, err)
	assert.Equal(t, nextErr, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "DeleteIssue", mock.Anything, mock.Anything)
}

func TestIssueStore_ListIssues_PassthroughToNext(t *testing.T) {
	mockNext := new(MockIssueRepo)
	mockCache := new(MockIssueCache)

	issues := []*domain.Issue{
		{ID: "issue-1"},
		{ID: "issue-2"},
	}

	mockNext.On("ListIssues", mock.Anything).Return(issues, nil)

	store := newIssueStoreForTest(mockNext, mockCache)

	result, err := store.ListIssues(context.Background())

	require.NoError(t, err)
	assert.Equal(t, issues, result)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "GetIssueByID", mock.Anything, mock.Anything)
	mockCache.AssertNotCalled(t, "SetIssues", mock.Anything, mock.Anything, mock.Anything)
}
