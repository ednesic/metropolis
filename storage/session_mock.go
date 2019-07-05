package storage

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type sessionMock struct {
	mock.Mock
}

func (s *sessionMock) WithTransaction(fn func(context.Context) error) error {
	args := s.Called(fn)
	return args.Error(0)
}

func (s *sessionMock) Commit() error {
	args := s.Called()
	return args.Error(0)
}

func (s *sessionMock) Abort() error {
	args := s.Called()
	return args.Error(0)
}