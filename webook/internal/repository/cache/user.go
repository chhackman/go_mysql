package cache

import (
	"awesomeProject/webook/internal/domain"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

//NewUserCache
//A用到了B, B一定是接口
//A用到了B，B一定是A字段
//A用到了B,A绝对不初始化，而是外面注入

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

//只要error 为nil,就认为缓存里面有数据
//如果没有数据,返回一个特定的error
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	//数据不存在,err =redis.Nil
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}
func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
