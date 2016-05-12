package main

import (
	"net/http"
	"github.com/ekonomizer/normalize_go/handlers"
	"fmt"
	"runtime"
)

func HandleRequest(queueClass func() handlers.HandlersQueue) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		hq := queueClass()
		hq.Run(w, r)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/address/normalize", HandleRequest(handlers.NormalizeHandlersQueue))

	err := http.ListenAndServe(":12345", nil) // setting listening port
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
