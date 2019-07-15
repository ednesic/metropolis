package handlers

import (
	"net/http"

	"github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/services/courseservice"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/labstack/echo/v4"
)

//GetCourse is a handler to get course passing a query parameter name
func GetCourse(c echo.Context) error {
	name := c.Param("name")
	cr, err := courseservice.GetInstance().FindOne(c.Request().Context(), name)
	httpStatus := http.StatusOK

	if serr, ok := err.(*cache.RedisErr); ok {
		c.Logger().Warn(serr)
		err = nil
	}
	if err == nil {
		return c.JSON(httpStatus, cr)
	}
	httpStatus = http.StatusInternalServerError
	if err == storage.ErrNotFound {
		httpStatus = http.StatusNotFound
	}
	_ = c.NoContent(httpStatus)
	return err
}

//GetCourses is a handler to get all courses
func GetCourses(c echo.Context) error {
	cs, err := courseservice.GetInstance().FindAll(c.Request().Context())

	if serr, ok := err.(*cache.RedisErr); ok {
		c.Logger().Warn(serr)
		err = nil
	}
	if cs == nil {
		cs = []types.Course{}
	}
	if err == nil {
		return c.JSON(http.StatusOK, cs)
	}
	_ = c.NoContent(http.StatusInternalServerError)
	return err
}

//SetCourse is a handler to create a course passing a type.Course in the body
func SetCourse(c echo.Context) error {
	var cr types.Course

	if err := c.Bind(&cr); err != nil {
		_ = c.NoContent(http.StatusBadRequest)
		return err
	}

	err := courseservice.GetInstance().Create(c.Request().Context(), cr)
	if serr, ok := err.(*cache.RedisErr); ok {
		c.Logger().Warn(serr)
		err = nil
	}
	if err == nil {
		return c.JSON(http.StatusOK, cr)
	}

	_ = c.NoContent(http.StatusInternalServerError)
	return err
}

//PutCourse is a handler to update a course passing a type.Course in the body
func PutCourse(c echo.Context) error {
	var cr types.Course

	if err := c.Bind(&cr); err != nil {
		_ = c.NoContent(http.StatusBadRequest)
		return err
	}

	err := courseservice.GetInstance().Update(c.Request().Context(), cr)
	if serr, ok := err.(*cache.RedisErr); ok {
		c.Logger().Warn(serr)
		err = nil
	}
	if err == nil {
		return c.JSON(http.StatusCreated, cr)
	}

	_ = c.NoContent(http.StatusInternalServerError)
	return err
}

//DelCourse is a handler that deletes a course tha has a the query parameter name
func DelCourse(c echo.Context) error {
	name := c.Param("name")
	httpStatus := http.StatusOK

	err := courseservice.GetInstance().Delete(c.Request().Context(), name)
	if serr, ok := err.(*cache.RedisErr); ok {
		c.Logger().Warn(serr)
		err = nil
	}
	if err != nil {
		httpStatus = http.StatusInternalServerError
	}
	if err == storage.ErrNotFound {
		httpStatus = http.StatusNotFound
	}
	_ = c.NoContent(httpStatus)
	return err
}
