package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/handlers"
	"github.com/ednesic/coursemanagement/metrics"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var err error
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	if os.Getenv("ENV") == "prod" {
		e.Logger.SetLevel(log.INFO)
	}

	cache.GetInstance().Initialize(map[string]string{"server1": os.Getenv("REDIS_HOST")})
	err = storage.GetInstance().Initialize(
		context.Background(),
		os.Getenv("DB_HOST"),
		os.Getenv("DB"),
	)

	if err != nil {
		e.Logger.Fatal("Could not resolve Data access layer: ", err)
	}

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(metrics.NewMetric())
	e.Use(middleware.Logger())

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/courses/:name", handlers.GetCourse)
	e.GET("/courses", handlers.GetCourses)
	e.PUT("/courses", handlers.PutCourse)
	e.POST("/courses", handlers.SetCourse)
	e.DELETE("/courses/:name", handlers.DelCourse)

	go func() {
		if err := e.Start(":9090"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	cache.GetInstance().Disconnect()
	storage.GetInstance().Disconnect()
}
