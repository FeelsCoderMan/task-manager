package httpError

import (
	"errors"
	"net/http"
)

type HttpError struct {
	error
	Code int
}

func newHttpError(code int, errorMessage string) *HttpError {
	return &HttpError{
		Code:  code,
		error: errors.New(errorMessage),
	}
}

func NotFound(errorMessage string) *HttpError {
	return newHttpError(http.StatusNotFound, errorMessage)
}

func BadRequest(errorMessage string) *HttpError {
	return newHttpError(http.StatusBadRequest, errorMessage)
}

func InternalServerError(errorMessage string) *HttpError {
	return newHttpError(http.StatusInternalServerError, errorMessage)
}
