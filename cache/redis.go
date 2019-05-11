package cache

import (
	"fmt"
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

type RedisErr struct {
	Msg string
}

func (e *RedisErr) Error() string {
	return fmt.Sprintf("Redis err: %s", e.Msg)

}

type RedisClient interface {
	Get(string, interface{}) error
	Set(*cache.Item) error
	Delete(string) error
}

type Redis struct {
	Codec *cache.Codec
}

func NewRedisClient(hosts map[string]string) RedisClient {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: hosts,
	})

	return &Redis{&cache.Codec{
		Redis: ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}}
}

func (rc *Redis) Get(key string, object interface{}) error {
	if err := rc.Codec.Get(key, object); err != nil {
		return &RedisErr{Msg: err.Error()}
	}
	return nil
}

func (rc *Redis) Set(item *cache.Item) error {
	if err := rc.Codec.Set(item); err != nil {
		return &RedisErr{Msg: err.Error()}
	}
	return nil
}

func (rc *Redis) Delete(key string) error {
	if err := rc.Codec.Delete(key); err != nil {
		return &RedisErr{Msg: err.Error()}
	}
	return nil
}
