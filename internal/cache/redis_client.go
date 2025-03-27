package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/enson89/sustainability-tracker-activity-service/internal/model"
	"github.com/go-redis/redis/v8"
)

const activityCacheKeyPrefix = "activity:"

// RedisClient wraps the redis.Client.
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient creates and pings a new Redis client.
func NewRedisClient(addr, password string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return &RedisClient{Client: client}, nil
}

// GetActivityCacheKey returns the cache key for a given activity ID.
func GetActivityCacheKey(id int64) string {
	return activityCacheKeyPrefix + strconv.FormatInt(id, 10)
}

// CacheActivity caches the activity data with the specified expiration.
func (r *RedisClient) CacheActivity(ctx context.Context, activity *model.Activity, expiration time.Duration) error {
	key := GetActivityCacheKey(activity.ID)
	data, err := json.Marshal(activity)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, data, expiration).Err()
}

// InvalidateActivityCache removes the cached activity.
func (r *RedisClient) InvalidateActivityCache(ctx context.Context, id int64) error {
	key := GetActivityCacheKey(id)
	return r.Client.Del(ctx, key).Err()
}

// GetCachedActivity retrieves the activity from the cache.
func (r *RedisClient) GetCachedActivity(ctx context.Context, id int64) (*model.Activity, error) {
	key := GetActivityCacheKey(id)
	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var activity model.Activity
	if err := json.Unmarshal([]byte(data), &activity); err != nil {
		return nil, err
	}
	return &activity, nil
}
