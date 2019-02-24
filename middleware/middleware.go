package middleware

import (
	"github.com/ednesic/coursemanagement/context"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Handler func(http.ResponseWriter, *http.Request) error

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context.AddRequestError(r, fn(w, r))

}

type loggerMiddleware struct {
	logger *zap.Logger
}

func NewLoggerMiddleware() *loggerMiddleware {
	logger, _:= zap.NewProduction()
	return &loggerMiddleware{
		logger: logger,
	}
}

func (l *loggerMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	err := context.GetRequestError(r)
	duration := time.Since(start)
	statusCode := rw.(negroni.ResponseWriter).Status()
	if statusCode == 0 {
		statusCode = 200
	}
	nowFormatted := time.Now().Format(time.RFC3339Nano)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	l.logger.Info("Request",
		zap.String("now", nowFormatted),
		zap.String("scheme", scheme),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Int("statuscode", statusCode),
		zap.String("useragent", r.UserAgent()),
		zap.Error(err),
		zap.Float64("duration", float64(duration)/float64(time.Millisecond)))
}


func ContextClearerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer context.Clear(r)
	next(w, r)
}