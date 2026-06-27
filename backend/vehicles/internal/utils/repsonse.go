package utils

import "github.com/go-playground/validator/v10"

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func ValidateErrors(errs validator.ValidationErrors) Response {
	var msgs []string

	for _, err := range errs {
		switch err.Tag() {

		}
	}
}
