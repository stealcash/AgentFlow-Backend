package exceptions

import (
	"fmt"
	"net/http"
)

// HttpException represents a custom HTTP error.
type HttpException struct {
	StatusCode int
	Message    string
}

// Error makes it compatible with the error interface.
func (e HttpException) Error() string {
	return e.Message
}

// New creates a new HttpException.
func New(statusCode int, message string) HttpException {
	return HttpException{
		StatusCode: statusCode,
		Message:    message,
	}
}

func BadRequest(message string) {
	panic(New(http.StatusBadRequest, message))
}

func Unauthorized(message string) {
	panic(New(http.StatusUnauthorized, message))
}

func Forbidden(message string) {
	panic(New(http.StatusForbidden, message))
}

func NotFound(message string) {
	panic(New(http.StatusNotFound, message))
}

func Conflict(message string) {
	panic(New(http.StatusConflict, message))
}

func Internal(message string) {
	panic(New(http.StatusInternalServerError, message))
}

func UnprocessableEntity(message string) {
	panic(New(http.StatusUnprocessableEntity, message))
}

// Custom lets you panic with any custom status code.
func Custom(statusCode int, message string) {
	panic(New(statusCode, message))
}

type HttpExceptionWithLog struct {
	HttpException
	LogMessage string
}

func CriticalWithLog(logMessage, message string) {
	panic(HttpExceptionWithLog{
		HttpException: New(http.StatusInternalServerError, message),
		LogMessage:    logMessage,
	})
}

func AnyError(err error) error {
	return fmt.Errorf("wrapped error: %w", err)
}
