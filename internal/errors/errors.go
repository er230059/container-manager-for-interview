package errors

import (
	"errors"
	"net/http"
)

type CustomError struct {
	Message string
	Status  int
	Cause   error
}

func newCustomError(status int, message string) *CustomError {
	return &CustomError{
		Status:  status,
		Message: message,
	}
}

func (e CustomError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Message
}

func (e CustomError) New(message string) error {
	e.Cause = errors.New(message)
	return &e
}

func (e CustomError) Wrap(err error) error {
	if err == nil {
		return nil
	}
	e.Cause = err
	return &e
}

func (e CustomError) Is(err error) bool {
	var ae CustomError
	if errors.As(err, &ae) {
		return e.Message == ae.Message
	}
	return e.Error() == err.Error()
}

var (
	BadRequest                 = newCustomError(http.StatusBadRequest, "bad request")
	Unauthorized               = newCustomError(http.StatusUnauthorized, "unauthorized")
	PermissionDenied           = newCustomError(http.StatusForbidden, "permission denied")
	EmptyPassword              = newCustomError(http.StatusBadRequest, "password cannot be empty")
	UserNotFound               = newCustomError(http.StatusNotFound, "user not found")
	JobNotFound                = newCustomError(http.StatusNotFound, "job not found")
	ContainerNotFound          = newCustomError(http.StatusNotFound, "container not found")
	ConflictContainerOperation = newCustomError(http.StatusConflict, "conflict container operation")
	InternalServerError        = newCustomError(http.StatusInternalServerError, "internal server error")
)
