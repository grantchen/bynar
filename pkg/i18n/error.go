package i18n

import (
	"errors"
	"fmt"
	"runtime"
)

var _ error = (*Error)(nil)

// Error is a i18n error.
type Error struct {
	Language string // language of the error message
	Message  string // translated error message

	messageId string // i18n message id
	cause     error  // underlying cause of the error
	stack     []byte // runtime stack
}

// NewError creates a new error with language and message.
func NewError(language, message string) *Error {
	err := &Error{
		Language: language,
		Message:  message,
	}
	err.withStack()
	return err
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

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// FullError returns the error message with the cause.
func (e *Error) FullError() string {
	return fmt.Sprintf("error: message = %scause = %v",
		e.Message, e.cause)
}

// Stack returns the runtime stack.
func (e *Error) Stack() string {
	return string(e.stack)
}

// with the runtime stack.
func (e *Error) withStack() {
	buf := make([]byte, 64<<10)
	buf = buf[:runtime.Stack(buf, false)]
	e.stack = buf
}

// withMessageId with the message id.
func (e *Error) withMessageId(messageId string) {
	e.messageId = messageId
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) {
	e.cause = cause
	// overwrite stack
	e.withStack()
}
