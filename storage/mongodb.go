package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"sync"
	"time"
)

var instance DataAccessLayer
var once sync.Once


type DataAccessLayer interface {
	Insert(context.Context, string, interface{}) error
	Find(context.Context, string, map[string]interface{}, interface{}) error
	FindOne(context.Context, string, map[string]interface{}, interface{}) error
	Count(context.Context, string, map[string]interface{}) (int64, error)
	Update(context.Context, string, map[string]interface{}, interface{}) error
	Remove(context.Context, string, map[string]interface{}) error
	Initialize(string, string, string) error
	StartSession() (Session, error)
	Disconnect()
}

func GetInstance() DataAccessLayer {
	once.Do(func() {
		if instance == nil {
			instance = &impl{}
		}
	})
	return instance
}

type impl struct {
	collection *mongo.Collection
	dbName  string
	client *mongo.Client
}

func (m *impl) Initialize(dbURI, dbName, collection string) error {
	ctx := context.Background()
	//ctx, _ := context.WithTimeout(context.Background(), 2 * time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	m.collection = client.Database(dbName).Collection(collection)
	m.dbName = dbName
	m.client = client
	return nil
}

func (m *impl) StartSession() (Session, error) {
	var err error
	var session mongo.Session
	if session, err = m.client.StartSession(); err != nil {
		return nil, err
	}
	if err = session.StartTransaction(); err != nil {
		return nil, err
	}
	return &sss{session}, nil
}

// Insert stores documents in the collection
func (m *impl) Insert(ctx context.Context, collName string, doc interface{}) error {
	_, err := m.collection.InsertOne(ctx, doc)
	return err
}

// Find finds all documents in the collection
func (m *impl) Find(ctx context.Context, collName string, query map[string]interface{}, doc interface{}) error {
	cur, err := m.collection.Find(ctx, query)
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
func (m *impl) FindOne(ctx context.Context, collName string, query map[string]interface{}, doc interface{}) error {
	return m.collection.FindOne(ctx, query).Decode(doc)
}

// Update updates one or more documents in the collection
func (m *impl) Update(ctx context.Context, collName string, selector map[string]interface{}, update interface{}) error {
	_, err := m.collection.UpdateOne(ctx, selector, map[string]interface{}{"$set": update})
	return err
}

// Remove one or more documents in the collection
func (m *impl) Remove(ctx context.Context, collName string, selector map[string]interface{}) error {
	_, err := m.collection.DeleteOne(ctx, selector)
	return err
}

// Count returns the number of documents of the query
func (m *impl) Count(ctx context.Context, collName string, query map[string]interface{}) (int64, error) {
	return m.collection.CountDocuments(ctx, query)
}

func (m *impl) Disconnect() {
	ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)
	_ = m.client.Disconnect(ctx)
}