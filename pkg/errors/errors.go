package errors

import (
	"errors"
	"fmt"
	"runtime"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

const (
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
)

// Error is a status error.
// code is for international
type Error struct {
	Reason   string
	Message  string
	Metadata map[string]string
	Code     string

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

// WithInternalCause with the underlying internal and cause of the error.
func (e *Error) WithInternalCause(cause error) *Error {
	err := Clone(e)
	err.internal = true
	err.cause = cause
	// overwrite stack
	err.withStack()
	err.setKnownErrorCode()
	return err
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
	err.setKnownErrorCode()
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
func New(reason, message, code string) *Error {
	err := &Error{
		Reason:  reason,
		Message: message,
		Code:    code,
	}
	err.withStack()
	return err
}

// NewUnknownError returns an error object for the unknown reason, message.
func NewUnknownError(message, code string) *Error {
	return New(UnknownReason, message, code)
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
		Code:     err.Code,
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

// set code value
func (e *Error) setKnownErrorCode() {
	if e.Code != "" || e.cause == nil {
		return
	}
	if errors.Is(e.cause, gip.ErrUserNotFound) {
		e.Code = ErrCodeNoUserFound
	} else if errors.Is(e.cause, gip.ErrPhoneNumberAlreadyExists) {
		e.Code = ErrCodeEmailAlreadyExists
	} else if errors.Is(e.cause, gip.ErrPhoneNumberAlreadyExists) {
		e.Code = ErrCodePhoneNumberAlreadyExists
	} else if errors.Is(e.cause, gip.ErrIDTokenInvalid) {
		e.Code = ErrCodeIDTokenInvalid
	} else {
		e.Code = ErrCode
	}

}
