package cache

import (
	"github.com/go-redis/cache"
	"github.com/stretchr/testify/mock"
)

type RedisMock struct {
	mock.Mock
}

func (rc *RedisMock) Get(key string, object interface{}) error {
	args := rc.Called(key, object)
	return args.Error(0)
}

func (rc *RedisMock) Set(item *cache.Item) error {
	args := rc.Called(item)
	return args.Error(0)
}

func (rc *RedisMock) Delete(key string) error {
	args := rc.Called(key)
	return args.Error(0)
}
