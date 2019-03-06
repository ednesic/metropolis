package main

import (
	"github.com/ednesic/coursemanagement/metrics"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"

	"github.com/ednesic/coursemanagement/handlers"
	"github.com/ednesic/coursemanagement/servicemanager"
	"github.com/ednesic/coursemanagement/services"
)

func main() {
	var err error
	e := echo.New()

	servicemanager.CourseService, err = services.NewCourseService(
		os.Getenv("COURSE_DB_HOST"),
		os.Getenv("COURSE_DB"),
		map[string]string{ "server1": os.Getenv("COURSE_REDIS_HOST")},
		)

	if err != nil {
		e.Logger.Fatal("Could not resolve course service", err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(metrics.NewMetric())

	//e.Server.ReadTimeout = time.Duration(1 * time.Second)
	//e.Server.WriteTimeout= time.Duration(1 * time.Second)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/courses/:name", handlers.GetCourse)
	e.GET("/courses", handlers.GetCourses)
	e.PUT("/courses", handlers.PutCourse)
	e.POST("/courses", handlers.SetCourse)
	e.DELETE("/courses/:name", handlers.DelCourse)

	e.Logger.Fatal(e.Start(":9090"))
}