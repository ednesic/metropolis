package opentracing

import (
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

const defaultKey = "github.com/hb-go/echo-web/middleware/opentracing"

func OpenTracing(comp string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var span opentracing.Span
			opName := comp + ":" + c.Request().URL.Path
			wireContext, err := opentracing.GlobalTracer().Extract(
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(c.Request().Header))
			if err != nil {
				span = opentracing.StartSpan(opName)
			} else {
				span = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
			}

			defer span.Finish()
			c.Set(defaultKey, span)

			span.SetTag("component", comp)
			span.SetTag("span.kind", "server")
			span.SetTag("http.url", c.Request().Host+c.Request().RequestURI)
			span.SetTag("http.method", c.Request().Method)

			if err := next(c); err != nil {
				span.SetTag("error", true)
				c.Error(err)
			}
			span.SetTag("error", false)
			span.SetTag("http.status_code", c.Response().Status)

			return nil
		}
	}
}
