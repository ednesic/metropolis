package cache

import (
"github.com/go-redis/cache"
"github.com/go-redis/redis"
"github.com/vmihailenco/msgpack"
)

type RedisClient interface {
	Get(string, interface{}) error
	Set(*cache.Item) error
	Delete(string) error
}

type Redis struct {
	Codec *cache.Codec
}

func NewRedisClient(hosts map[string]string) *cache.Codec {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: hosts,
	})

	return &cache.Codec{
		Redis: ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func (rc *Redis) Get(key string, object interface{}) error {
	return rc.Codec.Get(key, object)
}

func (rc *Redis) Set(item *cache.Item) error {
	return rc.Codec.Set(item)
}

func (rc *Redis) Delete(key string) error {
	return rc.Codec.Delete(key)
}
