package gomite

import (
	"encoding/json"
	"net/http"
)

func ParseJson(req *http.Request, dst interface{}) (err error) {
	err = json.NewDecoder(req.Body).Decode(dst)
	return
}

func SendJson(rw http.ResponseWriter, data interface{}, status ...int) {
	if len(status) == 0 {
		status = []int{http.StatusOK}
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status[0])
	json.NewEncoder(rw).Encode(data)
	return
}

func GetQueryValue(req *http.Request, key string) string {
	return req.URL.Query().Get(key)
}

func GetHeaderValue(req *http.Request, key string) string {
	return req.Header.Get(key)
}
