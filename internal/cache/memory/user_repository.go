package memory

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/cache/utils"
	"github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"
)

type UserCache struct {
	store *Store
}

func NewUserCache(store *Store) *UserCache {
	return &UserCache{store: store}
}

func (c *UserCache) SetUser(ctx context.Context, ttl time.Duration, user *domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	c.store.set(utils.UserKey(user.ID), data, ttl)
	return nil
}

func (c *UserCache) SetUsers(ctx context.Context, ttl time.Duration, users []*domain.User) error {
	for _, user := range users {
		if err := c.SetUser(ctx, ttl, user); err != nil {
			return err
		}
	}
	return nil
}

func (c *UserCache) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	data, err := c.store.get(utils.UserKey(id))
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *UserCache) DeleteUser(ctx context.Context, id string) error {
	c.store.delete(utils.UserKey(id))
	return nil
}

func (c *UserCache) DeleteAllUsers(ctx context.Context) error {
	c.store.deleteByPrefix("users:")
	return nil
}
