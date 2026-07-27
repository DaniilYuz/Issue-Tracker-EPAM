package user

import (
	"context"
	"testing"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockUserValidator struct {
	mock.Mock
}

func (m *MockUserValidator) ValidateCreate(req *gen.CreateUserRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockUserValidator) ValidateUpdate(req *gen.UpdateUserRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockUserValidator) ValidateGet(req *gen.ReadUserRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockUserValidator) ValidateDelete(req *gen.DeleteUserRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateUser(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockStore) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockStore) UpdateUser(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockStore) DeleteUser(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockStore) ListUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	req := &gen.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
	}

	mockValidator.On("ValidateCreate", req).Return(nil)
	mockStore.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	server := NewServer(mockValidator, mockStore)

	res, err := server.CreateUser(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.User)
	assert.NotEmpty(t, res.User.UserId, "Server must generate ULID for UserId")
	assert.Equal(t, req.FirstName, res.User.FirstName)
	assert.Equal(t, req.LastName, res.User.LastName)
	assert.Equal(t, req.EmailAddress, res.User.EmailAddress)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestCreateUserValidationFails(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	req := &gen.CreateUserRequest{}
	validationErr := status.Error(codes.InvalidArgument, "invalid email address")

	mockValidator.On("ValidateCreate", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	res, err := server.CreateUser(context.Background(), req)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
}

func TestReadUser(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	userID := "user-ulid-111"
	req := &gen.ReadUserRequest{UserId: userID}

	domainUser := &domain.User{
		ID:           userID,
		FirstName:    "Alice",
		LastName:     "Smith",
		EmailAddress: "alice@example.com",
	}

	mockValidator.On("ValidateGet", req).Return(nil)
	mockStore.On("GetUserByID", mock.Anything, userID).Return(domainUser, nil)

	server := NewServer(mockValidator, mockStore)

	res, err := server.ReadUser(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)

	require.NotNil(t, res.User)
	assert.Equal(t, userID, res.User.UserId)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestReadUserNotFound(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	userID := "missing-user-id"
	req := &gen.ReadUserRequest{UserId: userID}

	mockValidator.On("ValidateGet", req).Return(nil)
	mockStore.On("GetUserByID", mock.Anything, userID).Return(nil, status.Error(codes.NotFound, "not found"))

	server := NewServer(mockValidator, mockStore)

	res, err := server.ReadUser(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestReadUserValidationFails(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	userID := "bad-id"
	req := &gen.ReadUserRequest{UserId: userID}
	validationErr := status.Error(codes.InvalidArgument, "invalid user id format")

	mockValidator.On("ValidateGet", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	res, err := server.ReadUser(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	targetID := "update-user-id"
	reqUser := &gen.User{
		UserId:       targetID,
		FirstName:    "Bob",
		LastName:     "Updated",
		EmailAddress: "bob.new@example.com",
	}

	req := &gen.UpdateUserRequest{User: reqUser}

	existingUser := &domain.User{
		ID:           targetID,
		FirstName:    "Bob",
		LastName:     "Old",
		EmailAddress: "bob@example.com",
	}

	mockValidator.On("ValidateUpdate", req).Return(nil)
	mockStore.On("GetUserByID", mock.Anything, targetID).Return(existingUser, nil)
	mockStore.On("UpdateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	server := NewServer(mockValidator, mockStore)

	res, err := server.UpdateUser(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, reqUser.LastName, res.User.LastName)
	assert.Equal(t, reqUser.EmailAddress, res.User.EmailAddress)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestUpdateUserNotFound(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	reqUser := &gen.User{UserId: "ghost-id", FirstName: "Ghost"}
	req := &gen.UpdateUserRequest{User: reqUser}

	mockValidator.On("ValidateUpdate", req).Return(nil)
	mockStore.On("GetUserByID", mock.Anything, reqUser.UserId).Return(nil, status.Error(codes.NotFound, "not found"))

	server := NewServer(mockValidator, mockStore)

	res, err := server.UpdateUser(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestUpdateUserValidationFails(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	targetID := "update-user-id"
	reqUser := &gen.User{UserId: targetID}
	req := &gen.UpdateUserRequest{User: reqUser}
	validationErr := status.Error(codes.InvalidArgument, "last name cannot be empty")

	mockValidator.On("ValidateUpdate", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	res, err := server.UpdateUser(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	targetID := "delete-user-id"
	req := &gen.DeleteUserRequest{UserId: targetID}

	mockValidator.On("ValidateDelete", req).Return(nil)
	mockStore.On("DeleteUser", mock.Anything, targetID).Return(nil)

	server := NewServer(mockValidator, mockStore)

	_, err := server.DeleteUser(context.Background(), req)

	require.NoError(t, err)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestDeleteUserNotFound(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	userID := "non-existent-user-id"
	req := &gen.DeleteUserRequest{UserId: userID}

	mockValidator.On("ValidateDelete", req).Return(nil)
	mockStore.On("DeleteUser", mock.Anything, userID).Return(status.Error(codes.NotFound, "not found"))

	server := NewServer(mockValidator, mockStore)

	_, err := server.DeleteUser(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestDeleteUserValidationFails(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	targetID := "delete-user-id"
	req := &gen.DeleteUserRequest{UserId: targetID}
	validationErr := status.Error(codes.InvalidArgument, "invalid format")

	mockValidator.On("ValidateDelete", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	_, err := server.DeleteUser(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	mockValidator.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	mockValidator := new(MockUserValidator)
	mockStore := new(MockStore)

	domainUsers := []*domain.User{
		{ID: "id-1", FirstName: "User 1", LastName: "Last1", EmailAddress: "user1@example.com"},
		{ID: "id-2", FirstName: "User 2", LastName: "Last2", EmailAddress: "user2@example.com"},
	}
	mockStore.On("ListUsers", mock.Anything).Return(domainUsers, nil)

	server := NewServer(mockValidator, mockStore)

	req := &gen.ListUsersRequest{}
	res, err := server.ListUsers(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res.Users, 2)

	userIDs := []string{res.Users[0].UserId, res.Users[1].UserId}
	assert.Contains(t, userIDs, "id-1")
	assert.Contains(t, userIDs, "id-2")

	mockStore.AssertExpectations(t)
}
