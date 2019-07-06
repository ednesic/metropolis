package cache

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type RedisMock struct {
	mock.Mock
}

func (rc *RedisMock) Initialize(map[string]string) {
	instance = rc
}

func (rc *RedisMock) Get(key string, object interface{}) error {
	args := rc.Called(key, object)
	return args.Error(0)
}

func (rc *RedisMock) Set(k string, obj interface{}, d time.Duration) error {
	args := rc.Called(k, obj, d)
	return args.Error(0)
}

func (rc *RedisMock) Delete(key string) error {
	args := rc.Called(key)
	return args.Error(0)
}

func (rc *RedisMock) Disconnect() {}
