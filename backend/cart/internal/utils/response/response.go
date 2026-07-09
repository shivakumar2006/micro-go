package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, res any) error {
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

func ValidateErrors(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.Tag() {

		case "required":
			errMsg = append(errMsg,
				fmt.Sprintf("%s is required", err.Field()))

		case "gt":
			errMsg = append(errMsg,
				fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param()))

		case "gte":
			errMsg = append(errMsg,
				fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param()))

		case "lt":
			errMsg = append(errMsg,
				fmt.Sprintf("%s must be less than %s", err.Field(), err.Param()))

		case "lte":
			errMsg = append(errMsg,
				fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param()))

		default:
			errMsg = append(errMsg,
				fmt.Sprintf("%s is invalid", err.Field()))
		}
	}

	return Response{
		Status:  StatusError,
		Message: "validation failed",
		Error:   strings.Join(errMsg, ", "),
	}
}
