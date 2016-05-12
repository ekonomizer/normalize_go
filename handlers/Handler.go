package handlers

import "net/http"

type Handler interface {
	Handle(http.ResponseWriter, *http.Request, HandlersQueue)
	Error(http.ResponseWriter, *http.Request, HandlersQueue)
}
