package metrics

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	//PrometheusConfig promotheus configuration
	PrometheusConfig struct {
		Skipper   middleware.Skipper
		Namespace string
	}
)

var (
	//DefaultPrometheusConfig default prometheus configuration
	DefaultPrometheusConfig = PrometheusConfig{
		Skipper:   middleware.DefaultSkipper,
		Namespace: "echo",
	}
)

var (
	echoReqQPS      *prometheus.CounterVec
	echoReqDuration *prometheus.SummaryVec
	echoOutBytes    prometheus.Summary
)

func initCollector(namespace string) {
	echoReqQPS = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_request_total",
			Help:      "HTTP requests processed.",
		},
		[]string{"code", "method", "host", "path"},
	)
	echoReqDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
		},
		[]string{"method", "host", "url"},
	)
	echoOutBytes = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP response bytes.",
		},
	)
	prometheus.MustRegister(echoReqQPS, echoReqDuration, echoOutBytes)
}

//NewMetric is a middleware to get information for prometheus
func NewMetric() echo.MiddlewareFunc {
	return NewMetricWithConfig(DefaultPrometheusConfig)
}

//NewMetricWithConfig is a middleware to get information for prometheus. In this method is possible to pass config.
func NewMetricWithConfig(config PrometheusConfig) echo.MiddlewareFunc {
	initCollector(config.Namespace)
	if config.Skipper == nil {
		config.Skipper = DefaultPrometheusConfig.Skipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}
			status := strconv.Itoa(res.Status)
			elapsed := time.Since(start).Seconds()
			bytesOut := float64(res.Size)
			echoReqQPS.WithLabelValues(status, req.Method, req.Host, c.Path()).Inc()
			echoReqDuration.WithLabelValues(req.Method, req.Host, c.Path()).Observe(elapsed)
			echoOutBytes.Observe(bytesOut)
			return nil
		}
	}
}
