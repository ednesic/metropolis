package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DataAccessLayer interface {
	Insert(collectionName string, doc interface{}) error
	Find(collName string, query bson.M, doc interface{}) error
	FindOne(collName string, query bson.M, doc interface{}) error
	Count(collName string, query bson.M) (int, error)
	FindWithLimit(collName, sortQuery string, limit int, query bson.M, doc interface{}) error
	Update(collName string, selector bson.M, update interface{}) error
	Upsert(collName string, selector bson.M, update interface{}) error
	Remove(collName string, selector bson.M) error
	EnsureIndex(collName, field string, isUnique bool) error
}

type MongoDAL struct {
	session *mgo.Session
	dbName  string
}

// NewMongoDAL creates a MongoDAL
func NewMongoDAL(dbURI string, dbName string) (DataAccessLayer, error) {
	session, err := mgo.Dial(dbURI)
	mongo := &MongoDAL{
		session: session,
		dbName:  dbName,
	}
	return mongo, err
}

// Insert stores documents in the collection
func (m *MongoDAL) Insert(collName string, doc interface{}) error {
	session := m.session.Clone()
	defer session.Close()
	return session.DB(m.dbName).C(collName).Insert(doc)
}

// Find finds all documents in the collection
func (m *MongoDAL) Find(collName string, query bson.M, doc interface{}) error {
	session := m.session.Clone()
	defer session.Close()
	return session.DB(m.dbName).C(collName).Find(query).All(doc)
}

// FindWithLimit finds all documents in the collection with a limit
func (m *MongoDAL) FindWithLimit(collName, sortQuery string, limit int, query bson.M, doc interface{}) error {
	session := m.session.Clone()
	defer session.Close()
	return session.DB(m.dbName).C(collName).Find(query).Sort(sortQuery).Limit(limit).All(doc)
}

// FindOne finds one document in mongo
func (m *MongoDAL) FindOne(collName string, query bson.M, doc interface{}) error {
	session := m.session.Clone()
	defer session.Close()
	return session.DB(m.dbName).C(collName).Find(query).One(doc)
}

// Update updates one or more documents in the collection
func (m *MongoDAL) Update(collName string, selector bson.M, update interface{}) error {
	session := m.session.Clone()
	defer session.Close()
	return session.DB(m.dbName).C(collName).Update(selector, update)
}

// Upsert one or more documents in the collection
func (m *MongoDAL) Upsert(collName string, selector bson.M, update interface{}) error {
	session := m.session.Clone()
	defer session.Close()
	_, err := session.DB(m.dbName).C(collName).Upsert(selector, update)
	return err
}

// Remove one or more documents in the collection
func (m *MongoDAL) Remove(collName string, selector bson.M) error {
	session := m.session.Clone()
	defer session.Close()
	return session.DB(m.dbName).C(collName).Remove(selector)
}

// Count returns the number of documents of the query
func (m *MongoDAL) Count(collName string, query bson.M) (int, error) {
	session := m.session.Clone()
	defer session.Close()
	return session.DB(m.dbName).C(collName).Find(query).Count()
}

// EnsureIndex ensures that a key is indexed
func (m *MongoDAL) EnsureIndex(collName, field string, isUnique bool) error {
	session := m.session.Clone()
	defer session.Close()
	index := mgo.Index{
		Key:    []string{field},
		Unique: isUnique,
	}
	if err := session.DB(m.dbName).C(collName).EnsureIndex(index); err != nil {
		return err
	}
	return nil
}
