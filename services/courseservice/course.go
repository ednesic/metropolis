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

type impl struct {}

func GetInstance() CourseService {
	once.Do(func() {
		if instance == nil {
			instance = &impl{}
		}
	})
	return instance
}

func (s *impl) FindOne(name string) (c types.Course, err error) {
	var mgoErr error
	ctx := context.Background()
	if err := cache.GetInstance().Get(coll+name, &c); err != nil {
		if mgoErr = storage.GetInstance().FindOne(ctx, coll, map[string]interface{}{"name": name}, &c); mgoErr == nil {
			return c, cache.GetInstance().Set(coll + name, c, time.Minute)
		}
	}
	return c, mgoErr
}

func (s *impl) Create(course types.Course) error {
	ctx := context.Background()
	err := storage.GetInstance().Insert(ctx, coll, course)
	if err == nil {
		return cache.GetInstance().Set(coll + course.Name, course, time.Minute)
	}
	return err
}

func (s *impl) Update(course types.Course) error {
	ctx := context.Background()
	err := storage.GetInstance().Update(ctx, coll, map[string]interface{}{"name": course.Name}, &course)
	if err == nil {
		return cache.GetInstance().Delete(coll + course.Name)
	}
	return err
}

func (s *impl) FindAll() ([]types.Course, error) {
	cs := []types.Course{}
	ctx := context.Background()
	var mgoErr error
	suffixKey := "all"

	if cacheErr := cache.GetInstance().Get(coll + suffixKey, &cs); cacheErr != nil {
		if mgoErr = storage.GetInstance().Find(ctx, coll, map[string]interface{}{}, &cs); mgoErr == nil {
			return cs, cache.GetInstance().Set(coll + suffixKey, cs, time.Minute)
		}
	}
	return cs, mgoErr
}

func (s *impl) Delete(name string) error {
	ctx := context.Background()
	c, err := storage.GetInstance().StartSession()
	if err != nil {
		return err
	}
	err = c.WithTransaction(func(c context.Context) error {
		err = storage.GetInstance().Insert(c, coll, &types.Course{Name: "test1533"})
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return err
	}
	err = c.WithTransaction(func(c context.Context) error {
		err = storage.GetInstance().Insert(c, coll, &types.Course{Name: "test4355"})
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return err
	}
	err = c.Abort()
	if err != nil {
		return err
	}
	err = storage.GetInstance().Remove(ctx, coll, map[string]interface{}{"name": name})
	if err == nil {
		return cache.GetInstance().Delete(coll + name)
	}
	return err
}