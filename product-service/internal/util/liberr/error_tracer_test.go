package liberr_test

import (
	"fmt"
	"product-service/internal/util/liberr"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestNewTracer(t *testing.T) {
	type input struct {
		message string
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(*liberr.ErrorTracer)
	}{
		{
			name: "Success Create NewTracer",
			in: input{
				message: "Error message",
			},
			assertFn: func(result *liberr.ErrorTracer) {
				assert.Equal(t, &liberr.ErrorTracer{
					Message: "Error message",
				}, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(liberr.NewTracer(tc.in.message))
		})
	}
}

func TestErrorTracer_StackTrace(t *testing.T) {
	type input struct {
		err *liberr.ErrorTracer
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(*liberr.ErrorTracer)
	}{
		{
			name: "Success Retrieve StackTrace From NewTracerFromError",
			in: input{
				err: liberr.NewTracerFromError(errors.New("Error Happened")),
			},
			assertFn: func(in *liberr.ErrorTracer) {
				stacktrace := in.StackTrace()

				functionsName := make([]string, 0)
				sourceLocation := make([]string, 0)
				for _, stack := range stacktrace {
					functionsName = append(functionsName, fmt.Sprintf("%n", stack))
					sourceLocation = append(sourceLocation, fmt.Sprintf("%s:%d", stack, stack))
				}

				expectedFunctionNames := []string{
					"TestErrorTracer_StackTrace",
				}
				for i, expectedfuncName := range expectedFunctionNames {
					assert.Equal(t, expectedfuncName, functionsName[i])
				}

				expectedSourceLocations := []string{
					"error_tracer_test.go:55",
				}
				for i, expectedSourceLoc := range expectedSourceLocations {
					assert.Equal(t, expectedSourceLoc, sourceLocation[i])
				}
			},
		},
		{
			name: "Success Retrieve StackTrace From Wrap",
			in: input{
				err: func() *liberr.ErrorTracer {
					initialErr := errors.New("Error Happened")
					middleErr := liberr.NewTracer("Error On InitialErr").Wrap(initialErr)

					return liberr.NewTracer("Error On MiddleErr").Wrap(middleErr)
				}(),
			},
			assertFn: func(in *liberr.ErrorTracer) {
				stacktrace := in.StackTrace()

				functionsName := make([]string, 0)
				sourceLocation := make([]string, 0)
				for _, stack := range stacktrace {
					functionsName = append(functionsName, fmt.Sprintf("%n", stack))
					sourceLocation = append(sourceLocation, fmt.Sprintf("%s:%d", stack, stack))
				}

				fmt.Println(functionsName)
				fmt.Println(sourceLocation)

				expectedFunctionNames := []string{
					"TestErrorTracer_StackTrace.func2",
					"TestErrorTracer_StackTrace",
				}
				for i, expectedfuncName := range expectedFunctionNames {
					assert.Equal(t, expectedfuncName, functionsName[i])
				}

				expectedSourceLocations := []string{
					"error_tracer_test.go:87",
					"error_tracer_test.go:90",
				}
				for i, expectedSourceLoc := range expectedSourceLocations {
					assert.Equal(t, expectedSourceLoc, sourceLocation[i])
				}
			},
		},
		{
			name: "Success Retrieve StackTrace From NewTracer",
			in: input{
				err: liberr.NewTracer("Error Happened"),
			},
			assertFn: func(in *liberr.ErrorTracer) {
				assert.Nil(t, in.StackTrace())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.err)
		})
	}
}

func TestErrorTracer_Error(t *testing.T) {
	type input struct {
		err *liberr.ErrorTracer
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(string)
	}{
		{
			name: "Success Retrieve Error From NewTracerFromError",
			in: input{
				err: liberr.NewTracerFromError(errors.New("Error Happened")),
			},
			assertFn: func(result string) {
				assert.Equal(t, "Error Happened", result)
			},
		},
		{
			name: "Success Retrieve Error From Wrap",
			in: input{
				err: func() *liberr.ErrorTracer {
					initialErr := errors.New("Error Happened")
					middleErr := liberr.NewTracer("Error On InitialErr").Wrap(initialErr)

					return liberr.NewTracer("Error On MiddleErr").Wrap(middleErr)
				}(),
			},
			assertFn: func(result string) {
				assert.Equal(t, "Error On MiddleErr", result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.err.Error())
		})
	}
}

func TestErrorTracer_Unwrap(t *testing.T) {
	initialErr := errors.New("Error Happened")

	type input struct {
		err *liberr.ErrorTracer
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(error)
	}{
		{
			name: "Success Unwrap Error",
			in: input{
				err: &liberr.ErrorTracer{
					Err: initialErr,
				},
			},
			assertFn: func(result error) {
				assert.ErrorIs(t, initialErr, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.err.Unwrap())
		})
	}
}
