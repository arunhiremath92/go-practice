package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

// NewRedisStore initializes the connection
func NewRedisStore(opt *redis.Options) *RedisStore {
	client := redis.NewClient(opt)
	return &RedisStore{client: client}
}

func (r *RedisStore) GetObject(key string) (string, error) {
	keyVal := r.client.Get(context.Background(), key)
	if keyVal.Err() != nil {
		return "", fmt.Errorf("object not found in the db %s", keyVal.Err().Error())
	}
	return keyVal.Result()
}

func (r *RedisStore) SetObject(key string, value string) error {
	if keyVal := r.client.Get(context.Background(), key); keyVal != nil {
		fmt.Println("rewriting the key object")
	}
	if err := r.client.Set(context.Background(), key, string(value), 0).Err(); err != nil {
		return fmt.Errorf("faield to store object in to keystore %s", err)
	}
	return nil
}
