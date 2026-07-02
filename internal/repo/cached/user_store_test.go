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

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUserRepo) ListUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

type MockUserCache struct {
	mock.Mock
}

func (m *MockUserCache) SetUser(ctx context.Context, ttl time.Duration, user *domain.User) error {
	return m.Called(ctx, ttl, user).Error(0)
}

func (m *MockUserCache) SetUsers(ctx context.Context, ttl time.Duration, users []*domain.User) error {
	return m.Called(ctx, ttl, users).Error(0)
}

func (m *MockUserCache) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserCache) DeleteUser(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUserCache) DeleteAllUsers(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func newUserStoreForTest(next *MockUserRepo, cache *MockUserCache) *userStore {
	return &userStore{
		next:  next,
		cache: cache,
		ttl:   15 * time.Minute,
	}
}

func TestUserStore_CreateUser(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	user := &domain.User{ID: "user-1", FirstName: "Daniil"}

	mockNext.On("CreateUser", mock.Anything, user).Return(nil)
	mockCache.On("SetUser", mock.Anything, 15*time.Minute, user).Return(nil)

	store := newUserStoreForTest(mockNext, mockCache)

	err := store.CreateUser(context.Background(), user)

	require.NoError(t, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserStore_CreateUser_NextError(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	user := &domain.User{ID: "user-1"}
	nextErr := status.Error(codes.AlreadyExists, "duplicate email")

	mockNext.On("CreateUser", mock.Anything, user).Return(nextErr)

	store := newUserStoreForTest(mockNext, mockCache)

	err := store.CreateUser(context.Background(), user)

	require.Error(t, err)
	assert.Equal(t, nextErr, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "SetUser", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserStore_CreateUser_CacheErrorIsSwallowed(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	user := &domain.User{ID: "user-1"}

	mockNext.On("CreateUser", mock.Anything, user).Return(nil)
	mockCache.On("SetUser", mock.Anything, 15*time.Minute, user).Return(status.Error(codes.Unavailable, "redis down"))

	store := newUserStoreForTest(mockNext, mockCache)

	err := store.CreateUser(context.Background(), user)

	require.NoError(t, err, "cache errors must not fail the operation")
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserStore_GetUserByID_CacheHit(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	user := &domain.User{ID: "user-1", FirstName: "Cached"}

	mockCache.On("GetUserByID", mock.Anything, "user-1").Return(user, nil)

	store := newUserStoreForTest(mockNext, mockCache)

	result, err := store.GetUserByID(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, user, result)
	mockCache.AssertExpectations(t)
	mockNext.AssertNotCalled(t, "GetUserByID", mock.Anything, mock.Anything)
}

func TestUserStore_GetUserByID_CacheMiss(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	user := &domain.User{ID: "user-1", FirstName: "From DB"}

	mockCache.On("GetUserByID", mock.Anything, "user-1").Return(nil, status.Error(codes.NotFound, "cache miss"))
	mockNext.On("GetUserByID", mock.Anything, "user-1").Return(user, nil)
	mockCache.On("SetUser", mock.Anything, 15*time.Minute, user).Return(nil)

	store := newUserStoreForTest(mockNext, mockCache)

	result, err := store.GetUserByID(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, user, result)
	mockCache.AssertExpectations(t)
	mockNext.AssertExpectations(t)
}

func TestUserStore_GetUserByID_NextError(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	mockCache.On("GetUserByID", mock.Anything, "missing").Return(nil, status.Error(codes.NotFound, "cache miss"))
	mockNext.On("GetUserByID", mock.Anything, "missing").Return(nil, status.Error(codes.NotFound, "user \"missing\" not found"))

	store := newUserStoreForTest(mockNext, mockCache)

	result, err := store.GetUserByID(context.Background(), "missing")

	require.Error(t, err)
	assert.Nil(t, result)
	mockCache.AssertExpectations(t)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "SetUser", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserStore_UpdateUser(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	user := &domain.User{ID: "user-1", FirstName: "Updated"}

	mockNext.On("UpdateUser", mock.Anything, user).Return(nil)
	mockCache.On("SetUser", mock.Anything, 15*time.Minute, user).Return(nil)

	store := newUserStoreForTest(mockNext, mockCache)

	err := store.UpdateUser(context.Background(), user)

	require.NoError(t, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserStore_UpdateUser_NextError(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	user := &domain.User{ID: "user-1"}
	nextErr := status.Error(codes.NotFound, "user not found")

	mockNext.On("UpdateUser", mock.Anything, user).Return(nextErr)

	store := newUserStoreForTest(mockNext, mockCache)

	err := store.UpdateUser(context.Background(), user)

	require.Error(t, err)
	assert.Equal(t, nextErr, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "SetUser", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserStore_DeleteUser(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	mockNext.On("DeleteUser", mock.Anything, "user-1").Return(nil)
	mockCache.On("DeleteUser", mock.Anything, "user-1").Return(nil)

	store := newUserStoreForTest(mockNext, mockCache)

	err := store.DeleteUser(context.Background(), "user-1")

	require.NoError(t, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserStore_DeleteUser_NextError(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	nextErr := status.Error(codes.NotFound, "user not found")
	mockNext.On("DeleteUser", mock.Anything, "user-1").Return(nextErr)

	store := newUserStoreForTest(mockNext, mockCache)

	err := store.DeleteUser(context.Background(), "user-1")

	require.Error(t, err)
	assert.Equal(t, nextErr, err)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "DeleteUser", mock.Anything, mock.Anything)
}

func TestUserStore_ListUsers_PassthroughToNext(t *testing.T) {
	mockNext := new(MockUserRepo)
	mockCache := new(MockUserCache)

	users := []*domain.User{
		{ID: "user-1"},
		{ID: "user-2"},
	}

	mockNext.On("ListUsers", mock.Anything).Return(users, nil)

	store := newUserStoreForTest(mockNext, mockCache)

	result, err := store.ListUsers(context.Background())

	require.NoError(t, err)
	assert.Equal(t, users, result)
	mockNext.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "GetUserByID", mock.Anything, mock.Anything)
	mockCache.AssertNotCalled(t, "SetUsers", mock.Anything, mock.Anything, mock.Anything)
}
