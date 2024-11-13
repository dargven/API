package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

//Ответчик оставляем я думаю его сносить ненадо

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func ValidErrors(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required filed \n", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL \n", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid \n", err.Field()))
		}

	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
