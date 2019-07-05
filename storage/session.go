package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Session interface {
	WithTransaction(fn func(context.Context) error) error
	Commit() error
	Abort() error
}

type sss struct {
	mongo.Session
}

func (s *sss) WithTransaction(fn func(context.Context) error) error {
	return mongo.WithSession(context.Background(), s, func(sc mongo.SessionContext) error {
		return fn(sc)
	})
}

func (s *sss) Commit() error {
	return mongo.WithSession(context.Background(), s, func(sc mongo.SessionContext) error {
		return s.CommitTransaction(sc)
	})
}

func (s *sss) Abort() error {
	err := mongo.WithSession(context.Background(), s, func(sc mongo.SessionContext) error {
		return s.AbortTransaction(sc)
	})
	s.EndSession(context.Background())
	return err
}