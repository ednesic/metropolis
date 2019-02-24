package services

import (
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
)

type CourseService interface {
	Create(types.Course) error
	Update(types.Course) error
	FindAll() ([]types.Course,error)
	Delete(types.Course) error
	FindOne(string) (types.Course, error)
}

type CourseServiceImpl struct {
	storage storage.CourseStorage
}

func NewCourseService(dbUri, dbName string) (CourseService, error) {
	st, err := storage.NewCourseStorage(dbUri, dbName)
	return &CourseServiceImpl{
		st,
	}, err
}

func (s *CourseServiceImpl) FindOne(name string) (types.Course, error){
	return s.storage.FindOne(name)
}

func (s *CourseServiceImpl) Create(course types.Course) error {
	return s.storage.Create(course)
}

func (s *CourseServiceImpl) Update(course types.Course) error {
	return s.storage.Update(course)
}

func (s *CourseServiceImpl) FindAll() ([]types.Course, error) {
	return s.storage.FindAll()
}

func (s *CourseServiceImpl) Delete(course types.Course) error {
	return s.storage.Delete(course)
}