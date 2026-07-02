package issue

import (
	"context"
	"testing"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) ValidateCreate(req *gen.CreateIssueRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockValidator) ValidateUpdate(req, issue *gen.Issue) error {
	args := m.Called(req, issue)
	return args.Error(0)
}

func (m *MockValidator) ValidateGet(req *gen.ReadIssueRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockValidator) ValidateDelete(req *gen.DeleteIssueRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateIssue(ctx context.Context, issue *domain.Issue) error {
	return m.Called(ctx, issue).Error(0)
}

func (m *MockStore) GetIssueByID(ctx context.Context, id string) (*domain.Issue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Issue), args.Error(1)
}

func (m *MockStore) UpdateIssue(ctx context.Context, issue *domain.Issue) error {
	return m.Called(ctx, issue).Error(0)
}

func (m *MockStore) DeleteIssue(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockStore) ListIssues(ctx context.Context) ([]*domain.Issue, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Issue), args.Error(1)
}

func TestCreateIssue(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	req := &gen.CreateIssueRequest{
		Summary:     "Fix bug",
		Description: "Fix production issue",
		ProjectId:   "project-123",
		Status:      gen.Status_STATUS_NEW,
	}

	mockValidator.On("ValidateCreate", mock.Anything).Return(nil)
	mockStore.On("CreateIssue", mock.Anything, mock.AnythingOfType("*domain.Issue")).Return(nil)

	server := NewServer(mockValidator, mockStore)

	res, err := server.CreateIssue(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Issue)
	assert.NotEmpty(t, res.Issue.IssueId, "Server must generate ULID")
	assert.Equal(t, req.Summary, res.Issue.Summary)
	assert.Equal(t, gen.Resolution_RESOLUTION_RESOLUTION_UNSPECIFIED, res.Issue.Resolution)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestCreateIssueValidationFails(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	req := &gen.CreateIssueRequest{}
	validationErr := status.Error(codes.InvalidArgument, "invalid summary")

	mockValidator.On("ValidateCreate", mock.Anything).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	res, err := server.CreateIssue(context.Background(), req)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
}

func TestReadIssue(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	issueID := "existing-id"
	req := &gen.ReadIssueRequest{IssueId: issueID}

	domainIssue := &domain.Issue{
		ID:         issueID,
		Summary:    "Test",
		CreateDate: time.Now(),
		ModifyDate: time.Now(),
	}

	mockValidator.On("ValidateGet", req).Return(nil)
	mockStore.On("GetIssueByID", mock.Anything, issueID).Return(domainIssue, nil)

	server := NewServer(mockValidator, mockStore)

	res, err := server.ReadIssue(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, issueID, res.Issue.IssueId)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestReadIssueNotFound(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	issueID := "missing-id"
	req := &gen.ReadIssueRequest{IssueId: issueID}

	mockValidator.On("ValidateGet", req).Return(nil)
	mockStore.On("GetIssueByID", mock.Anything, issueID).Return(nil, status.Error(codes.NotFound, "not found"))

	server := NewServer(mockValidator, mockStore)

	res, err := server.ReadIssue(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestReadIssueValidationFails(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	issueID := "some-id"
	req := &gen.ReadIssueRequest{IssueId: issueID}
	validationErr := status.Error(codes.InvalidArgument, "bad request")

	mockValidator.On("ValidateGet", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	res, err := server.ReadIssue(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
}

func TestUpdateIssue(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	targetID := "update-id"
	reqIssue := &gen.Issue{IssueId: targetID, Summary: "Updated Summary"}

	existingIssue := &domain.Issue{
		ID:         targetID,
		Summary:    "Old Summary",
		CreateDate: time.Now().Add(-time.Hour),
		ModifyDate: time.Now().Add(-time.Hour),
	}

	mockStore.On("GetIssueByID", mock.Anything, targetID).Return(existingIssue, nil)
	mockValidator.On("ValidateUpdate", reqIssue, mock.Anything).Return(nil)
	mockStore.On("UpdateIssue", mock.Anything, mock.AnythingOfType("*domain.Issue")).Return(nil)

	server := NewServer(mockValidator, mockStore)

	req := &gen.UpdateIssueRequest{Issue: reqIssue}

	res, err := server.UpdateIssue(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, reqIssue.Summary, res.Issue.Summary)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestUpdateIssueNilEntity(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	server := NewServer(mockValidator, mockStore)

	req := &gen.UpdateIssueRequest{Issue: nil}

	res, err := server.UpdateIssue(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)
}

func TestUpdateIssueNotFound(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	reqIssue := &gen.Issue{IssueId: "ghost-id"}

	mockStore.On("GetIssueByID", mock.Anything, reqIssue.IssueId).Return(nil, status.Error(codes.NotFound, "not found"))

	server := NewServer(mockValidator, mockStore)

	req := &gen.UpdateIssueRequest{Issue: reqIssue}

	res, err := server.UpdateIssue(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, res)

	mockStore.AssertExpectations(t)
}

func TestUpdateIssueValidationFails(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	targetID := "update-id"
	reqIssue := &gen.Issue{IssueId: targetID}
	validationErr := status.Error(codes.InvalidArgument, "invalid status transition")

	existingIssue := &domain.Issue{
		ID:         targetID,
		Summary:    "Old Summary",
		CreateDate: time.Now().Add(-time.Hour),
		ModifyDate: time.Now().Add(-time.Hour),
	}

	mockStore.On("GetIssueByID", mock.Anything, targetID).Return(existingIssue, nil)
	mockValidator.On("ValidateUpdate", reqIssue, mock.Anything).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	req := &gen.UpdateIssueRequest{Issue: reqIssue}

	res, err := server.UpdateIssue(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, res)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestDeleteIssue(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	targetID := "delete-id"
	req := &gen.DeleteIssueRequest{IssueId: targetID}

	mockValidator.On("ValidateDelete", req).Return(nil)
	mockStore.On("DeleteIssue", mock.Anything, targetID).Return(nil)

	server := NewServer(mockValidator, mockStore)

	_, err := server.DeleteIssue(context.Background(), req)

	require.NoError(t, err)

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestDeleteIssueNotFound(t *testing.T) {

	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	issueID := "missing-id"
	req := &gen.DeleteIssueRequest{IssueId: issueID}

	mockValidator.On("ValidateDelete", req).Return(nil)
	mockStore.On("DeleteIssue", mock.Anything, issueID).Return(status.Error(codes.NotFound, "not found"))

	server := NewServer(mockValidator, mockStore)

	_, err := server.DeleteIssue(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())

	mockValidator.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestDeleteIssueValidationFails(t *testing.T) {
	mockValidator := new(MockValidator)
	mockStore := new(MockStore)

	targetID := "delete-id"
	req := &gen.DeleteIssueRequest{IssueId: targetID}
	validationErr := status.Error(codes.InvalidArgument, "invalid id")

	mockValidator.On("ValidateDelete", req).Return(validationErr)

	server := NewServer(mockValidator, mockStore)

	_, err := server.DeleteIssue(context.Background(), req)

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	mockValidator.AssertExpectations(t)
}
