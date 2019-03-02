package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ednesic/coursemanagement/servicemanager"
	"github.com/ednesic/coursemanagement/services"
	"github.com/ednesic/coursemanagement/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetCourse(t *testing.T) {
	type fields struct {
		name string
	}
	type wants struct {
		course     types.Course
		statusCode int
		err        error
	}
	tests := []struct {
		name   string
		fields fields
		want   wants
	}{
		{"Status ok", fields{"nameTest"}, wants{course: types.Course{Name: "nameTest", Price: 10, Picture: "pic.png", PreviewUrlVideo: "http://video"}, statusCode: http.StatusOK, err: nil}},
		{"Status notFound", fields{"nameNotFound"}, wants{course: types.Course{}, statusCode: http.StatusNotFound, err: mgo.ErrNotFound}},
		{"Status internal server error", fields{"nameInternal"}, wants{course: types.Course{}, statusCode: http.StatusInternalServerError, err: mgo.ErrCursor}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var courseServiceMngr = &services.CourseServiceMock{}
			courseServiceMngr.On("FindOne", tt.fields.name).Return(tt.want.course, tt.want.err).Once()
			servicemanager.CourseService = courseServiceMngr

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/course", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("name")
			c.SetParamValues(tt.fields.name)

			out, err := json.Marshal(tt.want.course)
			assert.NoError(t, err)
			assert.Equal(t, GetCourse(c), tt.want.err)
			assert.Equal(t, tt.want.statusCode, rec.Code)
			if tt.want.err == nil {
				assert.Equal(t, fmt.Sprintf("%s\n", out), rec.Body.String())
			}
			courseServiceMngr.AssertExpectations(t)
		})
	}
}

func Test_GetCourses(t *testing.T) {
	type wants struct {
		courses    []types.Course
		err        error
		statusCode int
	}
	type mocks struct {
		courses []types.Course
		err error
	}
	tests := []struct {
		name string
		mock mocks
		want wants
	}{
		{"Status ok", mocks{courses: []types.Course{{}, {}}, err: nil},wants{courses: []types.Course{{}, {}}, err: nil, statusCode: http.StatusOK}},
		{"Status ok(empty)", mocks{courses: nil, err: nil},wants{courses: []types.Course{}, err: nil, statusCode: http.StatusOK}},
		{"Status internal server error", mocks{courses: nil, err: mgo.ErrCursor},wants{courses: []types.Course{}, err: mgo.ErrCursor, statusCode: http.StatusInternalServerError}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var courseServiceMngr = &services.CourseServiceMock{}
			courseServiceMngr.On("FindAll").Return(tt.mock.courses, tt.mock.err).Once()
			servicemanager.CourseService = courseServiceMngr

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/course", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			out, err := json.Marshal(tt.want.courses)
			assert.NoError(t, err)
			assert.Equal(t, GetCourses(c), tt.want.err)
			assert.Equal(t, tt.want.statusCode, rec.Code)
			if tt.want.err == nil {
				assert.Equal(t, fmt.Sprintf("%s\n", out), rec.Body.String())
			}
			courseServiceMngr.AssertExpectations(t)
		})
	}
}

func TestSetCourse(t *testing.T) {

}

func Test_PutCourse(t *testing.T) {

}

func Test_DelCourse(t *testing.T) {
	type fields struct {
		name string
	}
	type wants struct {
		course     types.Course
		statusCode int
		err        error
	}
	tests := []struct {
		name   string
		fields fields
		want   wants
	}{
		{"Status ok", fields{"nameTest"}, wants{course: types.Course{Name: "nameTest", Price: 10, Picture: "pic.png", PreviewUrlVideo: "http://video"}, statusCode: http.StatusOK, err: nil}},
		{"Status notFound", fields{"nameNotFound"}, wants{course: types.Course{}, statusCode: http.StatusNotFound, err: mgo.ErrNotFound}},
		{"Status internal server error", fields{"nameInternal"}, wants{course: types.Course{}, statusCode: http.StatusInternalServerError, err: mgo.ErrCursor}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var courseServiceMngr = &services.CourseServiceMock{}
			courseServiceMngr.On("Delete", types.Course{Name: tt.fields.name}).Return(tt.want.err).Once()
			servicemanager.CourseService = courseServiceMngr

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/course", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("name")
			c.SetParamValues(tt.fields.name)

			assert.Equal(t, DelCourse(c), tt.want.err)
			assert.Equal(t, tt.want.statusCode, rec.Code)
			courseServiceMngr.AssertExpectations(t)
		})
	}
}
