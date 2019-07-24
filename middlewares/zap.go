package middlewares

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strconv"
	"time"

	"go.elastic.co/apm/module/apmzap"
)


var logger = zap.NewExample(zap.WrapCore((&apmzap.Core{}).WrapCore))

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
			}
			stop := time.Now()

			traceContextFields := apmzap.TraceContext(req.Context())
			fields := zap.Fields(zap.Error(err),
				zap.String("id", req.Header.Get(echo.HeaderXRequestID)),
				zap.String("remote_ip", c.RealIP()),
				zap.String("host", req.Host),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("user_agent", req.UserAgent()),
				zap.Int("status", res.Status),
				zap.Duration("latency", stop.Sub(start)),
				zap.String("bytes_in", req.Header.Get(echo.HeaderContentLength)),
				zap.String("bytes_out", strconv.FormatInt(res.Size, 10)),
			)

			switch {
			case res.Status >= 500:
				logger.With(traceContextFields...).WithOptions(fields).Error("")
			case res.Status >= 400:
				logger.With(traceContextFields...).WithOptions(fields).Warn("")
			case res.Status >= 200:
				logger.With(traceContextFields...).WithOptions(fields).Info("")
			}
			return
		}
	}
}