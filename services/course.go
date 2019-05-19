package services

import (
	redis "github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/go-redis/cache"
	"time"
)

const coll = "course"

type CourseService interface {
	Create(types.Course) error
	Update(types.Course) error
	FindAll() ([]types.Course, error)
	Delete(types.Course) error
	FindOne(string) (types.Course, error)
}

type CourseServiceImpl struct {
	dal   storage.DataAccessLayer
	cache redis.RedisClient
}

func NewCourseService(dbUri, dbName string, redisHosts map[string]string) (CourseService, error) {
	dal, err := storage.NewMongoConnectDAL(dbUri, dbName)
	return &CourseServiceImpl{
		dal,
		redis.NewRedisClient(redisHosts),
	}, err
}

func (s *CourseServiceImpl) FindOne(name string) (c types.Course, err error) {
	var mgoErr error
	if err := s.cache.Get(coll+name, &c); err != nil {
		if mgoErr = s.dal.FindOne(coll, map[string]interface{}{"name": name}, &c); mgoErr == nil {
			return c, s.cache.Set(&cache.Item{Key: coll + name, Object: c, Expiration: time.Hour})
		}
	}
	return c, mgoErr
}

func (s *CourseServiceImpl) Create(course types.Course) error {
	err := s.dal.Insert(coll, course)
	if err == nil {
		return s.cache.Set(&cache.Item{Key: coll + course.Name, Object: course, Expiration: time.Hour})
	}
	return err
}

func (s *CourseServiceImpl) Update(course types.Course) error {
	err := s.dal.Update(coll, map[string]interface{}{"name": course.Name}, &course)
	if err == nil {
		return s.cache.Delete(coll + course.Name)
	}
	return err
}

func (s *CourseServiceImpl) FindAll() ([]types.Course, error) {
	cs := []types.Course{}
	var mgoErr error
	suffixKey := "all"

	if cacheErr := s.cache.Get(coll + suffixKey, &cs); cacheErr != nil {
		if mgoErr = s.dal.Find(coll, map[string]interface{}{}, &cs); mgoErr == nil {
			return cs, s.cache.Set(&cache.Item{Key: coll + suffixKey, Object: cs, Expiration: time.Hour})
		}
	}
	return cs, mgoErr
}

func (s *CourseServiceImpl) Delete(course types.Course) error {
	err := s.dal.Remove(coll, map[string]interface{}{"name": course.Name})
	if err == nil {
		return s.cache.Delete(coll + course.Name)
	}
	return err
}
