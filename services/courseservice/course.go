package courseservice

import (
	"context"
	"sync"
	"time"

	"github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
)

const coll = "course"

var (
	instance CourseService
	once     sync.Once
)

//CourseService is an interface for course service
type CourseService interface {
	Create(context.Context, types.Course) error
	Update(context.Context, types.Course) error
	FindAll(context.Context) ([]types.Course, error)
	Delete(context.Context, string) error
	FindOne(context.Context, string) (types.Course, error)
}

type courseImpl struct{}

//GetInstance to get service instance
func GetInstance() CourseService {
	once.Do(func() {
		if instance == nil {
			instance = &courseImpl{}
		}
	})
	return instance
}

func (s courseImpl) FindOne(ctx context.Context, name string) (types.Course, error) {
	var c types.Course
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := cache.GetInstance().Get(ctx, coll+name, &c)
	if err == nil {
		return c, nil
	}
	mgoErr := storage.GetInstance().FindOne(ctx, coll, map[string]interface{}{"name": name}, &c)
	if mgoErr == nil {
		_ = cache.GetInstance().Set(ctx, coll+name, c, time.Minute)
		return c, nil
	}
	if mgoErr == storage.ErrNotFound {
		return c, nil
	}
	return c, mgoErr
}

func (s courseImpl) Create(ctx context.Context, course types.Course) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := storage.GetInstance().Insert(ctx, coll, course)
	if err == nil {
		return cache.GetInstance().Set(ctx, coll+course.Name, course, time.Minute)
	}
	return err
}

func (s courseImpl) Update(ctx context.Context, course types.Course) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := storage.
		GetInstance().
		Update(ctx, coll, map[string]interface{}{"name": course.Name}, map[string]interface{}{"$set": &course})
	if err == nil {
		return cache.GetInstance().Delete(ctx, coll + course.Name)
	}
	return err
}

func (s courseImpl) FindAll(ctx context.Context) ([]types.Course, error) {
	var cs []types.Course
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	suffixKey := "all"

	err := cache.GetInstance().Get(ctx, coll+suffixKey, &cs)
	if err == nil {
		return cs, nil
	}
	err = storage.GetInstance().Find(ctx, coll, map[string]interface{}{}, &cs)
	if err == nil {
		_ = cache.GetInstance().Set(ctx, coll+suffixKey, cs, time.Minute)
		return cs, nil
	}
	return cs, err
}

func (s courseImpl) Delete(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := storage.GetInstance().Remove(ctx, coll, map[string]interface{}{"name": name})
	if err == nil {
		return cache.GetInstance().Delete(ctx, coll + name)
	}
	return err
}
