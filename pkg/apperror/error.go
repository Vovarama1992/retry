package apperror

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(message string, code int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NotFound(message string) *AppError {
	return New(message, http.StatusNotFound)
}

func BadRequest(message string) *AppError {
	return New(message, http.StatusBadRequest)
}

func Internal(message string) *AppError {
	return New(message, http.StatusInternalServerError)
}
