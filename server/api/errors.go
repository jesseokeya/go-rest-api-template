package api

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// ApiError holds possible http api error fields
type ApiError struct {
	Err error `json:"-"`

	StatusCode int    `json:"-"`
	StatusText string `json:"status"`

	Location  string      `json:"location,omitempty"`
	AppCode   int64       `json:"code,omitempty"`
	ErrorText string      `json:"error,omitempty"`
	Cause     string      `json:"cause,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// Error return an error text
func (e *ApiError) Error() string {
	return e.ErrorText
}

// Render sends error message to the client
func (e *ApiError) Render(w http.ResponseWriter, r *http.Request) error {
	pc := make([]uintptr, 5) // maximum 5 levels to go
	runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc)
	next := false
	for {
		frame, more := frames.Next()
		if next {
			e.Location = fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}
		if strings.Contains(frame.File, "api/renderer.go") {
			next = true
		}
		if !more {
			break
		}
	}
	w.WriteHeader(e.StatusCode)
	return nil
}

var (
	// ErrBadID error message for bad or invalid id
	ErrBadID = &ApiError{StatusCode: http.StatusBadRequest, ErrorText: "bad or invalid id"}
	// ErrPermissionDenied error message for permission denied
	ErrPermissionDenied = &ApiError{StatusCode: http.StatusUnauthorized, ErrorText: "permission denied"}
	// ErrInvalidSession error message for invalid session
	ErrInvalidSession = &ApiError{StatusCode: http.StatusUnauthorized, ErrorText: "invalid session"}
	// ErrEncryptionError error message for bcrypt password encryption error
	ErrEncryptionError = &ApiError{StatusCode: http.StatusInternalServerError, ErrorText: "issue with email signup"}
	// ErrUserExists throws error if user already exists
	ErrUserExists = &ApiError{StatusCode: http.StatusConflict, ErrorText: "user already exists"}
)

// ErrUnauthorized is error message for Unauthorized
func ErrUnauthorized(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusUnauthorized,
		StatusText: "Unauthorized",
		ErrorText:  err.Error(),
	}
}

// ErrDatabase is error message for inability to perform database operation
func ErrDatabase(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		StatusText: "Couldn't perform operation, please try again later.",
		ErrorText:  err.Error(),
	}
}

// ErrPermission is error message for Unauthorized
func ErrPermission(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusUnauthorized,
		StatusText: "Permission denied.",
		ErrorText:  err.Error(),
	}
}

// ErrInvalidRequest is error message for Unauthorized
func ErrInvalidRequest(err error, data ...interface{}) *ApiError {
	v := &ApiError{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid request.",
		ErrorText:  err.Error(),
	}
	if len(data) > 0 {
		v.Data = data[0]
	}
	return v
}

// ErrServiceUnavailable is error message for Service Unavailable
func ErrServiceUnavailable(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusServiceUnavailable,
		StatusText: "Service Unavailable.",
		ErrorText:  err.Error(),
	}
}

// ErrInternalServerError is error message for Internal Server.
func ErrInternalServerError(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		StatusText: "Internal Server.",
		ErrorText:  err.Error(),
	}
}

// ErrRequestEntityTooLarge is error message for Request Entity Too Large
func ErrRequestEntityTooLarge(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusRequestEntityTooLarge,
		StatusText: "Request Entity Too Large",
		ErrorText:  err.Error(),
	}
}

// ErrInvalidEmailSignup shapes error message if the users email address is invalid
func ErrInvalidEmailSignup(cause error) *ApiError {
	return &ApiError{
		Err:        cause,
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid email address",
		ErrorText:  fmt.Sprintf("Invalid email: %s", cause.Error()),
		Data:       map[string]string{"field": "email"},
	}
}

// IgnoreError ignores error
func IgnoreError(v ...interface{}) {}
