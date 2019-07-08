package storage

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	instance DataAccessLayer
	once     sync.Once
)

//MongoContext for specific session context for mongo-go-driver transaction
type MongoContext mongo.SessionContext

//DataAccessLayer is an interface for db connection
type DataAccessLayer interface {
	Insert(context.Context, string, interface{}) error
	Find(context.Context, string, map[string]interface{}, interface{}) error
	FindOne(context.Context, string, map[string]interface{}, interface{}) error
	Count(context.Context, string, map[string]interface{}) (int64, error)
	Update(context.Context, string, map[string]interface{}, interface{}) error
	Remove(context.Context, string, map[string]interface{}) error
	WithTransaction(context.Context, func(MongoContext) error) error
	Initialize(context.Context, string, string) error
	Disconnect()
}

//GetInstance to get database instance
func GetInstance() DataAccessLayer {
	once.Do(func() {
		if instance == nil {
			instance = &mongodbImpl{}
		}
	})
	return instance
}

type mongodbImpl struct {
	client *mongo.Client
	dbName string
}

func (m *mongodbImpl) Initialize(ctx context.Context, dbURI, dbName string) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	m.dbName = dbName
	m.client = client
	return nil
}

func (m *mongodbImpl) WithTransaction(ctx context.Context, fn func(MongoContext) error) error {
	return m.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err = fn(sessionContext)
		if err != nil {
			return sessionContext.AbortTransaction(sessionContext)
		}
		return sessionContext.CommitTransaction(sessionContext)
	})
}

// Insert stores documents in the collection
func (m *mongodbImpl) Insert(ctx context.Context, collName string, doc interface{}) error {
	_, err := m.client.Database(m.dbName).Collection(collName).InsertOne(ctx, doc)
	return err
}

// Find finds all documents in the collection
func (m *mongodbImpl) Find(ctx context.Context, collName string, query map[string]interface{}, doc interface{}) error {
	cur, err := m.client.Database(m.dbName).Collection(collName).Find(ctx, query)
	if err != nil {
		return err
	}

	resultv := reflect.ValueOf(doc)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		return errors.New("failed to return array response")
	}

	slicev := resultv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elem := slicev.Type().Elem()

	i := 0
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		elemp := reflect.New(elem)
		err := cur.Decode(elemp.Interface())
		if err != nil {
			return err
		}
		slicev = reflect.Append(slicev, elemp.Elem())
		slicev = slicev.Slice(0, slicev.Cap())
		i++
	}

	resultv.Elem().Set(slicev.Slice(0, i))
	return nil
}

// FindOne finds one document in mongo
func (m *mongodbImpl) FindOne(ctx context.Context, collName string, query map[string]interface{}, doc interface{}) error {
	return m.client.Database(m.dbName).Collection(collName).FindOne(ctx, query).Decode(doc)
}

// Update updates one or more documents in the collection
func (m *mongodbImpl) Update(ctx context.Context, collName string, selector map[string]interface{}, update interface{}) error {
	_, err := m.client.Database(m.dbName).Collection(collName).UpdateOne(ctx, selector, update)
	return err
}

// Remove one or more documents in the collection
func (m *mongodbImpl) Remove(ctx context.Context, collName string, selector map[string]interface{}) error {
	_, err := m.client.Database(m.dbName).Collection(collName).DeleteOne(ctx, selector)
	return err
}

// Count returns the number of documents of the query
func (m *mongodbImpl) Count(ctx context.Context, collName string, query map[string]interface{}) (int64, error) {
	return m.client.Database(m.dbName).Collection(collName).CountDocuments(ctx, query)
}

func (m *mongodbImpl) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_ = m.client.Disconnect(ctx)
}
