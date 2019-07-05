package cache

import "fmt"

type RedisErr struct {
	Msg string
}

func (e *RedisErr) Error() string {
	return fmt.Sprintf("redis err: %s", e.Msg)
}