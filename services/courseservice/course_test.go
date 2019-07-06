package courseservice

import (
	"context"
	"errors"
	"testing"

	redis "github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/go-redis/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCourseFindOne_FindsCourseCached(t *testing.T) {
	redisMock := &redis.Mock{}
	testName := "test01"
	redisCourseMock := types.Course{Name: testName}
	redisMock.Initialize(map[string]string{})

	redisMock.On("Get", coll+testName, mock.Anything).Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*types.Course)
			*arg = redisCourseMock
		}).Once()
	courseService := courseImpl{}

	c, err := courseService.FindOne(testName)
	assert.Nil(t, err)
	assert.Equal(t, c, redisCourseMock)

	redisMock.AssertExpectations(t)
}

func TestCourseFindOne_DoNotFindCourseCached(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	testName := "test01"
	mongoCourseMock := types.Course{Name: testName}
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	redisMock.On("Get", coll+testName, mock.Anything).Return(cache.ErrCacheMiss).Once()
	mongoMock.On("FindOne", mock.Anything, coll, mock.Anything, mock.AnythingOfType("*types.Course")).
		Run(func(args mock.Arguments) {
			arg := args.Get(3).(*types.Course)
			*arg = mongoCourseMock
		}).Return(nil).Once()
	redisMock.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	courseService := courseImpl{}

	c, err := courseService.FindOne(testName)
	assert.Nil(t, err)
	assert.Equal(t, c, mongoCourseMock)

	mongoMock.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}

func TestCourseCreate_ErrOnInsert(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	testCourse := types.Course{Name: "test02"}
	errMock := errors.New("insert err")
	_ = mongoMock.Initialize(context.Background(), "", "")

	mongoMock.On("Insert", mock.Anything, coll, mock.AnythingOfType("types.Course")).
		Return(errMock).Once()

	courseService := courseImpl{}

	err := courseService.Create(testCourse)
	assert.Equal(t, err, errMock)
	mongoMock.AssertExpectations(t)
}

func TestCourseCreate_ErrOnCache(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	testCourse := types.Course{Name: "test02"}
	errMock := errors.New("insert err")
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	mongoMock.On("Insert", mock.Anything, coll, mock.AnythingOfType("types.Course")).
		Return(nil).Once()
	redisMock.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(errMock).Once()

	courseService := courseImpl{}

	err := courseService.Create(testCourse)
	assert.Equal(t, errMock, err)
	mongoMock.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}

func TestCourseCreate_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	testCourse := types.Course{Name: "test02"}
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	mongoMock.On("Insert", mock.Anything, coll, mock.AnythingOfType("types.Course")).
		Return(nil).Once()
	redisMock.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	courseService := courseImpl{}

	err := courseService.Create(testCourse)
	assert.Nil(t, err)
	mongoMock.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}

func TestCourseUpdate_ErrUpdate(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	testCourse := types.Course{Name: "test02"}
	errMock := errors.New("err update")
	_ = mongoMock.Initialize(context.Background(), "", "")

	mongoMock.On("Update", mock.Anything, coll, mock.Anything, mock.Anything).
		Return(errMock).Once()

	courseService := courseImpl{}

	err := courseService.Update(testCourse)
	assert.Equal(t, err, errMock)
	mongoMock.AssertExpectations(t)
}

func TestCourseUpdate_ErrCache(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	testCourse := types.Course{Name: "test02"}
	errMock := errors.New("err update")
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	mongoMock.On("Update", mock.Anything, coll, mock.Anything, mock.Anything).
		Return(nil).Once()
	redisMock.On("Delete", mock.Anything).Return(errMock).Once()

	courseService := courseImpl{}

	err := courseService.Update(testCourse)
	assert.Equal(t, err, errMock)
	mongoMock.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}

func TestCourseUpdate_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	testCourse := types.Course{Name: "test02"}
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	mongoMock.On("Update", mock.Anything, coll, mock.Anything, mock.Anything).
		Return(nil).Once()
	redisMock.On("Delete", mock.Anything).Return(nil).Once()

	courseService := courseImpl{}

	err := courseService.Update(testCourse)
	assert.Nil(t, err)
	mongoMock.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}

