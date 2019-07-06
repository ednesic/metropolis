package courseservice

import (
	"github.com/ednesic/coursemanagement/types"
	"github.com/stretchr/testify/mock"
)

//Mock is a mocked structure for course service
type Mock struct {
	mock.Mock
}

//InitMock to initialize mock befores tests
func (s *Mock) InitMock() {
	instance = s
}

//FindOne is a mock for course service findOne
func (s *Mock) FindOne(name string) (c types.Course, err error) {
	args := s.Called(name)
	return args.Get(0).(types.Course), args.Error(1)
}

//Create is a mock for course service create
func (s *Mock) Create(course types.Course) error {
	args := s.Called(course)
	return args.Error(0)
}

//Update is a mock for course service update
func (s *Mock) Update(course types.Course) error {
	args := s.Called(course)
	return args.Error(0)
}

//FindAll is a mock for course service finaAll
func (s *Mock) FindAll() (cs []types.Course, err error) {
	args := s.Called()
	return args.Get(0).([]types.Course), args.Error(1)
}

//Delete is a mock for course service delete
func (s *Mock) Delete(name string) error {
	args := s.Called(name)
	return args.Error(0)
}
