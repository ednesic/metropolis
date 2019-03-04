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
	cache *cache.Codec
}

func NewCourseService(dbUri, dbName string, redisHosts map[string]string) (CourseService, error) {
	dal, err := storage.NewMongoDAL(dbUri, dbName)
	return &CourseServiceImpl{
		dal,
		redis.NewRedisClient(redisHosts),
	}, err
}

func (s *CourseServiceImpl) FindOne(name string) (c types.Course, err error) {
	if cacheErr := s.cache.Get(coll+name, &c); cacheErr != nil {
		err = s.dal.FindOne(coll, map[string]interface{}{"name": name}, &c)
		cacheErr = s.cache.Set(&cache.Item{Key: coll + name, Object: c, Expiration: time.Hour})
	}
	return
}

func (s *CourseServiceImpl) Create(course types.Course) error {
	err := s.dal.Insert(coll, course)
	_ = s.cache.Set(&cache.Item{Key: coll + course.Name, Object: course, Expiration: time.Hour}) //resolve cache error
	return err

}

func (s *CourseServiceImpl) Update(course types.Course) error {
	err := s.dal.Upsert(coll, map[string]interface{}{"name": course.Name}, &course)
	_ = s.cache.Delete(coll + course.Name)
	return err
}

func (s *CourseServiceImpl) FindAll() (cs []types.Course, err error) {
	suffixKey := "all"
	if cacheErr := s.cache.Get(coll + suffixKey, &cs); cacheErr != nil {
		err = s.dal.Find(coll, map[string]interface{}{}, &cs)
		cacheErr = s.cache.Set(&cache.Item{Key: coll + suffixKey, Object: cs, Expiration: time.Hour})
	}
	return
}

func (s *CourseServiceImpl) Delete(course types.Course) error {
	err := s.dal.Remove(coll, map[string]interface{}{"name": course.Name})
	_ = s.cache.Delete(coll + course.Name)
	return err
}
