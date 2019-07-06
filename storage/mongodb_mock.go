package storage

import (
	"context"

	"github.com/stretchr/testify/mock"
)

//DataAccessLayerMock is a mock for db connection
type DataAccessLayerMock struct {
	mock.Mock
}

//WithTransaction is a mock for db WithTransaction
func (m *DataAccessLayerMock) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

//Initialize is a mock for db Initialize
func (m *DataAccessLayerMock) Initialize(ctx context.Context, dbURI, dbName string) error {
	instance = m
	return nil
}

//Insert is a mock for db Insert
func (m *DataAccessLayerMock) Insert(ctx context.Context, collName string, doc interface{}) error {
	args := m.Called(ctx, collName, doc)
	return args.Error(0)
}

//FindOne is a mock for db FindOne
func (m *DataAccessLayerMock) FindOne(ctx context.Context, collName string, query map[string]interface{}, doc interface{}) error {
	args := m.Called(ctx, collName, query, doc)
	return args.Error(0)
}

//Find is a mock for db Find
func (m *DataAccessLayerMock) Find(ctx context.Context, collName string, query map[string]interface{}, doc interface{}) error {
	args := m.Called(ctx, collName, query, doc)
	return args.Error(0)
}

//Count is a mock for db Count
func (m *DataAccessLayerMock) Count(ctx context.Context, collName string, query map[string]interface{}) (int64, error) {
	args := m.Called(ctx, collName, query)
	return int64(args.Int(0)), args.Error(1)
}

//Update is a mock for Update
func (m *DataAccessLayerMock) Update(ctx context.Context, collName string, selector map[string]interface{}, update interface{}) error {
	args := m.Called(ctx, collName, selector, update)
	return args.Error(0)
}

//Remove is a mock for Remove
func (m *DataAccessLayerMock) Remove(ctx context.Context, collName string, selector map[string]interface{}) error {
	args := m.Called(ctx, collName, selector)
	return args.Error(0)
}

//Disconnect is a mock for Disconnect
func (m *DataAccessLayerMock) Disconnect() {}
