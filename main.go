package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ednesic/coursemanagement/handlers"
	"github.com/ednesic/coursemanagement/middleware"
	"github.com/ednesic/coursemanagement/services"
	"github.com/ednesic/coursemanagement/servivemanager"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

func main() {
	var err error
	servivemanager.CourseService, err = services.NewCourseService(os.Getenv("COURSE_DB_HOST"), os.Getenv("COURSE_DB"))

	if err != nil {
		log.Fatal("Could not resolve course service", zap.Error(err))
	}

	srv := &http.Server{
		Addr:    ":" + "9090",
		Handler: initRoutes(),
	}

	log.Fatal("Shutting down", srv.ListenAndServe())
}

func initRoutes() *negroni.Negroni {
	r := mux.NewRouter()
	r.NewRoute().Path("/courses/{name}").Handler(middleware.Handler(handlers.GetCourse)).Methods(http.MethodGet)
	r.NewRoute().Path("/courses/{name}").Handler(middleware.Handler(handlers.DelCourse)).Methods(http.MethodDelete)
	r.NewRoute().Path("/courses").Handler(middleware.Handler(handlers.PutCourse)).Methods(http.MethodPut)
	r.NewRoute().Path("/courses").Handler(middleware.Handler(handlers.GetCourses)).Methods(http.MethodGet)
	r.NewRoute().Path("/courses").Handler(middleware.Handler(handlers.SetCourse)).Methods(http.MethodPost)
	n := negroni.New(
		//negroni.NewRecovery(),
		//negroni.HandlerFunc(middleware.ContextClearerMiddleware),
		middleware.NewLoggerMiddleware(),
		)
	n.UseHandler(r)

	return n
}

