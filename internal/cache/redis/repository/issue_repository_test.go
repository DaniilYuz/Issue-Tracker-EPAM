package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

func TestNewIssueRepository(t *testing.T) {
	client, _ := setupTestRedis(t)

	repo := NewIssueRepository(client)

	assert.NotNil(t, repo)
	assert.IsType(t, &issueRepository{}, repo)
}

func TestSetIssue(t *testing.T) {
	client, _ := setupTestRedis(t)
	repo := NewIssueRepository(client)

	issue := &domain.Issue{
		ID:          "test-id",
		Summary:     "Test Issue",
		Description: "Test Description",
		Status:      domain.IssueStatusNew,
		Type:        domain.IssueTypeBug,
		Priority:    domain.PriorityCritical,
		Resolution:  domain.ResolutionUnspecified,
		ProjectID:   "project-1",
		AssigneeID:  "user-1",
		CreateDate:  time.Now(),
		ModifyDate:  time.Now(),
	}

	ttl := 15 * time.Minute

	err := repo.SetIssue(context.Background(), ttl, issue)

	assert.NoError(t, err)

	result, err := repo.GetIssueByID(context.Background(), issue.ID)
	assert.NoError(t, err)
	assert.Equal(t, issue.ID, result.ID)
	assert.Equal(t, issue.Summary, result.Summary)
}

func TestSetIssueError(t *testing.T) {
	client, mr := setupTestRedis(t)
	repo := NewIssueRepository(client)

	issue := &domain.Issue{
		ID:      "test-id",
		Summary: "Test Issue",
		Status:  domain.IssueStatusNew,
	}

	mr.Close()

	err := repo.SetIssue(context.Background(), time.Minute, issue)

	assert.Error(t, err)
}

func TestSetIssues(t *testing.T) {
	client, _ := setupTestRedis(t)
	repo := NewIssueRepository(client)

	issues := []*domain.Issue{
		{
			ID:      "issue-1",
			Summary: "Issue 1",
			Status:  domain.IssueStatusNew,
			Type:    domain.IssueTypeBug,
		},
		{
			ID:      "issue-2",
			Summary: "Issue 2",
			Status:  domain.IssueStatusAssigned,
			Type:    domain.IssueTypeFeature,
		},
	}

	ttl := 15 * time.Minute

	err := repo.SetIssues(context.Background(), ttl, issues)
	assert.NoError(t, err)

	for _, issue := range issues {
		result, err := repo.GetIssueByID(context.Background(), issue.ID)
		assert.NoError(t, err)
		assert.Equal(t, issue.ID, result.ID)
		assert.Equal(t, issue.Summary, result.Summary)
	}
}

func TestSetIssuesError(t *testing.T) {
	client, mr := setupTestRedis(t)
	repo := NewIssueRepository(client)

	issues := []*domain.Issue{
		{
			ID:      "issue-1",
			Summary: "Issue 1",
			Status:  domain.IssueStatusNew,
		},
	}

	mr.Close()

	err := repo.SetIssues(context.Background(), time.Minute, issues)

	assert.Error(t, err)
}

func TestGetIssueByID(t *testing.T) {
	client, _ := setupTestRedis(t)
	repo := NewIssueRepository(client)

	issue := &domain.Issue{
		ID:       "test-id",
		Summary:  "Test Issue",
		Status:   domain.IssueStatusNew,
		Type:     domain.IssueTypeBug,
		Priority: domain.PriorityMajor,
	}

	err := repo.SetIssue(context.Background(), time.Hour, issue)
	require.NoError(t, err)

	result, err := repo.GetIssueByID(context.Background(), issue.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, issue.ID, result.ID)
	assert.Equal(t, issue.Summary, result.Summary)
	assert.Equal(t, issue.Status, result.Status)
}

func TestGetIssueByIDNotFound(t *testing.T) {
	client, _ := setupTestRedis(t)
	repo := NewIssueRepository(client)

	result, err := repo.GetIssueByID(context.Background(), "missing-id")

	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, result)
}

