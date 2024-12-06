package redis

import (
	"fmt"
	"time"
)

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
