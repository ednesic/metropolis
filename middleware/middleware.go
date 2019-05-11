package middleware

import (
	"github.com/labstack/echo/v4"
)

const RedisContext = "redisContext"

func RedisWarn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		warn := c.Get(RedisContext)
		if warn != nil {
			c.Logger().Warn(warn)
		}
		return err
	}
}