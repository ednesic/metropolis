package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type DataAccessLayer interface {
	Insert(collName string, doc interface{}) error
	Find(collName string, query map[string]interface{}, doc interface{}) error
	FindOne(collName string, query map[string]interface{}, doc interface{}) error
	Count(collName string, query map[string]interface{}) (int64, error)
	Update(collName string, selector map[string]interface{}, update interface{}) error
	Remove(collName string, selector map[string]interface{}) error
}


type MongoConnectDAL struct {
	collection *mongo.Collection
	dbName  string
}

func NewMongoConnectDAL(dbURI string, dbName string) (DataAccessLayer, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, err
	}
	collection := client.Database(dbName).Collection("course")

	newMongo := &MongoConnectDAL{
		collection: collection,
		dbName:  dbName,
	}
	return newMongo, err
}

// Insert stores documents in the collection
func (m *MongoConnectDAL) Insert(collName string, doc interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)
	_, err := m.collection.InsertOne(ctx, doc)
	return err
}

// Find finds all documents in the collection
func (m *MongoConnectDAL) Find(collName string, query map[string]interface{}, doc interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)
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
func (m *MongoConnectDAL) FindOne(collName string, query map[string]interface{}, doc interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)
	return m.collection.FindOne(ctx, query).Decode(doc)
}

// Update updates one or more documents in the collection
func (m *MongoConnectDAL) Update(collName string, selector map[string]interface{}, update interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)

	_, err := m.collection.UpdateOne(ctx, selector, map[string]interface{}{"$set": update})
	return err
}

// Remove one or more documents in the collection
func (m *MongoConnectDAL) Remove(collName string, selector map[string]interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)

	_, err := m.collection.DeleteOne(ctx, selector)
	return err
}

// Count returns the number of documents of the query
func (m *MongoConnectDAL) Count(collName string, query map[string]interface{}) (int64, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)
	return m.collection.CountDocuments(ctx, query)
}
