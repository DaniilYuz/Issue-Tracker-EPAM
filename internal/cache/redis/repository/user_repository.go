package repository

import (
	"context"
	"encoding/json"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/cache"
	"git.epam.com/go-language-global-mentoring-program/internal/cache/utils"
	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"github.com/redis/go-redis/v9"
)

type userRepository struct {
	redis redis.Cmdable
}

func NewUserRepository(redis *redis.Client) cache.UserCacheRepository {
	return &userRepository{redis: redis}
}

func (r *userRepository) SetUser(ctx context.Context, ttl time.Duration, user *domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	pipe := r.redis.TxPipeline()

	pipe.Set(ctx, utils.UserKey(user.ID), data, ttl)
	pipe.SAdd(ctx, "users", user.ID)

	_, err = pipe.Exec(ctx)

	return err
}

func (r *userRepository) SetUsers(ctx context.Context, ttl time.Duration, users []*domain.User) error {
	pipe := r.redis.TxPipeline()
	for _, user := range users {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		pipe.Set(ctx, utils.UserKey(user.ID), data, ttl)
		pipe.SAdd(ctx, "users", user.ID)
	}

	_, err := pipe.Exec(ctx)

	return err
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	val, err := r.redis.Get(ctx, utils.UserKey(id)).Result()
	if err == redis.Nil {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	pipe := r.redis.TxPipeline()
	pipe.Del(ctx, utils.UserKey(id))
	pipe.SRem(ctx, "users", id)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *userRepository) DeleteAllUsers(ctx context.Context) error {
	ids, err := r.redis.SMembers(ctx, "users").Result()
	if err != nil {
		return err
	}

	pipe := r.redis.TxPipeline()
	for _, id := range ids {
		pipe.Del(ctx, utils.UserKey(id))
	}
	pipe.Del(ctx, "users")
	_, err = pipe.Exec(ctx)
	return err
}
