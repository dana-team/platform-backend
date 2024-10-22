package customerrors

import (
	"errors"
	"fmt"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

// ErrorWithStatusCode is a common interface for errors that include a status code.
type ErrorWithStatusCode interface {
	Error() string
	StatusCode() int
	StatusReason() metav1.StatusReason
}

// APIError is a generic error type that includes an HTTP status code.
type APIError struct {
	Message string
	Code    int
	Reason  metav1.StatusReason
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) StatusCode() int {
	return e.Code
}

func (e *APIError) StatusReason() metav1.StatusReason {
	return e.Reason
}

// NewAPIError creates a new APIError based on a k8s StatusError or defaults to InternalServerError.
func NewAPIError(message string, err error) ErrorWithStatusCode {
	errMessage := fmt.Sprintf("%s, %v", message, err)

	var errWithStatusCode ErrorWithStatusCode
	if errors.As(err, &errWithStatusCode) {
		return errWithStatusCode
	}

	var statusErr *k8serrors.StatusError
	if errors.As(err, &statusErr) {
		return &APIError{
			Message: errMessage,
			Code:    int(statusErr.ErrStatus.Code),
			Reason:  statusErr.ErrStatus.Reason,
		}
	}

	return NewInternalServerError(errMessage)
}

// ValidationError represents an error due to invalid input.
type ValidationError struct {
	Message string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) StatusCode() int {
	return http.StatusBadRequest
}

func (e *ValidationError) StatusReason() metav1.StatusReason {
	return metav1.StatusReasonBadRequest
}

// NotFoundError represents an error when a resource is not found.
type NotFoundError struct {
	Message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

func (e *NotFoundError) StatusReason() metav1.StatusReason {
	return metav1.StatusReasonNotFound
}

// UnauthorizedError represents an error when access is denied.
type UnauthorizedError struct {
	Message string
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		Message: message,
	}
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}

func (e *UnauthorizedError) StatusReason() metav1.StatusReason {
	return metav1.StatusReasonUnauthorized
}

// ConflictError represents an error due to a conflict in state.
type ConflictError struct {
	Message string
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{
		Message: message,
	}
}

func (e *ConflictError) Error() string {
	return e.Message
}

func (e *ConflictError) StatusCode() int {
	return http.StatusConflict
}

func (e *ConflictError) StatusReason() metav1.StatusReason {
	return metav1.StatusReasonConflict
}

// InternalServerError represents an error due to an unknown error.
type InternalServerError struct {
	Message string
}

func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{
		Message: message,
	}
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func (e *InternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}

func (e *InternalServerError) StatusReason() metav1.StatusReason {
	return metav1.StatusReasonInternalError
}
