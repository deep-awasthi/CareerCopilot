package errors

import (
	"errors"
	"net/http"
)

// Domain error types
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Constructor helpers
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// Common error factories
func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, message)
}

func Internal(message string) *AppError {
	return New(http.StatusInternalServerError, message)
}

func InternalWrap(message string, err error) *AppError {
	return Wrap(http.StatusInternalServerError, message, err)
}

func TooManyRequests(message string) *AppError {
	return New(http.StatusTooManyRequests, message)
}

// Standard error variables
var (
	ErrUserNotFound        = NotFound("user not found")
	ErrUserAlreadyExists   = Conflict("user already exists")
	ErrInvalidCredentials  = Unauthorized("invalid credentials")
	ErrTokenExpired        = Unauthorized("token expired")
	ErrTokenInvalid        = Unauthorized("invalid token")
	ErrJobNotFound         = NotFound("job not found")
	ErrApplicationNotFound = NotFound("application not found")
	ErrInterviewNotFound   = NotFound("interview not found")
	ErrReferralNotFound    = NotFound("referral not found")
	ErrCompanyNotFound     = NotFound("company not found")
	ErrBookmarkNotFound    = NotFound("bookmark not found")
	ErrAlertNotFound       = NotFound("keyword alert not found")
)

// IsNotFound checks if error is a not found error
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == http.StatusNotFound
	}
	return false
}
