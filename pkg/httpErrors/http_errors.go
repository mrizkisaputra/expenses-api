package httpErrors

import (
	"fmt"
	"net/http"
)

const (
	EmailAlreadyExistsMsg     = "User with given email already exists"
	InvalidEmailOrPasswordMsg = "Invalid email or password"
)

const (
	BadRequestErrorMsg     = "Bad request"
	UnauthorizedErrorMsg   = "Unauthorized"
	ForbiddenErrorMsg      = "Forbidden"
	NotFoundErrorMsg       = "Not found"
	StatusConflictErrorMsg = "Status conflict"
	InternalServerErrorMsg = "Internal server error"
)

// costum error http request
type Error struct {
	Status  int
	Message string
	Causes  interface{} // tidak untuk di expose ke response, hanya untuk logger
}

func (e Error) Error() string {
	return fmt.Sprintf("status: %d - message: %s", e.Status, e.Message)
}

// create a new error inctance with an optional message
func NewError(status int, message string, causes interface{}) *Error {
	return &Error{
		Status:  status,
		Message: message,
		Causes:  causes,
	}
}

// create a new Error Bad Request instance
func NewBadRequestError(causes interface{}) *Error {
	return &Error{
		Status:  http.StatusBadRequest,
		Message: BadRequestErrorMsg,
		Causes:  causes,
	}
}

// create a new Error Unauthorized instance
func NewUnauthorizedError(causes interface{}) *Error {
	return &Error{
		Status:  http.StatusUnauthorized,
		Message: UnauthorizedErrorMsg,
		Causes:  causes,
	}
}

// create a new Error Notfound instance
func NewNotFoundError(causes interface{}) *Error {
	return &Error{
		Status:  http.StatusNotFound,
		Message: NotFoundErrorMsg,
		Causes:  causes,
	}
}

// create a new Error Internal Server instance
func NewInternalServerError(causes interface{}) *Error {
	return &Error{
		Status:  http.StatusInternalServerError,
		Message: InternalServerErrorMsg,
		Causes:  causes,
	}
}
