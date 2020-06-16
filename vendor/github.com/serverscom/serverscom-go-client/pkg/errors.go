package serverscom

import (
	"fmt"
)

type responseErrorWrapper struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

// ParsingError represents any error related to the problem with parsing a body
type ParsingError struct {
	StatusCode   int
	Body         string
	ParsingError error
}

func newParsingError(statusCode int, body string, err error) error {
	return &ParsingError{
		StatusCode:   statusCode,
		Body:         body,
		ParsingError: err,
	}
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf(
		"Parsing error: %s, for body: %s, with status code: %d",
		e.ParsingError.Error(),
		e.Body,
		e.StatusCode,
	)
}

// BadRequestError represents an errors related to 400 response status code
type BadRequestError struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

func newBadRequestError(statusCode int, errorCode, message string) error {
	return &BadRequestError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Bad request: %s", e.Message)
}

// UnauthorizedError represents an errors related to 401 response status code
type UnauthorizedError struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

func newUnauthorizedError(statusCode int, errorCode, message string) error {
	return &UnauthorizedError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("Unauthorized: %s", e.Message)
}

// ForbiddenError represents an errors related to 403 response status code
type ForbiddenError struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

func newForbiddenError(statusCode int, errorCode, message string) error {
	return &ForbiddenError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("Forbidden: %s", e.Message)
}

// NotFoundError represents an errors related to 404 response status code
type NotFoundError struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

func newNotFoundError(statusCode int, errorCode, message string) error {
	return &NotFoundError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not found: %s", e.Message)
}

// ConflictError represents an errors related to 409 response status code
type ConflictError struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

func newConflictError(statusCode int, errorCode, message string) error {
	return &ConflictError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("Conflict: %s", e.Message)
}

// UnprocessableEntityError represents an errors related to 422 response status code
type UnprocessableEntityError struct {
	StatusCode int
	ErrorCode  string
	Message    string
	Errors     map[string]string
}

func newUnprocessableEntityError(statusCode int, errorCode, message string, errors map[string]string) error {
	return &UnprocessableEntityError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
		Errors:     errors,
	}
}

func (e *UnprocessableEntityError) Error() string {
	return fmt.Sprintf("Unprocessable entity: %s, with errors: %v", e.Message, e.Errors)
}

// InternalServerError represents an errors related to 500 response status code
type InternalServerError struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

func newInternalServerError(statusCode int, errorCode, message string) error {
	return &InternalServerError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("Internal server error: %s", e.Message)
}
