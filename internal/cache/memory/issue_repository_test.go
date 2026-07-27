package memory

import (
	"context"
	"testing"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueCache_SetGetDeleteIssue(t *testing.T) {
	store := NewStore()
	cache := NewIssueCache(store)
	ctx := context.Background()

	issue := &domain.Issue{
		ID:          "i1",
		CreateDate:  time.Now(),
		ModifyDate:  time.Now(),
		Summary:     "Test summary",
		Description: "Test description",
		Status:      domain.IssueStatusNew,
		Resolution:  domain.ResolutionUnspecified,
		Type:        domain.IssueTypeBug,
		Priority:    domain.PriorityMajor,
		ProjectID:   "p1",
		AssigneeID:  "u1",
	}
	err := cache.SetIssue(ctx, time.Minute, issue)
	require.NoError(t, err)

	got, err := cache.GetIssueByID(ctx, "i1")
	require.NoError(t, err)
	assert.Equal(t, issue.ID, got.ID)
	assert.Equal(t, issue.Summary, got.Summary)
	assert.Equal(t, issue.Description, got.Description)
	assert.Equal(t, issue.Status, got.Status)
	assert.Equal(t, issue.Resolution, got.Resolution)
	assert.Equal(t, issue.Type, got.Type)
	assert.Equal(t, issue.Priority, got.Priority)
	assert.Equal(t, issue.ProjectID, got.ProjectID)
	assert.Equal(t, issue.AssigneeID, got.AssigneeID)

	err = cache.DeleteIssue(ctx, "i1")
	require.NoError(t, err)

	got, err = cache.GetIssueByID(ctx, "i1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestIssueCache_SetIssues_DeleteAllIssues(t *testing.T) {
	store := NewStore()
	cache := NewIssueCache(store)
	ctx := context.Background()

	issues := []*domain.Issue{
		{
			ID:          "i1",
			CreateDate:  time.Now(),
			ModifyDate:  time.Now(),
			Summary:     "Issue 1",
			Description: "Desc 1",
			Status:      domain.IssueStatusNew,
			Resolution:  domain.ResolutionUnspecified,
			Type:        domain.IssueTypeBug,
			Priority:    domain.PriorityMajor,
			ProjectID:   "p1",
			AssigneeID:  "u1",
		},
		{
			ID:          "i2",
			CreateDate:  time.Now(),
			ModifyDate:  time.Now(),
			Summary:     "Issue 2",
			Description: "Desc 2",
			Status:      domain.IssueStatusAssigned,
			Resolution:  domain.ResolutionFixed,
			Type:        domain.IssueTypeFeature,
			Priority:    domain.PriorityCritical,
			ProjectID:   "p2",
			AssigneeID:  "u2",
		},
	}
	err := cache.SetIssues(ctx, time.Minute, issues)
	require.NoError(t, err)

	got1, err := cache.GetIssueByID(ctx, "i1")
	require.NoError(t, err)
	assert.Equal(t, "Issue 1", got1.Summary)
	assert.Equal(t, "Desc 1", got1.Description)

	got2, err := cache.GetIssueByID(ctx, "i2")
	require.NoError(t, err)
	assert.Equal(t, "Issue 2", got2.Summary)
	assert.Equal(t, "Desc 2", got2.Description)

	err = cache.DeleteAllIssues(ctx)
	require.NoError(t, err)

	got1, err = cache.GetIssueByID(ctx, "i1")
	assert.Error(t, err)
	assert.Nil(t, got1)

	got2, err = cache.GetIssueByID(ctx, "i2")
	assert.Error(t, err)
	assert.Nil(t, got2)
}

func TestIssueCache_TTL(t *testing.T) {
	store := NewStore()
	cache := NewIssueCache(store)
	ctx := context.Background()

	issue := &domain.Issue{
		ID:          "i1",
		CreateDate:  time.Now(),
		ModifyDate:  time.Now(),
		Summary:     "Test summary",
		Description: "Test description",
		Status:      domain.IssueStatusNew,
		Resolution:  domain.ResolutionUnspecified,
		Type:        domain.IssueTypeBug,
		Priority:    domain.PriorityMajor,
		ProjectID:   "p1",
		AssigneeID:  "u1",
	}
	err := cache.SetIssue(ctx, 10*time.Millisecond, issue)
	require.NoError(t, err)

	time.Sleep(20 * time.Millisecond)
	got, err := cache.GetIssueByID(ctx, "i1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestIssueCache_NotFound(t *testing.T) {
	store := NewStore()
	cache := NewIssueCache(store)
	ctx := context.Background()

	got, err := cache.GetIssueByID(ctx, "not-exist")
	assert.Error(t, err)
	assert.Nil(t, got)
}
