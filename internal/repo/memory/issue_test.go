package memory

import (
	"context"
	"testing"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateIssue(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	issue := &domain.Issue{
		ID:          "issue-123",
		Summary:     "Test Issue",
		Description: "A test issue",
		ProjectID:   "project-123",
		AssigneeID:  "user-123",
		CreateDate:  time.Now(),
		ModifyDate:  time.Now(),
	}

	err = store.CreateIssue(context.Background(), issue)

	require.NoError(t, err)

	// Verify issue was created
	retrieved, err := store.GetIssueByID(context.Background(), "issue-123")
	require.NoError(t, err)
	assert.Equal(t, issue.ID, retrieved.ID)
	assert.Equal(t, issue.Summary, retrieved.Summary)
	assert.Equal(t, issue.Description, retrieved.Description)
	assert.Equal(t, issue.ProjectID, retrieved.ProjectID)
	assert.Equal(t, issue.AssigneeID, retrieved.AssigneeID)
}

func TestGetIssueByID(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	issue := &domain.Issue{
		ID:          "issue-456",
		Summary:     "Another Issue",
		Description: "Another test issue",
		ProjectID:   "project-456",
		AssigneeID:  "user-456",
		CreateDate:  time.Now(),
		ModifyDate:  time.Now(),
	}

	// Create issue first
	err = store.CreateIssue(context.Background(), issue)
	require.NoError(t, err)

	retrieved, err := store.GetIssueByID(context.Background(), "issue-456")

	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, issue.ID, retrieved.ID)
	assert.Equal(t, issue.Summary, retrieved.Summary)
	assert.Equal(t, issue.Description, retrieved.Description)
	assert.Equal(t, issue.ProjectID, retrieved.ProjectID)
	assert.Equal(t, issue.AssigneeID, retrieved.AssigneeID)
}

func TestGetIssueByIDNotFound(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	retrieved, err := store.GetIssueByID(context.Background(), "non-existent-issue")

	require.Error(t, err)
	assert.Nil(t, retrieved)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "issue with this id did not found")
}

func TestUpdateIssue(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	// Create original issue
	originalIssue := &domain.Issue{
		ID:          "issue-789",
		Summary:     "Original Issue",
		Description: "Original description",
		ProjectID:   "project-789",
		AssigneeID:  "user-789",
		CreateDate:  time.Now(),
		ModifyDate:  time.Now(),
	}
	err = store.CreateIssue(context.Background(), originalIssue)
	require.NoError(t, err)

	// Update issue
	updatedIssue := &domain.Issue{
		ID:          "issue-789",
		Summary:     "Updated Issue",
		Description: "Updated description",
		ProjectID:   "project-789",
		AssigneeID:  "user-updated",
		CreateDate:  originalIssue.CreateDate,
		ModifyDate:  time.Now(),
	}

	err = store.UpdateIssue(context.Background(), updatedIssue)

	require.NoError(t, err)

	// Verify update
	retrieved, err := store.GetIssueByID(context.Background(), "issue-789")
	require.NoError(t, err)
	assert.Equal(t, "Updated Issue", retrieved.Summary)
	assert.Equal(t, "Updated description", retrieved.Description)
	assert.Equal(t, "user-updated", retrieved.AssigneeID)
}

func TestUpdateIssueNotFound(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	issue := &domain.Issue{
		ID:          "non-existent-issue",
		Summary:     "Ghost Issue",
		Description: "This issue doesn't exist",
		ProjectID:   "project-ghost",
		AssigneeID:  "user-ghost",
		CreateDate:  time.Now(),
		ModifyDate:  time.Now(),
	}

	err = store.UpdateIssue(context.Background(), issue)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "issue not found")
}

func TestDeleteIssue(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	// Create issue first
	issue := &domain.Issue{
		ID:          "issue-delete",
		Summary:     "Delete Me",
		Description: "This issue will be deleted",
		ProjectID:   "project-delete",
		AssigneeID:  "user-delete",
		CreateDate:  time.Now(),
		ModifyDate:  time.Now(),
	}
	err = store.CreateIssue(context.Background(), issue)
	require.NoError(t, err)

	err = store.DeleteIssue(context.Background(), "issue-delete")

	require.NoError(t, err)

	// Verify issue was deleted
	retrieved, err := store.GetIssueByID(context.Background(), "issue-delete")
	require.Error(t, err)
	assert.Nil(t, retrieved)
}

func TestListIssues(t *testing.T) {
	store, err := NewStore()
	require.NoError(t, err)

	// Create multiple issues
	issues := []*domain.Issue{
		{
			ID:          "issue-1",
			Summary:     "Issue One",
			Description: "First issue",
			ProjectID:   "project-1",
			AssigneeID:  "user-1",
			CreateDate:  time.Now(),
			ModifyDate:  time.Now(),
		},
		{
			ID:          "issue-2",
			Summary:     "Issue Two",
			Description: "Second issue",
			ProjectID:   "project-2",
			AssigneeID:  "user-2",
			CreateDate:  time.Now(),
			ModifyDate:  time.Now(),
		},
		{
			ID:          "issue-3",
			Summary:     "Issue Three",
			Description: "Third issue",
			ProjectID:   "project-3",
			AssigneeID:  "user-3",
			CreateDate:  time.Now(),
			ModifyDate:  time.Now(),
		},
	}

	for _, issue := range issues {
		err = store.CreateIssue(context.Background(), issue)
		require.NoError(t, err)
	}

	retrieved, err := store.ListIssues(context.Background())

	require.NoError(t, err)
	assert.Len(t, retrieved, 3)

	// Verify all issues are present
	issueIDs := make([]string, len(retrieved))
	for i, issue := range retrieved {
		issueIDs[i] = issue.ID
	}

	assert.Contains(t, issueIDs, "issue-1")
	assert.Contains(t, issueIDs, "issue-2")
	assert.Contains(t, issueIDs, "issue-3")
}