func TestGetIssueByIDError(t *testing.T) {
	client, mr := setupTestRedis(t)
	repo := NewIssueRepository(client)

	mr.Close()

	result, err := repo.GetIssueByID(context.Background(), "test-id")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetIssueByIDUnmarshalError(t *testing.T) {
	client, mr := setupTestRedis(t)
	repo := NewIssueRepository(client)

	err := mr.Set("issues:test-id", "invalid json data")
	require.NoError(t, err)

	result, err := repo.GetIssueByID(context.Background(), "test-id")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteIssue(t *testing.T) {
	client, _ := setupTestRedis(t)
	repo := NewIssueRepository(client)

	issue := &domain.Issue{
		ID:      "test-id",
		Summary: "Test Issue",
		Status:  domain.IssueStatusNew,
	}

	err := repo.SetIssue(context.Background(), time.Hour, issue)
	require.NoError(t, err)

	result, err := repo.GetIssueByID(context.Background(), issue.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	err = repo.DeleteIssue(context.Background(), issue.ID)
	assert.NoError(t, err)

	result, err = repo.GetIssueByID(context.Background(), issue.ID)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, result)
}

func TestDeleteIssueError(t *testing.T) {
	client, mr := setupTestRedis(t)
	repo := NewIssueRepository(client)

	mr.Close()

	err := repo.DeleteIssue(context.Background(), "test-id")

	assert.Error(t, err)
}

func TestDeleteAllIssues(t *testing.T) {
	client, mr := setupTestRedis(t)
	repo := NewIssueRepository(client)

	issues := []*domain.Issue{
		{ID: "issue-1", Summary: "Issue 1", Status: domain.IssueStatusNew},
		{ID: "issue-2", Summary: "Issue 2", Status: domain.IssueStatusAssigned},
		{ID: "issue-3", Summary: "Issue 3", Status: domain.IssueStatusInProgress},
	}

	for _, issue := range issues {
		err := repo.SetIssue(context.Background(), time.Hour, issue)
		require.NoError(t, err)
	}

	assert.True(t, mr.Exists("issues"))
	members, _ := mr.SMembers("issues")
	assert.Len(t, members, 3)
	assert.Contains(t, members, "issue-1")
	assert.Contains(t, members, "issue-2")
	assert.Contains(t, members, "issue-3")

	err := repo.DeleteAllIssues(context.Background())
	assert.NoError(t, err)

	for _, issue := range issues {
		result, err := repo.GetIssueByID(context.Background(), issue.ID)
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
		assert.Nil(t, result)
	}

	assert.False(t, mr.Exists("issues"))
}

func TestDeleteAllIssuesError(t *testing.T) {
	client, mr := setupTestRedis(t)
	repo := NewIssueRepository(client)

	mr.Close()

	err := repo.DeleteAllIssues(context.Background())

	assert.Error(t, err)
}

func TestDeleteAllIssuesEmpty(t *testing.T) {
	client, _ := setupTestRedis(t)
	repo := NewIssueRepository(client)

	err := repo.DeleteAllIssues(context.Background())

	assert.NoError(t, err)
}

func TestIssueExpiration(t *testing.T) {
	client, mr := setupTestRedis(t)
	repo := NewIssueRepository(client)

	issue := &domain.Issue{
		ID:      "test-id",
		Summary: "Test Issue",
		Status:  domain.IssueStatusNew,
	}

	err := repo.SetIssue(context.Background(), 100*time.Millisecond, issue)
	require.NoError(t, err)

	result, err := repo.GetIssueByID(context.Background(), issue.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	mr.FastForward(200 * time.Millisecond)

	_, err = repo.GetIssueByID(context.Background(), issue.ID)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisState(t *testing.T) {
	client, mr := setupTestRedis(t)
	repo := NewIssueRepository(client)

	issue := &domain.Issue{
		ID:      "test-id",
		Summary: "Test Issue",
		Status:  domain.IssueStatusNew,
	}

	err := repo.SetIssue(context.Background(), time.Hour, issue)
	require.NoError(t, err)

	assert.True(t, mr.Exists("issues:test-id"))
	assert.True(t, mr.Exists("issues"))

	members, _ := mr.SMembers("issues")
	assert.Contains(t, members, "test-id")
}
