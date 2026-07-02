package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(config *ConfigRedis) (*redis.Client, error) {
	options := &redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	}
	client := redis.NewClient(options)

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	if pong != "PONG" {
		return nil, fmt.Errorf("unexpected ping responce: %s", pong)
	}
	return client, nil
}
