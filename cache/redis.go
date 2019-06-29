package cache

import (
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"sync"
)

var instance *Redis
var once sync.Once

type RedisClient interface {
	Get(string, interface{}) error
	Set(*cache.Item) error
	Delete(string) error
	Initialize(map[string]string)
	Disconnect()
}

type Redis struct {
	Codec *cache.Codec
	ring *redis.Ring
}

func GetInstance() RedisClient {
	once.Do(func() {
		instance = &Redis{}
	})
	return instance
}

func (rc *Redis) Initialize(hosts map[string]string) {
	instance.ring = redis.NewRing(&redis.RingOptions{
		Addrs: hosts,
	})
	instance.Codec = &cache.Codec{
		Redis: instance.ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
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

func (rc * Redis) Disconnect() {
	_ = rc.ring.Close()
}