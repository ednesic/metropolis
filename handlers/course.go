package handlers

import (
	"github.com/ednesic/coursemanagement/cache"
	internalMiddleware "github.com/ednesic/coursemanagement/middleware"
	"github.com/ednesic/coursemanagement/servicemanager"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetCourse(c echo.Context) error {
	name := c.Param("name")
	course, err := servicemanager.CourseService.FindOne(name)
	httpStatus := http.StatusOK

	if serr, ok := err.(*cache.RedisErr); ok {
		c.Set(internalMiddleware.RedisContext, serr)
		err = nil
	}
	if err == nil {
		return c.JSON(httpStatus, course)
	}
	httpStatus = http.StatusInternalServerError
	if err == storage.ErrNotFound {
		httpStatus = http.StatusNotFound
	}
	_ = c.NoContent(httpStatus)
	return err
}

func GetCourses(c echo.Context) error {
	courses, err := servicemanager.CourseService.FindAll()

	if courses == nil {
		courses = []types.Course{}
	}
	if err == nil {
		return c.JSON(http.StatusOK, courses)
	}
	_ = c.NoContent(http.StatusInternalServerError)
	return err
}

func SetCourse(c echo.Context) error {
	var course types.Course

	if err := c.Bind(&course); err != nil {
		return err
	}

	err := servicemanager.CourseService.Create(course)
	if err == nil {
		return c.JSON(http.StatusCreated, course)
	}

	_ = c.NoContent(http.StatusInternalServerError)
	return err
}

func PutCourse(c echo.Context) error {
	var course types.Course

	if err := c.Bind(&course); err != nil {
		return err
	}

	err := servicemanager.CourseService.Update(course)
	if err == nil {
		return c.JSON(http.StatusCreated, course)
	}

	_ = c.NoContent(http.StatusInternalServerError)
	return err
	}

func DelCourse(c echo.Context) error {
	name := c.Param("name")
	httpStatus := http.StatusOK

	err:= servicemanager.CourseService.Delete(types.Course{Name: name})
	if err != nil {
		httpStatus = http.StatusInternalServerError
	}
	if err == storage.ErrNotFound {
		httpStatus = http.StatusNotFound
	}
	_ = c.NoContent(httpStatus)
	return err
}