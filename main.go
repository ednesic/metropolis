package main

import (
	"context"
	"github.com/ednesic/coursemanagement/metrics"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"os/signal"
	"time"

	"github.com/ednesic/coursemanagement/handlers"
	"github.com/ednesic/coursemanagement/servicemanager"
	"github.com/ednesic/coursemanagement/services"
)

func main() {
	var err error
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	if os.Getenv("ENV") == "prod" {
		e.Logger.SetLevel(log.INFO)
	}

	servicemanager.CourseService, err = services.NewCourseService(
		os.Getenv("COURSE_DB_HOST"),
		os.Getenv("COURSE_DB"),
		map[string]string{ "server1": os.Getenv("COURSE_REDIS_HOST")},
		)

	if err != nil {
		e.Logger.Fatal("Could not resolve course service: ", err)
	}

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(metrics.NewMetric())
	e.Use(middleware.Logger())

	//e.Server.ReadTimeout = time.Duration(1 * time.Second)
	//e.Server.WriteTimeout= time.Duration(1 * time.Second)

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
	servicemanager.CourseService.Shutdown()
}