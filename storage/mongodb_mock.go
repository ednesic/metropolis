package storage

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type DataAccessLayerMock struct {
	mock.Mock
}

func (m *DataAccessLayerMock) StartSession() (Session, error) {
	args := m.Called()
	return &sessionMock{} , args.Error(1)
}

func (m *DataAccessLayerMock) Initialize(dbURI, dbName, collection string) error {
	 instance = m
	 return nil
}

func (m *DataAccessLayerMock) Insert(ctx context.Context, collName string, doc interface{}) error {
	args := m.Called(ctx, collName, doc)
	return args.Error(0)
}

func (m *DataAccessLayerMock) FindOne(ctx context.Context, collName string, query map[string]interface{}, doc interface{}) error {
	args := m.Called(ctx, collName, query, doc)
	return args.Error(0)
}

func (m *DataAccessLayerMock) Find(ctx context.Context, collName string, query map[string]interface{}, doc interface{}) error {
	args := m.Called(ctx, collName, query, doc)
	return args.Error(0)
}

func (m *DataAccessLayerMock) Count(ctx context.Context, collName string, query map[string]interface{}) (int64, error) {
	args := m.Called(ctx, collName, query)
	return int64(args.Int(0)), args.Error(1)
}

func (m *DataAccessLayerMock) Update(ctx context.Context, collName string, selector map[string]interface{}, update interface{}) error {
	args := m.Called(ctx, collName, selector, update)
	return args.Error(0)
}

func (m *DataAccessLayerMock) Remove(ctx context.Context, collName string, selector map[string]interface{}) error {
	args := m.Called(ctx, collName, selector)
	return args.Error(0)
}

func (m *DataAccessLayerMock) Disconnect() {}