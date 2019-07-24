package cache

import (
	"context"
	"go.elastic.co/apm/module/apmgoredis"
	"sync"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

var (
	instance Cache
	once     sync.Once
)

//Cache is an interface to handle cache
type Cache interface {
	Get(context.Context, string, interface{}) error
	Set(context.Context, string, interface{}, time.Duration) error
	Delete(context.Context, string) error
	Initialize(map[string]string)
	Disconnect()
}

type rImpl struct {
	codec *cache.Codec
	ring  *redis.Ring
}

//GetInstance to return a redis client
func GetInstance() Cache {
	once.Do(func() {
		if instance == nil {
			instance = &rImpl{}
		}
	})
	return instance
}

func (rc *rImpl) Initialize(hosts map[string]string) {
	rc.ring = redis.NewRing(&redis.RingOptions{
		Addrs: hosts,
	})
	rc.codec = &cache.Codec{
		Redis: rc.ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func (rc *rImpl) Get(ctx context.Context, key string, object interface{}) error {
	var b []byte
	client := apmgoredis.Wrap(rc.ring).WithContext(ctx)
	b, err := client.Get(key).Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(b, object)
}

func (rc *rImpl) Set(ctx context.Context, k string, obj interface{}, d time.Duration) error {
	client := apmgoredis.Wrap(rc.ring).WithContext(ctx)
	if err := client.Set(k, obj, d).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *rImpl) Delete(ctx context.Context, key string) error {
	client := apmgoredis.Wrap(rc.ring).WithContext(ctx)
	if err := client.Del(key).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *rImpl) Disconnect() {
	_ = rc.ring.Close()
}
