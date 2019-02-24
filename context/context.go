package context

import (
	"context"
	"fmt"
	gcontext "github.com/gorilla/context"
	"net/http"
)

type ctxKey int


const (
	tokenContextKey ctxKey = iota
	errorContextKey
	delayedHandlerKey
	preventUnlockKey
	appContextKey
	reqBodyKey
)

func GetRequestError(r *http.Request) error {
	if r == nil {
		return nil
	}
	val := gcontext.GetAll(r)
	fmt.Println(val)
	if v, ok := r.Context().Value(errorContextKey).(error); ok {
		return v
	}
	return nil
}

func AddRequestError(r *http.Request, err error) {
	if err == nil {
		return
	}
	fmt.Println(gcontext.Get(r, delayedHandlerKey))
	gcontext.Set(r, delayedHandlerKey, err)
	fmt.Println(gcontext.Get(r, delayedHandlerKey))
	*r = *r.WithContext(context.WithValue(r.Context(), errorContextKey, err))
}


func Clear(r *http.Request) {
	if r == nil {
		return
	}
	newReq := r.WithContext(context.Background())
	*r = *newReq
}