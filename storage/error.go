package storage

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotFound = mongo.ErrNoDocuments
)

type ErrorReponse struct {
	Message string `json:"message,omitempty"`
}