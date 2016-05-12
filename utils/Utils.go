package utils

import (
	"io"
	"io/ioutil"
)
const READ_LIMIT_BYTES = 1e6
func ArrayIndexOfStr(arr []string, str string)(int) {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return i
		}
	}
	return -1
}

func ParseHttpResponse(body io.ReadCloser) (string, error) {
	lr := &io.LimitedReader{body, READ_LIMIT_BYTES}
	contents, err := ioutil.ReadAll(lr)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

//// Parse google api response string
//contents, errParseHttpResponse := utils.ParseHttpResponse(googleApiResponse.Body)
//if errParseHttpResponse != nil {
//	return "", "error when ParseHttpResponse google api: " + errParseHttpResponse.Error()
//}
//
//// Parse response string(json) to struct Response
//response := &Response{}
//reader := strings.NewReader(contents)
//errJson := json.NewDecoder(reader).Decode(&response)
//if errJson != nil {
//	return "", "error when parse json address: " + errJson.Error()
//}