package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	ctx context.Context
	rdb *redis.Client
}

func NewRedisStore() (*RedisStore, error) {
	var ctx = context.Background()

	Addr := os.Getenv("REDIS_URI")

	fmt.Printf("Addr is : %v \n ", Addr)
	rdb := redis.NewClient(&redis.Options{
		Addr: Addr,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStore{
		rdb: rdb,
		ctx: ctx,
	}, nil
}

func (rc *RedisStore) SetEmailVerificationCode(email, s string) error {
	// ev- stands for email verification
	key := fmt.Sprintf("ev-%v", email)
	err := rc.rdb.Set(rc.ctx, key, s, time.Second*60*60*24).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisStore) GetEmailVerificationCode(email string) (string, error) {
	// ev- stands for email verification
	key := fmt.Sprintf("ev-%v", email)
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (rc *RedisStore) DeleteEmailVerificationCode(email string) error {
	// ev- stands for email verification
	key := fmt.Sprintf("ev-%v", email)
	_, err := rc.rdb.Del(rc.ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
