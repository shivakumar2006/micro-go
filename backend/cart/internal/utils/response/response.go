package response

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, res *Response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res)
}

func GeneralError(err error) *Response {
	return &Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}
