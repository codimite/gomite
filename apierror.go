package gomite

import (
	"encoding/json"
	"net/http"
)

type ErrorVal struct {
	Code    int
	Message string
}

type ApiError struct {
	ErrorMap map[string]ErrorVal
	Debug    bool
}

type errorResponse struct {
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Exception string `json:"exception,omitempty"`
}

func (apiErr ApiError) SendError(rw http.ResponseWriter, err error, status int, exceptions ...string) {
	var exception string
	if len(exceptions) > 0 {
		exception = exceptions[0]
	}
	errVal, found := apiErr.ErrorMap[err.Error()]
	if !found {
		errVal = ErrorVal{Code: 100, Message: "Unexpected Error"}
		exception = err.Error()
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	if apiErr.Debug {
		json.NewEncoder(rw).Encode(errorResponse{
			Message:   errVal.Message,
			Code:      errVal.Code,
			Exception: exception,
		})
	} else {
		json.NewEncoder(rw).Encode(errorResponse{
			Message: errVal.Message,
			Code:    errVal.Code,
		})
	}
}
