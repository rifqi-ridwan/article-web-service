package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepository interface {
	Read(ctx context.Context, key string) (string, error)
	ReadJSON(ctx context.Context, key string, obj interface{}) error
	Write(ctx context.Context, key, value string, exp time.Duration) error
	WriteJSON(ctx context.Context, key string, obj interface{}, exp time.Duration) error
}

type repository struct {
	client *redis.Client
}

func NewCache(client *redis.Client) CacheRepository {
	return &repository{client}
}

func (r *repository) Read(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *repository) ReadJSON(ctx context.Context, key string, obj interface{}) error {
	result, err := r.Read(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), obj)
}

func (r *repository) Write(ctx context.Context, key, value string, exp time.Duration) error {
	return r.client.Set(ctx, key, value, exp).Err()
}

func (r *repository) WriteJSON(ctx context.Context, key string, obj interface{}, exp time.Duration) error {
	result, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return r.Write(ctx, key, string(result), exp)
}
