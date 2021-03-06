package error

import (
	"fmt"
)

// ErrInvalidContent is default invalid-content-error retured
// when function or receiver receives invalid content in args or context.
var ErrInvalidContent = invalidContentError{}

type invalidContentError struct {
	message string
	inner   error
}

func (e invalidContentError) Error() string {
	return fmt.Sprintf("InvalidContentError: %v\n  %v", e.message, e.inner)
}

func (e invalidContentError) Unwrap() error {
	return e.inner
}

func (invalidContentError) Is(target error) bool {
	_, ok := target.(invalidContentError)
	return ok
}

// NewInvalidContentError generates a invalid-content-error
func NewInvalidContentError(inner error, message string) invalidContentError {
	return invalidContentError{
		message: message,
		inner:   inner,
	}
}

// ErrNotFound is default not-found-error retured
// when function or receiver can not find item in repogitory.
var ErrNotFound = notFoundError{}

type notFoundError struct {
	message string
	inner   error
}

func (e notFoundError) Error() string {
	return fmt.Sprintf("NotFoundError: %v\n  %v", e.message, e.inner)
}

func (e notFoundError) Unwrap() error {
	return e.inner
}

func (notFoundError) Is(target error) bool {
	_, ok := target.(notFoundError)
	return ok
}

// NewNotFoundError generates a not-found-error
func NewNotFoundError(inner error, message string) notFoundError {
	return notFoundError{
		message: message,
		inner:   inner,
	}
}

// ErrDuplication is default duplication-error retured
// when same item alraedy exists.
var ErrDuplication = duplicationError{}

type duplicationError struct {
	message string
	inner   error
}

func (e duplicationError) Error() string {
	return fmt.Sprintf("DuplicationError: %v\n  %v", e.message, e.inner)
}

func (e duplicationError) Unwrap() error {
	return e.inner
}

func (duplicationError) Is(target error) bool {
	_, ok := target.(duplicationError)
	return ok
}

// NewDuplicationError generates a duplication-error
func NewDuplicationError(inner error, message string) duplicationError {
	return duplicationError{
		message: message,
		inner:   inner,
	}
}

// ErrInternal is default duplication-error retured
// when internal server error occued.
var ErrInternal = internalError{}

type internalError struct {
	message string
	inner   error
}

func (e internalError) Error() string {
	return fmt.Sprintf("InternalError: %v\n  %v", e.message, e.inner)
}

func (e internalError) Unwrap() error {
	return e.inner
}

func (internalError) Is(target error) bool {
	_, ok := target.(internalError)
	return ok
}

// NewInternalError generates a duplication-error
func NewInternalError(inner error, message string) internalError {
	return internalError{
		message: message,
		inner:   inner,
	}
}

// ErrAuthorization is default authorization-error retured
// when authorization failed.
var ErrAuthorization = authorizationError{}

type authorizationError struct {
	message string
	inner   error
}

func (e authorizationError) Error() string {
	return fmt.Sprintf("AuthorizationError: %v\n  %v", e.message, e.inner)
}

func (e authorizationError) Unwrap() error {
	return e.inner
}

func (authorizationError) Is(target error) bool {
	_, ok := target.(authorizationError)
	return ok
}

// NewAuthorizationError generates a duplication-error
func NewAuthorizationError(inner error, message string) authorizationError {
	return authorizationError{
		message: message,
		inner:   inner,
	}
}
