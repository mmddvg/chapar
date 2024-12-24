package redisregistery

import (
	"context"
	"fmt"
	"mmddvg/chapar/pkg/errs"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

type RedisRegistery struct {
	client *redis.Client
}

func NewRedisRegister() (*RedisRegistery, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     lo.Ternary(os.Getenv("REDIS_URI") != "", os.Getenv("REDIS_URI"), "localhost:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	res := rdb.Ping(context.Background())

	return &RedisRegistery{
		client: rdb,
	}, res.Err()
}

func (r *RedisRegistery) Register(userId uint64, serverId uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("user:%d", userId)
	member := fmt.Sprintf("%d", serverId)

	_, err := r.client.SAdd(ctx, key, member).Result()
	if err != nil {
		return errs.NewUnexpected(fmt.Errorf("failed to register server %d for user %d: %w", serverId, userId, err))
	}

	return nil
}

func (r *RedisRegistery) Retrive(userId uint64) ([]uint, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:%d", userId)

	serverIds, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errs.NewUnexpected(fmt.Errorf("failed to retrieve servers for user %d: %w", userId, err))
	}

	var results []uint
	for _, id := range serverIds {
		var serverId uint
		_, err := fmt.Sscanf(id, "%d", &serverId)
		if err != nil {
			return nil, errs.NewUnexpected(fmt.Errorf("failed to parse server ID %s: %w", id, err))
		}
		results = append(results, serverId)
	}

	return results, nil
}

func (r *RedisRegistery) UnRegister(userId uint64, serverId uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("user:%d", userId)
	member := fmt.Sprintf("%d", serverId)

	_, err := r.client.SRem(ctx, key, member).Result()
	if err != nil {
		return errs.NewUnexpected(fmt.Errorf("failed to unregister server %d for user %d: %w", serverId, userId, err))
	}

	return nil
}
