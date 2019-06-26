package services

import (
	"github.com/ednesic/coursemanagement/repositorymanager"
	"github.com/ednesic/coursemanagement/types"
	"github.com/go-redis/cache"
	"time"
)

const coll = "course"

type CourseService interface {
	Create(types.Course) error
	Update(types.Course) error
	FindAll() ([]types.Course, error)
	Delete(string) error
	FindOne(string) (types.Course, error)
}

type CourseServiceImpl struct {}

func NewCourseService() CourseService {
	return &CourseServiceImpl{}
}

func (s *CourseServiceImpl) FindOne(name string) (c types.Course, err error) {
	var mgoErr error
	if err := repositorymanager.Redis.Get(coll+name, &c); err != nil {
		if mgoErr = repositorymanager.Dal.FindOne(coll, map[string]interface{}{"name": name}, &c); mgoErr == nil {
			return c, repositorymanager.Redis.Set(&cache.Item{Key: coll + name, Object: c, Expiration: time.Hour})
		}
	}
	return c, mgoErr
}

func (s *CourseServiceImpl) Create(course types.Course) error {
	err := repositorymanager.Dal.Insert(coll, course)
	if err == nil {
		return repositorymanager.Redis.Set(&cache.Item{Key: coll + course.Name, Object: course, Expiration: time.Hour})
	}
	return err
}

func (s *CourseServiceImpl) Update(course types.Course) error {
	err := repositorymanager.Dal.Update(coll, map[string]interface{}{"name": course.Name}, &course)
	if err == nil {
		return repositorymanager.Redis.Delete(coll + course.Name)
	}
	return err
}

func (s *CourseServiceImpl) FindAll() ([]types.Course, error) {
	cs := []types.Course{}
	var mgoErr error
	suffixKey := "all"

	if cacheErr := repositorymanager.Redis.Get(coll + suffixKey, &cs); cacheErr != nil {
		if mgoErr = repositorymanager.Dal.Find(coll, map[string]interface{}{}, &cs); mgoErr == nil {
			return cs, repositorymanager.Redis.Set(&cache.Item{Key: coll + suffixKey, Object: cs, Expiration: time.Hour})
		}
	}
	return cs, mgoErr
}

func (s *CourseServiceImpl) Delete(name string) error {
	err := repositorymanager.Dal.Remove(coll, map[string]interface{}{"name": name})
	if err == nil {
		return repositorymanager.Redis.Delete(coll + name)
	}
	return err
}