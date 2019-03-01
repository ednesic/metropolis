package services

import (
	"github.com/ednesic/coursemanagement/storage/mongodb"
	"github.com/ednesic/coursemanagement/types"
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
	dal mongodb.DataAccessLayer
}

func NewCourseService(dbUri, dbName string) (CourseService, error) {
	dal, err := mongodb.NewMongoDAL(dbUri, dbName)
	return &CourseServiceImpl{
		dal,
	}, err
}

func (s *CourseServiceImpl) FindOne(name string) (c types.Course, err error) {
	err = s.dal.FindOne(coll, map[string]interface{}{"name": name}, &c)
	return
}

func (s *CourseServiceImpl) Create(course types.Course) error {
	return s.dal.Insert(coll, course)
}

func (s *CourseServiceImpl) Update(course types.Course) error {
	return s.dal.Upsert(coll, map[string]interface{}{"name": course.Name}, &course)
}

func (s *CourseServiceImpl) FindAll() (cs []types.Course, err error) {
	err = s.dal.Find(coll, map[string]interface{}{}, &cs)
	return
}

func (s *CourseServiceImpl) Delete(course types.Course) error {
	return s.dal.Remove(coll, map[string]interface{}{"name": course.Name})
}
