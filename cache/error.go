package cache

import "fmt"

//RedisErr is the err that every cache err should return.
type RedisErr struct {
	Msg string
}

func (e *RedisErr) Error() string {
	return fmt.Sprintf("redis err: %s", e.Msg)
}
