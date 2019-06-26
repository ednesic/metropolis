package services

import (
	"github.com/ednesic/coursemanagement/types"
	"github.com/stretchr/testify/mock"
)

type CourseServiceMock struct {
	mock.Mock
}

func (s *CourseServiceMock) FindOne(name string) (c types.Course, err error) {
	args := s.Called(name)
	return args.Get(0).(types.Course), args.Error(1)
}

func (s *CourseServiceMock) Create(course types.Course) error {
	args := s.Called(course)
	return args.Error(0)
}

func (s *CourseServiceMock) Update(course types.Course) error {
	args := s.Called(course)
	return args.Error(0)
}

func (s *CourseServiceMock) FindAll() (cs []types.Course, err error) {
	args := s.Called()
	return args.Get(0).([]types.Course), args.Error(1)
}

func (s *CourseServiceMock) Delete(name string) error {
	args := s.Called(name)
	return args.Error(0)
}