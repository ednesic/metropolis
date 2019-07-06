package storage

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	//ErrNotFound for database not found documents
	ErrNotFound = mongo.ErrNoDocuments
)
