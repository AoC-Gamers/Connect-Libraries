package errors

import (
	"fmt"
	"net/http"
)

// InternalError represents an error from an internal API call
type InternalError struct {
	StatusCode int
	Service    string
	Endpoint   string
	Message    string
	Details    map[string]interface{}
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("%s %s: %d %s", e.Service, e.Endpoint, e.StatusCode, e.Message)
}

// NewInternalError creates a new internal error
func NewInternalError(statusCode int, service, endpoint, message string) *InternalError {
	return &InternalError{
		StatusCode: statusCode,
		Service:    service,
		Endpoint:   endpoint,
		Message:    message,
		Details:    make(map[string]interface{}),
	}
}

// Common error constructors
func NewBadRequest(service, endpoint, message string) *InternalError {
	return NewInternalError(http.StatusBadRequest, service, endpoint, message)
}

func NewUnauthorized(service, endpoint, message string) *InternalError {
	return NewInternalError(http.StatusUnauthorized, service, endpoint, message)
}

func NewForbidden(service, endpoint, message string) *InternalError {
	return NewInternalError(http.StatusForbidden, service, endpoint, message)
}

func NewNotFound(service, endpoint, message string) *InternalError {
	return NewInternalError(http.StatusNotFound, service, endpoint, message)
}

func NewInternalServerError(service, endpoint, message string) *InternalError {
	return NewInternalError(http.StatusInternalServerError, service, endpoint, message)
}

func NewServiceUnavailable(service, endpoint, message string) *InternalError {
	return NewInternalError(http.StatusServiceUnavailable, service, endpoint, message)
}

// IsInternalError checks if an error is an InternalError
func IsInternalError(err error) bool {
	_, ok := err.(*InternalError)
	return ok
}

// GetStatusCode extracts status code from error, returns 500 if not InternalError
func GetStatusCode(err error) int {
	if ierr, ok := err.(*InternalError); ok {
		return ierr.StatusCode
	}
	return http.StatusInternalServerError
}
