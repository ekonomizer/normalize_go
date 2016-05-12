package handlers

import (
	"net/http"
	"fmt"
	"strings"
	"net/url"
	"encoding/json"
	"normalize_go/utils"
)

const GoogleApiKey = "AIzaSyBhOuB_Vzl34lQzuzGrD3OiNG-xyDJKElg"
const GoogleMapApiLink = "https://maps.googleapis.com/maps/api/geocode/json?address="

type NormalizeHandler struct {

}
func makeGoogleMapApiUrl(address string) (string, string){
	u, err := url.Parse(GoogleMapApiLink)
	if err != nil {
		return "", err.Error()
	}
	//u.Scheme = "https"
	//u.Host = "google.com"
	q := u.Query()
	q.Set("address", address)
	q.Set("key", GoogleApiKey)
	u.RawQuery = q.Encode()
	return u.String(), ""
}

type Response struct {
	Results []Result
}

type Result struct {
	Address_components []Address `json:"address_components"`
}

type Address struct {
	Long_name string `json:"long_name"`
	Short_name string  `json:"short_name"`
	Types []string  `json:"types"`
}

type State struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type NormalizeResponse struct {
	Address string `json:"address"`
	City string `json:"city"`
	State State `json:"state"`
	Zipcode string `json:"zipcode"`
	Country string `json:"country"`
}


func getGoogleApiAddress(urlStr string)(string, string) {
	// Http request to ggolge api
	googleApiResponse, err := http.Get(urlStr)
	if err != nil {
		return "", "error google map api request: " + err.Error()
	}
	defer googleApiResponse.Body.Close()

	// Decode google api response string
	response := &Response{}
	errJson := json.NewDecoder(googleApiResponse.Body).Decode(&response)
	if errJson != nil {
		return "", "error when parse json address: " + errJson.Error()
	}

	fmt.Println(response)

	// Make response
	normResponse := normalizeResponse(*response)
	resp, errMarshal := json.Marshal(normResponse)
	if errMarshal != nil {
		return "", "error when Marshal normResponse: " + errMarshal.Error()
	}

	return string(resp), ""
}

func normalizeResponse(response Response) (NormalizeResponse) {
	var result Result = response.Results[0]
	address := ""
	var v Address
	var street_number string
	var route string
	var subpremise string
	var locality string
	var administrative_area_level_1_short string
	var administrative_area_level_1_long string
	var country string
	var postal_code string

	for _, v = range result.Address_components{
		if utils.ArrayIndexOfStr(v.Types, "street_number") != -1 {
			street_number = v.Short_name
		} else if utils.ArrayIndexOfStr(v.Types, "route") != -1 {
			route = v.Short_name
		} else if utils.ArrayIndexOfStr(v.Types, "subpremise") != -1 {
			subpremise = v.Short_name
		} else if utils.ArrayIndexOfStr(v.Types, "locality") != -1 {
			locality = v.Short_name
		} else if utils.ArrayIndexOfStr(v.Types, "administrative_area_level_1") != -1 {
			administrative_area_level_1_short = v.Short_name
			administrative_area_level_1_long = v.Long_name
		} else if utils.ArrayIndexOfStr(v.Types, "country") != -1 {
			country = v.Short_name
		} else if utils.ArrayIndexOfStr(v.Types, "postal_code") != -1 {
			postal_code = v.Short_name
		}
	}
	address += street_number + " " + route + "., Sutie " + subpremise
	state := State{administrative_area_level_1_short, administrative_area_level_1_long}

	return NormalizeResponse{address, locality, state, postal_code, country}
}

//https://maps.googleapis.com/maps/api/geocode/json?address=8150+nw+53rd+street+suite+350-140,+doral,+33166,%20usa&key=AIzaSyBhOuB_Vzl34lQzuzGrD3OiNG-xyDJKElg
func (s NormalizeHandler) Handle(w http.ResponseWriter, r *http.Request, hq HandlersQueue)  {
	r.ParseForm()
	params := r.Form
	fmt.Println("Start NormalizeHandler")

	if params["raw_address"] == nil {
		hq.Errors[http.StatusNotFound] = "raw_address required"
		return
	} else if len(strings.Split(params["raw_address"][0], ", ")) != 5 {
		hq.Errors[http.StatusNotFound] = "wrong address format"
		return
	}

	urlStr, makeUrlErr := makeGoogleMapApiUrl(params["raw_address"][0])
	fmt.Println("Make url:", urlStr)
	if makeUrlErr != "" {
		hq.Errors[http.StatusNotFound] = "Can't make url"
		return
	}

	result , addressErr := getGoogleApiAddress(urlStr)
	if addressErr != "" {
		hq.Errors[http.StatusNotFound] = "Can't convert google json"
		return
	}

	w.Write([]byte(result))
}

func (s NormalizeHandler) Error(w http.ResponseWriter, r *http.Request, hq HandlersQueue)  {
	fmt.Println("NormalizeHandler: response error", hq.Errors)
	for k, v := range hq.Errors {
		http.Error(w, v, k)
	}
}


