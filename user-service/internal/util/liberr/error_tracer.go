package liberr

import (
	"github.com/pkg/errors"
)

// StackTracer is interface for retrieve stacktrace from pkg/errors
type StackTracer interface {
	StackTrace() errors.StackTrace
}

type ErrorTracer struct {
	Message    string
	Err        error
	stackTrace errors.StackTrace
}

func NewTracer(message string) *ErrorTracer {
	return &ErrorTracer{
		Message: message,
	}
}

func NewTracerFromError(err error) *ErrorTracer {
	tracer := NewTracer(err.Error())
	tracer.Err = err
	tracer.assignStackTrace()
	return tracer
}

func (e *ErrorTracer) Error() string {
	return e.Message
}

func (e *ErrorTracer) Wrap(err error) *ErrorTracer {
	e.Err = err
	e.assignStackTrace()
	return e
}

func (e *ErrorTracer) Unwrap() error {
	return e.Err
}

func (e *ErrorTracer) StackTrace() errors.StackTrace {
	return e.stackTrace
}

func (e *ErrorTracer) assignStackTrace() {
	// Check if the error that wrapped is already implement StackTracer
	if stracer, ok := e.Err.(StackTracer); ok {
		e.stackTrace = stracer.StackTrace()
	} else {
		// Use pkg/errors to wrapped error with stacktrace
		e.Err = errors.WithStack(e.Err)
		// Assert the e.Err to StacktTracer so we can call the StackTrace() function
		if stracer, ok = e.Err.(StackTracer); ok {
			// exclude the first 2 stacktrace
			// because we want to remove both this method (assignStackTrace) and the caller (`NewTracerFromError` or `Wrap`)
			e.stackTrace = stracer.StackTrace()[2:]
		}
	}
}
