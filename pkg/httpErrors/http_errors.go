package httpErrors

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mrizkisaputra/expenses-api/pkg/contextutils"
	"net/http"
)

const (
	EmailAlreadyExistsMsg      = "User with given email already exists"
	InvalidEmailOrPasswordMsg  = "Invalid email or password"
	InvalidJwtTokenMsg         = "Invalid JWT token"
	MaxFileSizeMsg             = "File size exceeds 1MB"
	NotAllowedImageHeaderMsg   = "Not allowed image header"
	NotAllowedFileExtensionMsg = "Not allowed file extension"
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
	return fmt.Sprintf("status: %d - message: %s - causes: %v", e.Status, e.Message, e.Causes)
}

func (e Error) GetCauses() interface{} {
	return e.Causes
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

func NewInvalidJwtTokenError(causes interface{}) *Error {
	return &Error{
		Status:  http.StatusUnauthorized,
		Message: InvalidJwtTokenMsg,
		Causes:  causes,
	}
}

// ErrorResponse is a function return response error json
func ErrorResponse(ctx *gin.Context, err error) (int, any) {
	apiErrorResponse := ParseErrors(ctx, err)
	return apiErrorResponse.ErrorInfo.Status, apiErrorResponse
}

func ParseErrors(ctx *gin.Context, err error) ApiErrorResponse {
	requestID := contextutils.GetRequestId(ctx)

	var er *Error
	var errValidation validator.ValidationErrors

	switch {
	case errors.As(err, &er):
		{
			return NewApiErrorResponse(ErrorInfo{
				Status:   er.Status,
				Message:  er.Message,
				SubError: nil,
			}, requestID)
		}
	case errors.As(err, &errValidation):
		{
			return NewApiErrorResponse(ErrorInfo{
				Status:   http.StatusBadRequest,
				Message:  BadRequestErrorMsg,
				SubError: validationError(errValidation),
			}, requestID)
		}
	default:
		return NewApiErrorResponse(ErrorInfo{
			Status:   http.StatusInternalServerError,
			Message:  InternalServerErrorMsg,
			SubError: nil,
		}, requestID)
	}
}

func validationError(err validator.ValidationErrors) *[]ApiValidationError {
	var apiValidationErr []ApiValidationError
	fieldTagMessage := map[string]string{
		"required": "REQUIRED",
		"email":    "EMAIL_FORMAT",
		"max":      "TO_LONG",
		"min":      "TO_SHORT",
		"alpha":    "MUST_ALPHA",
		"numeric":  "MUST_NUMERIC",
		"uuid":     "MUST_UUID",
	}
	for _, e := range err {
		if msg, ok := fieldTagMessage[e.Tag()]; ok {
			apiValidationErr = append(apiValidationErr, ApiValidationError{
				Field:         e.Field(),
				RejectedValue: e.Value(),
				Message:       msg,
			})
		}
	}

	return &apiValidationErr
}
