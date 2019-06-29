package courseservice

import (
	redisrepository "github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/go-redis/cache"
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

type Impl struct {}

func GetInstance() CourseService {
	once.Do(func() {
		if instance == nil {
			instance = &Impl{}
		}
	})
	return instance
}

func (s *Impl) FindOne(name string) (c types.Course, err error) {
	var mgoErr error
	if err := redisrepository.GetInstance().Get(coll+name, &c); err != nil {
		if mgoErr = storage.GetInstance().FindOne(coll, map[string]interface{}{"name": name}, &c); mgoErr == nil {
			return c, redisrepository.GetInstance().Set(&cache.Item{Key: coll + name, Object: c, Expiration: time.Minute})
		}
	}
	return c, mgoErr
}

func (s *Impl) Create(course types.Course) error {
	err := storage.GetInstance().Insert(coll, course)
	if err == nil {
		return redisrepository.GetInstance().Set(&cache.Item{Key: coll + course.Name, Object: course, Expiration: time.Minute})
	}
	return err
}

func (s *Impl) Update(course types.Course) error {
	err := storage.GetInstance().Update(coll, map[string]interface{}{"name": course.Name}, &course)
	if err == nil {
		return redisrepository.GetInstance().Delete(coll + course.Name)
	}
	return err
}

func (s *Impl) FindAll() ([]types.Course, error) {
	cs := []types.Course{}
	var mgoErr error
	suffixKey := "all"

	if cacheErr := redisrepository.GetInstance().Get(coll + suffixKey, &cs); cacheErr != nil {
		if mgoErr = storage.GetInstance().Find(coll, map[string]interface{}{}, &cs); mgoErr == nil {
			return cs, redisrepository.GetInstance().Set(&cache.Item{Key: coll + suffixKey, Object: cs, Expiration: time.Minute})
		}
	}
	return cs, mgoErr
}

func (s *Impl) Delete(name string) error {
	err := storage.GetInstance().Remove(coll, map[string]interface{}{"name": name})
	if err == nil {
		return redisrepository.GetInstance().Delete(coll + name)
	}
	return err
}