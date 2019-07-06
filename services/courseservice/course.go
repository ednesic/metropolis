package courseservice

import (
	"context"
	"github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"sync"
	"time"
)

const coll = "course"

var instance CourseService
var once sync.Once

type CourseService interface {
	Create(types.Course) error
	Update(types.Course) error
	FindAll() ([]types.Course, error)
	Delete(string) error
	FindOne(string) (types.Course, error)
}

type courseImpl struct{}

func GetInstance() CourseService {
	once.Do(func() {
		if instance == nil {
			instance = &courseImpl{}
		}
	})
	return instance
}

func (s *courseImpl) FindOne(name string) (c types.Course, err error) {
	var mgoErr error
	ctx := context.Background()
	if err := cache.GetInstance().Get(coll+name, &c); err != nil {
		if mgoErr = storage.GetInstance().FindOne(ctx, coll, map[string]interface{}{"name": name}, &c); mgoErr == nil {
			return c, cache.GetInstance().Set(coll+name, c, time.Minute)
		}
	}
	return c, mgoErr
}

func (s *courseImpl) Create(course types.Course) error {
	ctx := context.Background()
	err := storage.GetInstance().Insert(ctx, coll, course)
	if err == nil {
		return cache.GetInstance().Set(coll+course.Name, course, time.Minute)
	}
	return err
}

func (s *courseImpl) Update(course types.Course) error {
	ctx := context.Background()
	err := storage.
		GetInstance().
		Update(ctx, coll, map[string]interface{}{"name": course.Name}, map[string]interface{}{"$set": &course})
	if err == nil {
		return cache.GetInstance().Delete(coll + course.Name)
	}
	return err
}

func (s *courseImpl) FindAll() ([]types.Course, error) {
	var mgoErr error
	var cs []types.Course
	ctx := context.Background()
	suffixKey := "all"

	if cacheErr := cache.GetInstance().Get(coll+suffixKey, &cs); cacheErr != nil {
		if mgoErr = storage.GetInstance().Find(ctx, coll, map[string]interface{}{}, &cs); mgoErr == nil {
			return cs, cache.GetInstance().Set(coll+suffixKey, cs, time.Minute)
		}
	}
	return cs, mgoErr
}

func (s *courseImpl) Delete(name string) error {
	ctx := context.Background()
	err := storage.GetInstance().Remove(ctx, coll, map[string]interface{}{"name": name})
	if err == nil {
		return cache.GetInstance().Delete(coll + name)
	}
	return err
}
