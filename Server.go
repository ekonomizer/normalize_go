package main

import (
	"net/http"
	Handlers "normalize_go/handlers"
	//"flag"
	//"io/ioutil"
	"fmt"
	"encoding/json"
	"strings"
	"strconv"

)

type Error struct {
	Error string `json:"error"`
}

type State struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Response struct {
	Address string `json:"address"`
	City string `json:"city"`
	State State `json:"state"`
	Zipcode int `json:"zipcode"`
	Country string `json:"country"`
}
var (
	States = map[string]string {
		"FL": "Florida",
	}
	Countries = map[string]string {
		"usa": "United States",
	}
)


//
// Request string - raw_address=8150+nw+53rd+street+suite+350-140,+doral,+fl,+33166, usa
// Response string -
// "address": "8150 NW 53rd St., Suite 350-140",
// "city": "Doral",
// "state": {
// "code": "FL",
// "name": "Florida"
//},
// "zipcode": 33166,
// "country": "United States"
func makeResponseStruct(params []string) Response {
	splitParams := strings.Split(params[0], ", ")
	address := strings.Split(splitParams[0], " ")
	addressStr := address[0] + " " + strings.ToUpper(address[1]) + " " + address[2] + " " + strings.Replace(address[3], "street", "St.", 1) + ", " + strings.Title(address[4]) + " " + address[5]
	cityStr := strings.Title(splitParams[1])
	zipCodeStr, err :=  strconv.Atoi(strings.Trim(splitParams[3], " "))
	if err != nil {
		panic(err)
	}
	countryStr := Countries[strings.Trim(splitParams[4], " ")]
	state := strings.ToUpper(splitParams[2])
	stateStr := State{ state, States[state] }

	fmt.Println("makeResponseStruct", addressStr)
	return Response{addressStr, cityStr, stateStr, zipCodeStr, countryStr}
}

func checkRequestErrors(r *http.Request, w http.ResponseWriter) string{
	r.ParseForm()
	params := r.Form
	fmt.Println("!!!!!", r.Form )
	fmt.Println("!!!!!", r.PostForm )
	var err string
	if !tryAuth(r) {
		err = "Can't authenticate!"
	}else if params["raw_address"] == nil {
		err = "raw_address required"
	} else if len(strings.Split(params["raw_address"][0], ", ")) != 5 {
		err = "wrong address format"
	}
	return err
}

func tryAuth(r *http.Request) bool {
	fmt.Println("Header", r.Header, r.Header["Login"],r.Header["Pass"])
	return r.Header["Login"] != nil && r.Header["Password"] != nil && r.Header["Login"][0] == "admin" && r.Header["Password"][0] == "admin"
}

func normalize(w http.ResponseWriter, r *http.Request) {
	paramsErr := checkRequestErrors(r, w)
	var err error
	var resp []byte

	if paramsErr != "" {
		w.WriteHeader(http.StatusNotFound)
		resp, err = json.Marshal(Error{paramsErr})
	} else {
		respStruct := makeResponseStruct(r.Form["raw_address"])
		resp, err = json.Marshal(respStruct)
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Can't make json"))
		return
	}
	w.Write(resp)
}

func HandleRequest(queueClass func() Handlers.HandlersQueue) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		hq := queueClass()
		hq.Run(w, r)
	}
}

func main() {
	http.HandleFunc("/address/normalize", HandleRequest(Handlers.NormalizeHandlersQueue))

	err := http.ListenAndServe(":12345", nil) // setting listening port
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
