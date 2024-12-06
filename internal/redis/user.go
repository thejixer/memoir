package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thejixer/memoir/internal/models"
)

func (rc *RedisStore) CacheUser(u *models.User) error {
	key := fmt.Sprintf("u-%v", u.ID)

	st, merr := json.Marshal(u)
	if merr != nil {
		return merr
	}

	err := rc.rdb.Set(rc.ctx, key, string(st), time.Second*60*10).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisStore) GetUser(id int) *models.User {

	key := fmt.Sprintf("u-%v", id)
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		return nil
	}
	var u models.User

	err = json.Unmarshal([]byte(val), &u)

	if err != nil {
		return nil
	}

	return &u
}
func (rc *RedisStore) DelUser(id int) error {
	key := fmt.Sprintf("u-%v", id)
	_, err := rc.rdb.Del(rc.ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
