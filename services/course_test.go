package services

import (
	redis "github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/go-redis/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCourseFindOne_FindsCourseCached(t *testing.T) {
	redisMock := &redis.RedisMock{}
	testName := "test01"
	redisCourseMock := types.Course{Name: testName}

	redisMock.On("Get", coll+testName, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*types.Course)
		*arg = redisCourseMock
	}).Once()
	redisMock.On("Set", mock.Anything).Return(nil).Once()

	courseService := CourseServiceImpl{
		cache: redisMock,
	}

	c, err := courseService.FindOne(testName)
	assert.Nil(t, err)
	assert.Equal(t, c, redisCourseMock)

	redisMock.AssertExpectations(t)

}

func TestCourseFindOne_DoNotFindCourseCached(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.RedisMock{}
	testName := "test01"
	mongoCourseMock := types.Course{Name: testName}

	redisMock.On("Get", coll+testName, mock.Anything).Return(cache.ErrCacheMiss).Once()
	mongoMock.On("FindOne", coll, mock.AnythingOfType("bson.M"), mock.AnythingOfType("*types.Course")).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*types.Course)
		*arg = mongoCourseMock
	}).Return(nil).Once()
	redisMock.On("Set", mock.Anything).Return(nil).Once()

	courseService := CourseServiceImpl{
		cache: redisMock,
		dal:   mongoMock,
	}

	c, err := courseService.FindOne(testName)
	assert.Nil(t, err)
	assert.Equal(t, c, mongoCourseMock)

	redisMock.AssertExpectations(t)

}

func TestCourseCreate(t *testing.T) {

}

func TestCourseUpdate(t *testing.T) {

}

func TestCourseFindAll(t *testing.T) {

}

func TestCourseDelete(t *testing.T) {

}
