package memory

import "github.com/DaniilYuz/Issue-Tracker-EPAM/internal/cache"

type MemoryCache struct {
	*UserCache
	*ProjectCache
	*IssueCache
}

func NewCache() cache.Cache {
	store := NewStore()
	return &MemoryCache{
		UserCache:    NewUserCache(store),
		ProjectCache: NewProjectCache(store),
		IssueCache:   NewIssueCache(store),
	}
}
