package storage

import (
	"github.com/ednesic/coursemanagement/storage/mongodb"
	"github.com/ednesic/coursemanagement/types"
)

const coll = "course"



type CourseStorage interface {
	Create(types.Course) error
	Update(types.Course) error
	FindAll() ([]types.Course, error)
	Delete(types.Course) error
	FindOne(string) (types.Course, error)
}

type CourseStorageImpl struct {
	dal mongodb.DataAccessLayer
}

func NewCourseStorage(dbUri, dbName string) (CourseStorage, error) {
	dal, err := mongodb.NewMongoDAL(dbUri, dbName)
	return &CourseStorageImpl{
		dal,
	}, err
}

func (s *CourseStorageImpl) Create(c types.Course) error {
	return s.dal.Insert(coll, c)
}

func (s *CourseStorageImpl) Update(c types.Course) error {
	return s.dal.Upsert(coll, map[string]interface{}{"name": c.Name}, &c)
}

func (s *CourseStorageImpl) FindOne(name string) (c types.Course, err error) {
	err = s.dal.FindOne(coll, map[string]interface{}{"name": name}, &c)
	return
}

func (s *CourseStorageImpl) FindAll() (cs []types.Course, err error) {
	err = s.dal.Find(coll, map[string]interface{}{}, &cs)
	return
}

func (s *CourseStorageImpl) Delete(c types.Course) error {
	return s.dal.Remove(coll, map[string]interface{}{"name": c.Name})
}
