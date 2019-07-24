package courseservice

import (
	"context"
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
func (s *Mock) FindOne(ctx context.Context, name string) (c types.Course, err error) {
	args := s.Called(ctx, name)
	return args.Get(0).(types.Course), args.Error(1)
}

//Create is a mock for course service create
func (s *Mock) Create(ctx context.Context, course types.Course) error {
	args := s.Called(ctx, course)
	return args.Error(0)
}

//Update is a mock for course service update
func (s *Mock) Update(ctx context.Context, course types.Course) error {
	args := s.Called(ctx, course)
	return args.Error(0)
}

//FindAll is a mock for course service finaAll
func (s *Mock) FindAll(ctx context.Context) (cs []types.Course, err error) {
	args := s.Called(ctx)
	return args.Get(0).([]types.Course), args.Error(1)
}

//Delete is a mock for course service delete
func (s *Mock) Delete(ctx context.Context, name string) error {
	args := s.Called(ctx, name)
	return args.Error(0)
}
