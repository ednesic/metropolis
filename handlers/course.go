package handlers

import (
	"github.com/ednesic/coursemanagement/servivemanager"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetCourse(c echo.Context) error {
	name := c.Param("name")
	course, err := servivemanager.CourseService.FindOne(name)
	httpStatus := http.StatusOK

	if err == nil {
		return c.JSON(httpStatus, course)
	}

	httpStatus = http.StatusInternalServerError
	if err == storage.ErrNotFound {
		httpStatus = http.StatusNotFound
	}
	if err := c.JSON(httpStatus, storage.ErrorReponse{ Message: "Failed to get course" }); err != nil {
		return err
	}
	return err
}

func GetCourses(c echo.Context) error {
	courses, err := servivemanager.CourseService.FindAll()

	if courses == nil {
		courses = []types.Course{}
	}
	if err == nil {
		return c.JSON(http.StatusOK, courses)
	}
	if err := c.JSON(http.StatusInternalServerError,
		storage.ErrorReponse{ Message: "Failed to get courses" }); err != nil {
		return err
	}
	return err
}

func SetCourse(c echo.Context) error {
	var course types.Course
	errMsg := storage.ErrorReponse{ Message: "Failed to set course" }

	err := c.Bind(&course)
	if err != nil {
		if err := c.JSON(http.StatusBadRequest, errMsg); err != nil {
			return err
		}
		return err
	}

	err = servivemanager.CourseService.Create(course)
	if err == nil {
		return c.JSON(http.StatusCreated, course)
	}

	if err := c.JSON(http.StatusInternalServerError, errMsg); err != nil {
		return err
	}
	return err
}

func PutCourse(c echo.Context) error {
	var course types.Course
	errMsg := storage.ErrorReponse{ Message: "Failed to put course" }

	err := c.Bind(&course)
	if err != nil {
		if err := c.JSON(http.StatusBadRequest,errMsg); err != nil {
			return  err
		}
		return err
	}

	err = servivemanager.CourseService.Update(course)
	if err == nil {
		return c.JSON(http.StatusCreated, course)
	}

	if err := c.JSON(http.StatusInternalServerError, errMsg); err != nil {
		return err
	}
	return err
	}

func DelCourse(c echo.Context) error {
	name := c.Param("name")
	httpStatus := http.StatusOK

	err:= servivemanager.CourseService.Delete(types.Course{Name:name})
	if err != nil {
		httpStatus = http.StatusInternalServerError
	}
	if err == storage.ErrNotFound {
		httpStatus = http.StatusNotFound
	}
	if err := c.NoContent(httpStatus); err != nil {
		return err
	}
	return err
}