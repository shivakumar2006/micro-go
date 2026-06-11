package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidateErrors(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required filed", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid email", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at least %s characters", err.Field(), err.Param()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at most %s characters", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
