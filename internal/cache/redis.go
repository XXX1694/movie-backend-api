package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang/pkg/modules"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisCache{client: client}
}

func (c *RedisCache) GetUser(ctx context.Context, id int) (*modules.User, error) {
	key := fmt.Sprintf("user:%d", id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var user modules.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *RedisCache) SetUser(ctx context.Context, user *modules.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, 5*time.Minute).Err()
}

func (c *RedisCache) DeleteUser(ctx context.Context, id int) error {
	key := fmt.Sprintf("user:%d", id)
	return c.client.Del(ctx, key).Err()
}
