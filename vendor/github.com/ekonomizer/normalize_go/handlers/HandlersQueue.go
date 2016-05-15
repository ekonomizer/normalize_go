package handlers

import (
	"net/http"
	"fmt"
)

type HandlersQueue struct {
	Handlers []Handler
	Errors map[int]string
}

func (s HandlersQueue) Run(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("start HandlersQueue", "handler len: ", len(s.Handlers), " errors len: ",len(s.Errors))
	var previusHandler Handler
	for _, handler := range s.Handlers {
		if len(s.Errors) == 0 {
			previusHandler = handler
			handler.Handle(w, r, s)
		} else {
			fmt.Println("Handle Error")
			previusHandler.Error(w, r, s)
			return
		}
	}
	if len(s.Errors) != 0 {
		fmt.Println("Handle Error")
		previusHandler.Error(w, r, s)
	}
}

func NormalizeHandlersQueue() HandlersQueue {
	return HandlersQueue{
		[]Handler{
			AuthHandler{},
			NormalizeHandler{},
		},map[int]string {},
	}
}