func TestCourseFindAll_SuccessGetCache(t *testing.T) {
	redisMock := &redis.Mock{}
	suffix := "all"
	redisCourseMock := []types.Course{{Name: "test03"}, {Name: "test04"}}
	redisMock.Initialize(map[string]string{})

	redisMock.On("Get", coll+suffix, mock.Anything).Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*[]types.Course)
			*arg = redisCourseMock
		}).Once()

	courseService := courseImpl{}

	c, err := courseService.FindAll()
	assert.Nil(t, err)
	assert.Equal(t, c, redisCourseMock)

	redisMock.AssertExpectations(t)
}

func TestCourseFindAll_ErrGet(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	suffix := "all"
	errMock := errors.New("err find")
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	redisMock.On("Get", coll+suffix, mock.Anything).Return(cache.ErrCacheMiss).Once()
	mongoMock.On("Find", mock.Anything, coll, mock.Anything, mock.AnythingOfType("*[]types.Course")).Return(errMock).Once()

	courseService := courseImpl{}

	c, err := courseService.FindAll()
	assert.Equal(t, err, errMock)
	assert.Len(t, c, 0)

	redisMock.AssertExpectations(t)
	mongoMock.AssertExpectations(t)
}

func TestCourseFindAll_ErrSetCache(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	suffix := "all"
	errMock := errors.New("err set cache")
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	redisMock.On("Get", coll+suffix, mock.Anything).Return(cache.ErrCacheMiss).Once()
	mongoMock.On("Find", mock.Anything, coll, mock.Anything, mock.AnythingOfType("*[]types.Course")).Return(nil).Once()
	redisMock.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(errMock).Once()

	courseService := courseImpl{}

	c, err := courseService.FindAll()
	assert.Equal(t, errMock, err)
	assert.Len(t, c, 0)

	redisMock.AssertExpectations(t)
	mongoMock.AssertExpectations(t)
}

func TestCourseFindAll_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	suffix := "all"
	mongoCourseMock := []types.Course{{Name: "test03"}, {Name: "test04"}}
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	redisMock.On("Get", coll+suffix, mock.Anything).Return(cache.ErrCacheMiss).Once()
	mongoMock.On("Find", mock.Anything, coll, mock.Anything, mock.AnythingOfType("*[]types.Course")).Run(func(args mock.Arguments) {
		arg := args.Get(3).(*[]types.Course)
		*arg = mongoCourseMock
	}).Return(nil).Once()
	redisMock.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	courseService := courseImpl{}

	c, err := courseService.FindAll()
	assert.Nil(t, err)
	assert.Equal(t, c, mongoCourseMock)

	redisMock.AssertExpectations(t)
	mongoMock.AssertExpectations(t)
}

func TestCourseDelete_ErrDelete(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	errMock := errors.New("err delete")
	mongoMock.On("Remove", mock.Anything, coll, mock.Anything).Return(errMock).Once()
	testCourse := "test02"
	_ = mongoMock.Initialize(context.Background(), "", "")

	courseService := courseImpl{}

	err := courseService.Delete(testCourse)
	assert.Equal(t, err, errMock)

	mongoMock.AssertExpectations(t)
}

func TestCourseDelete_ErrCache(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	errMock := errors.New("err delete")
	testCourse := "test02"
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	redisMock.On("Delete", coll+testCourse).Return(errMock).Once()
	mongoMock.On("Remove", mock.Anything, coll, mock.Anything).Return(nil).Once()

	courseService := courseImpl{}

	err := courseService.Delete(testCourse)
	assert.Equal(t, err, errMock)

	redisMock.AssertExpectations(t)
	mongoMock.AssertExpectations(t)
}

func TestCourseDelete_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	redisMock := &redis.Mock{}
	testCourse := "test02"
	redisMock.Initialize(map[string]string{})
	_ = mongoMock.Initialize(context.Background(), "", "")

	redisMock.On("Delete", coll+testCourse).Return(nil).Once()
	mongoMock.On("Remove", mock.Anything, coll, mock.Anything).Return(nil).Once()

	courseService := courseImpl{}

	err := courseService.Delete(testCourse)
	assert.Nil(t, err)

	redisMock.AssertExpectations(t)
	mongoMock.AssertExpectations(t)
}
