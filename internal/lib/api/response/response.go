package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Response базовая структура ответа
// @Description Базовый ответ API
type Response struct {
	Status string `json:"status" example:"OK"`
	Error  string `json:"error,omitempty" example:"error message"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// Error создает ответ с ошибкой
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// OK создает успешный ответ
func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

// ValidationError создает ответ с ошибками валидации
func ValidationError(errs validator.ValidationErrors) Response {
	return ValidErrors(errs)
}

// ValidErrors создает ответ с ошибками валидации
func ValidErrors(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at least %s characters", err.Field(), err.Param()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at most %s characters", err.Field(), err.Param()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be a valid email", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
