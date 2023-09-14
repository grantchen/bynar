package errors

import (
	"errors"
	"fmt"
	"runtime"
)

const (
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
)

// Error is a status error.
type Error struct {
	Reason   string
	Message  string
	Metadata map[string]string

	internal bool // is internal error or not
	cause    error
	stack    []byte // runtime stack
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) FullError() string {
	return fmt.Sprintf("error: reason = %s message = %s metadata = %v internal = %v cause = %v",
		e.Reason, e.Message, e.Metadata, e.internal, e.cause)
}

func (e *Error) Stack() string {
	return string(e.stack)
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Reason == e.Reason
	}
	return false
}

// WithInternal with the underlying internal of the error.
func (e *Error) WithInternal() *Error {
	err := Clone(e)
	err.internal = true
	return err
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	// overwrite stack
	err.withStack()
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := Clone(e)
	err.Metadata = md
	return err
}

// IsInternal return internal
func (e *Error) IsInternal() bool {
	if e == nil {
		return false
	}
	return e.internal
}

// with the runtime stack.
func (e *Error) withStack() {
	buf := make([]byte, 64<<10)
	buf = buf[:runtime.Stack(buf, false)]
	e.stack = buf
}

// New returns an error object for the reason, message.
func New(reason, message string) *Error {
	err := &Error{
		Reason:  reason,
		Message: message,
	}
	err.withStack()
	return err
}

// NewUnknownError returns an error object for the unknown reason, message.
func NewUnknownError(message string) *Error {
	return New(UnknownReason, message)
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(reason, format string, a ...interface{}) *Error {
	return New(reason, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(reason, format string, a ...interface{}) error {
	return New(reason, fmt.Sprintf(format, a...))
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// Clone deep clone error to a new error.
func Clone(err *Error) *Error {
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &Error{
		Reason:   err.Reason,
		Message:  err.Message,
		Metadata: metadata,
		internal: err.internal,
		cause:    err.cause,
		stack:    err.stack,
	}
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	return nil
}
