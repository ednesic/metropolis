package storage

import (
	"gopkg.in/mgo.v2"
)

var (
	ErrNotFound = mgo.ErrNotFound
)

type ErrorReponse struct {
	Message string `json:"message,omitempty"`
}