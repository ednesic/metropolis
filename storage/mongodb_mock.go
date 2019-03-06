package storage

import (
"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

type DataAccessLayerMock struct {
	mock.Mock
}

func (m *DataAccessLayerMock) Insert(collName string, doc interface{}) error {
	args := m.Called(collName, doc)
	return args.Error(0)
}

func (m *DataAccessLayerMock) FindOne(collName string, query bson.M, doc interface{}) error {
	args := m.Called(collName, query, doc)
	return args.Error(0)
}

func (m *DataAccessLayerMock) FindWithLimit(collName, sortQuery string, limit int, query bson.M, doc interface{}) error {
	args := m.Called(collName, sortQuery, query, doc, limit)
	return args.Error(0)
}

func (m *DataAccessLayerMock) Find(collName string, query bson.M, doc interface{}) error {
	args := m.Called(collName, query, doc)
	return args.Error(0)
}

func (m *DataAccessLayerMock) Count(collName string, query bson.M) (int, error) {
	args := m.Called(collName, query)
	return args.Int(0), args.Error(1)
}

func (m *DataAccessLayerMock) Update(collName string, selector bson.M, update interface{}) error {
	args := m.Called(collName, selector, update)
	return args.Error(0)
}

func (m *DataAccessLayerMock) Upsert(collName string, selector bson.M, update interface{}) error {
	args := m.Called(collName, selector, update)
	return args.Error(0)
}

func (m *DataAccessLayerMock) Remove(collName string, selector bson.M) error {
	args := m.Called(collName, selector)
	return args.Error(0)
}

func (m *DataAccessLayerMock) EnsureIndex(collName, field string, isUnique bool) error {
	args := m.Called(collName, field, isUnique)
	return args.Error(0)
}
