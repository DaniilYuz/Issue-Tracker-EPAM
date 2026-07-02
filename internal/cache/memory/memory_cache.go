package memory

import "git.epam.com/go-language-global-mentoring-program/internal/cache"

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
