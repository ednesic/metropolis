package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/services/courseservice"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2"
)

func TestGetCourse(t *testing.T) {
	type fields struct {
		name    string
		mockErr error
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
		{"Status ok", fields{name: "nameTest"}, wants{course: types.Course{Name: "nameTest", Price: 10, Picture: "pic.png", PreviewURLVideo: "http://video"}, statusCode: http.StatusOK, err: nil}},
		{"Status ok but redis err", fields{mockErr: &cache.RedisErr{}, name: "nameTest"}, wants{course: types.Course{Name: "nameTest", Price: 10, Picture: "pic.png", PreviewURLVideo: "http://video"}, statusCode: http.StatusOK, err: nil}},
		{"Status notFound", fields{mockErr: storage.ErrNotFound, name: "nameNotFound"}, wants{course: types.Course{}, statusCode: http.StatusNotFound, err: storage.ErrNotFound}},
		{"Status internal server error", fields{mockErr: mgo.ErrCursor, name: "nameInternal"}, wants{course: types.Course{}, statusCode: http.StatusInternalServerError, err: mgo.ErrCursor}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var courseServiceMngr = &courseservice.Mock{}
			courseServiceMngr.On("FindOne", tt.fields.name).Return(tt.want.course, tt.fields.mockErr).Once()
			courseServiceMngr.InitMock()

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

func BenchmarkGetCourse(b *testing.B) {
	var courseServiceMngr = &courseservice.Mock{}
	courseServiceMngr.On("FindOne", mock.Anything).Return(types.Course{Name: "bench"}, nil)
	courseServiceMngr.InitMock()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/course", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("Bench")

	for i := 0; i < b.N; i++ {
		_ = GetCourse(c)
	}
}

func TestGetCourses(t *testing.T) {
	type wants struct {
		courses    []types.Course
		err        error
		statusCode int
	}
	type mocks struct {
		courses []types.Course
		err     error
	}
	tests := []struct {
		name string
		mock mocks
		want wants
	}{
		{"Status ok", mocks{courses: []types.Course{{}, {}}, err: nil}, wants{courses: []types.Course{{}, {}}, err: nil, statusCode: http.StatusOK}},
		{"Status ok but redis err", mocks{courses: []types.Course{{}, {}}, err: &cache.RedisErr{}}, wants{courses: []types.Course{{}, {}}, err: nil, statusCode: http.StatusOK}},
		{"Status ok(empty)", mocks{courses: nil, err: nil}, wants{courses: []types.Course{}, err: nil, statusCode: http.StatusOK}},
		{"Status internal server error", mocks{courses: nil, err: mgo.ErrCursor}, wants{courses: []types.Course{}, err: mgo.ErrCursor, statusCode: http.StatusInternalServerError}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var courseServiceMngr = &courseservice.Mock{}
			courseServiceMngr.On("FindAll").Return(tt.mock.courses, tt.mock.err).Once()
			courseServiceMngr.InitMock()

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

func BenchmarkGetCourses(b *testing.B) {
	var courseServiceMngr = &courseservice.Mock{}
	courseServiceMngr.On("FindAll").Return([]types.Course{}, nil)
	courseServiceMngr.InitMock()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/course", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for i := 0; i < b.N; i++ {
		_ = GetCourses(c)
	}
}

func TestSetCourse(t *testing.T) {
	type wants struct {
		err        error
		statusCode int
	}
	type mocks struct {
		mongoMockTimes int
		course         types.Course
		err            error
	}
	type fields struct { //passar interface para bad request
		body interface{}
	}
	tests := []struct {
		name  string
		want  wants
		mock  mocks
		field fields
	}{
		{"Status ok", wants{statusCode: http.StatusOK}, mocks{mongoMockTimes: 1, course: types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}, fields{types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}},
		{"Status ok but redis err", wants{statusCode: http.StatusOK}, mocks{err: &cache.RedisErr{}, mongoMockTimes: 1, course: types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}, fields{types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}},
		{"Status bad request", wants{statusCode: http.StatusBadRequest, err: &echo.HTTPError{}}, mocks{}, fields{"{err}"}},
		{"Status internal server error", wants{statusCode: http.StatusInternalServerError, err: errors.New("")}, mocks{mongoMockTimes: 1, err: mgo.ErrCursor}, fields{types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var courseServiceMngr = &courseservice.Mock{}
			courseServiceMngr.On("Create", tt.field.body).Return(tt.mock.err).Maybe().Times(tt.mock.mongoMockTimes)
			courseServiceMngr.InitMock()

			out, err := json.Marshal(tt.field.body)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/course", strings.NewReader(string(out)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = SetCourse(c)
			assert.IsType(t, err, tt.want.err)
			er, ok := err.(*echo.HTTPError)
			if ok {
				assert.Equal(t, tt.want.statusCode, er.Code)
			} else {
				assert.Equal(t, tt.want.statusCode, rec.Code)
			}
			courseServiceMngr.AssertExpectations(t)
		})
	}
}

func BenchmarkSetCourse(b *testing.B) {
	var courseServiceMngr = &courseservice.Mock{}
	courseServiceMngr.On("Create", mock.Anything).Return(nil)
	courseServiceMngr.InitMock()

	out, _ := json.Marshal(types.Course{Name: "BEnch1", Price: 10, Picture: "bench", PreviewURLVideo: "bench"})
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/course", strings.NewReader(string(out)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for i := 0; i < b.N; i++ {
		_ = SetCourse(c)
	}
}

func TestPutCourse(t *testing.T) {
	type wants struct {
		err        error
		statusCode int
	}
	type mocks struct {
		mongoMockTimes int
		course         types.Course
		err            error
	}
	type fields struct {
		body interface{}
	}
	tests := []struct {
		name  string
		want  wants
		mock  mocks
		field fields
	}{
		{"Status ok", wants{statusCode: http.StatusCreated}, mocks{mongoMockTimes: 1, course: types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}, fields{types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}},
		{"Status ok but redis err", wants{statusCode: http.StatusCreated}, mocks{err: &cache.RedisErr{}, mongoMockTimes: 1, course: types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}, fields{types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}},
		{"Status bad request", wants{statusCode: http.StatusBadRequest, err: &echo.HTTPError{}}, mocks{}, fields{"{err}"}},
		{"Status internal server error", wants{statusCode: http.StatusInternalServerError, err: errors.New("")}, mocks{mongoMockTimes: 1, err: mgo.ErrCursor}, fields{types.Course{Name: "Test123", Price: 10, Picture: "test.png", PreviewURLVideo: "http://video"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var courseServiceMngr = &courseservice.Mock{}
			courseServiceMngr.On("Update", tt.field.body).Return(tt.mock.err).Maybe().Times(tt.mock.mongoMockTimes)
			courseServiceMngr.InitMock()

			out, err := json.Marshal(tt.field.body)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/course", strings.NewReader(string(out)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = PutCourse(c)
			assert.IsType(t, err, tt.want.err)
			er, ok := err.(*echo.HTTPError)
			if ok {
				assert.Equal(t, tt.want.statusCode, er.Code)
			} else {
				assert.Equal(t, tt.want.statusCode, rec.Code)
			}
			courseServiceMngr.AssertExpectations(t)
		})
	}
}

func BenchmarkPutCourse(b *testing.B) {
	var courseServiceMngr = &courseservice.Mock{}
	courseServiceMngr.On("Update", mock.Anything).Return(nil)
	courseServiceMngr.InitMock()

	out, _ := json.Marshal(types.Course{Name: "BEnch1", Price: 10, Picture: "bench", PreviewURLVideo: "bench"})
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/course", strings.NewReader(string(out)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for i := 0; i < b.N; i++ {
		_ = PutCourse(c)
	}
}

func TestDelCourse(t *testing.T) {
	type fields struct {
		name string
		err  error
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
		{"Status ok", fields{name: "nameTest"}, wants{course: types.Course{Name: "nameTest", Price: 10, Picture: "pic.png", PreviewURLVideo: "http://video"}, statusCode: http.StatusOK, err: nil}},
		{"Status ok but redis err", fields{err: &cache.RedisErr{}, name: "nameTest"}, wants{course: types.Course{Name: "nameTest", Price: 10, Picture: "pic.png", PreviewURLVideo: "http://video"}, statusCode: http.StatusOK, err: nil}},
		{"Status notFound", fields{err: storage.ErrNotFound, name: "nameNotFound"}, wants{course: types.Course{}, statusCode: http.StatusNotFound, err: storage.ErrNotFound}},
		{"Status internal server error", fields{err: mgo.ErrCursor, name: "nameInternal"}, wants{course: types.Course{}, statusCode: http.StatusInternalServerError, err: mgo.ErrCursor}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var courseServiceMngr = &courseservice.Mock{}
			courseServiceMngr.On("Delete", tt.fields.name).Return(tt.fields.err).Once()
			courseServiceMngr.InitMock()

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

func BenchmarkDelCourse(b *testing.B) {
	var courseServiceMngr = &courseservice.Mock{}
	courseServiceMngr.On("Delete", mock.Anything).Return(nil)
	courseServiceMngr.InitMock()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/course", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("Bench")

	for i := 0; i < b.N; i++ {
		_ = DelCourse(c)
	}
}
