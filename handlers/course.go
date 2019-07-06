package handlers

import (
	"net/http"

	"github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/services/courseservice"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/labstack/echo/v4"
)

func GetCourse(c echo.Context) error {
	name := c.Param("name")
	cr, err := courseservice.GetInstance().FindOne(name)
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

func GetCourses(c echo.Context) error {
	cs, err := courseservice.GetInstance().FindAll()

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

func SetCourse(c echo.Context) error {
	var cr types.Course

	if err := c.Bind(&cr); err != nil {
		_ = c.NoContent(http.StatusBadRequest)
		return err
	}

	err := courseservice.GetInstance().Create(cr)
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

func PutCourse(c echo.Context) error {
	var cr types.Course

	if err := c.Bind(&cr); err != nil {
		_ = c.NoContent(http.StatusBadRequest)
		return err
	}

	err := courseservice.GetInstance().Update(cr)
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

func DelCourse(c echo.Context) error {
	name := c.Param("name")
	httpStatus := http.StatusOK

	err := courseservice.GetInstance().Delete(name)
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
