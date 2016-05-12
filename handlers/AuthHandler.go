package handlers

import (
	"net/http"
	"fmt"
)

const AuthorizationHeaderLength = 500

var AuthData = map[string]string {
	"admin": "admin",
}
//admin:admin
//YWRtaW46YWRtaW4=

type AuthHandler struct {

}

func (s AuthHandler) Handle(w http.ResponseWriter, r *http.Request, hq HandlersQueue)  {
	fmt.Println("Start AuthHandler")
	username, password, ok := r.BasicAuth()
	fmt.Println("AuthHandler", "login:", username, "password:", password, ok)

	if !ok || AuthData[username] != password {
		hq.Errors[http.StatusUnauthorized] = http.StatusText(http.StatusUnauthorized)
		return
	}
	fmt.Println("Success basic auth")
}

func (s AuthHandler) Error(w http.ResponseWriter, r *http.Request, hq HandlersQueue)  {
	fmt.Println("AuthHandler: response error", hq.Errors)
	for k, v := range hq.Errors {
		http.Error(w, v, k)
	}
}