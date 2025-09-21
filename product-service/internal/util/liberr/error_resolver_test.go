package liberr_test

import (
	"fmt"
	"product-service/internal/util/liberr"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestResolveError(t *testing.T) {
	errDetail := liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "field")

	type input struct {
		err error
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(err error)
	}{
		{
			name: "Success Resolve Error From Empty Error",
			in:   input{},
			assertFn: func(result error) {
				assert.Nil(t, result)
			},
		},
		{
			name: "Success Resolve Error From BaseError",
			in: input{
				err: liberr.NewBaseError(errDetail),
			},
			assertFn: func(result error) {
				berr, ok := result.(*liberr.BaseError)

				assert.True(t, ok)

				stacktrace := berr.StackTrace()

				functionsName := make([]string, 0)
				sourceLocation := make([]string, 0)
				for _, stack := range stacktrace {
					functionsName = append(functionsName, fmt.Sprintf("%n", stack))
					sourceLocation = append(sourceLocation, fmt.Sprintf("%s:%d", stack, stack))
				}

				expectedFunctionNames := []string{
					"TestResolveError",
				}
				for i, expectedfuncName := range expectedFunctionNames {
					assert.Equal(t, expectedfuncName, functionsName[i])
				}

				expectedSourceLocations := []string{
					"error_resolver_test.go:34",
				}
				for i, expectedSourceLoc := range expectedSourceLocations {
					assert.Equal(t, expectedSourceLoc, sourceLocation[i])
				}
			},
		},
		{
			name: "Success Resolve Error From ErrorTracer",
			in: input{
				err: liberr.NewTracerFromError(errors.New("Error Happened")),
			},
			assertFn: func(result error) {
				et, ok := result.(*liberr.ErrorTracer)

				assert.True(t, ok)

				stacktrace := et.StackTrace()

				functionsName := make([]string, 0)
				sourceLocation := make([]string, 0)
				for _, stack := range stacktrace {
					functionsName = append(functionsName, fmt.Sprintf("%n", stack))
					sourceLocation = append(sourceLocation, fmt.Sprintf("%s:%d", stack, stack))
				}

				expectedFunctionNames := []string{
					"TestResolveError",
				}
				for i, expectedfuncName := range expectedFunctionNames {
					assert.Equal(t, expectedfuncName, functionsName[i])
				}

				expectedSourceLocations := []string{
					"error_resolver_test.go:68",
				}
				for i, expectedSourceLoc := range expectedSourceLocations {
					assert.Equal(t, expectedSourceLoc, sourceLocation[i])
				}
			},
		},
		{
			name: "Success Resolve Error From Any Error",
			in: input{
				err: errors.New("Error Happened"),
			},
			assertFn: func(result error) {
				et, ok := result.(*liberr.ErrorTracer)

				assert.True(t, ok)

				stacktrace := et.StackTrace()

				functionsName := make([]string, 0)
				sourceLocation := make([]string, 0)
				for _, stack := range stacktrace {
					functionsName = append(functionsName, fmt.Sprintf("%n", stack))
					sourceLocation = append(sourceLocation, fmt.Sprintf("%s:%d", stack, stack))
				}

				expectedFunctionNames := []string{
					"TestResolveError",
				}
				for i, expectedfuncName := range expectedFunctionNames {
					assert.Equal(t, expectedfuncName, functionsName[i])
				}

				expectedSourceLocations := []string{
					"error_resolver_test.go:102",
				}
				for i, expectedSourceLoc := range expectedSourceLocations {
					assert.Equal(t, expectedSourceLoc, sourceLocation[i])
				}
			},
		},
		{
			name: "Success Resolve Error From ErrorDetail",
			in: input{
				err: errDetail,
			},
			assertFn: func(result error) {
				berr, ok := result.(*liberr.BaseError)

				assert.True(t, ok)

				stacktrace := berr.StackTrace()

				functionsName := make([]string, 0)
				sourceLocation := make([]string, 0)
				for _, stack := range stacktrace {
					functionsName = append(functionsName, fmt.Sprintf("%n", stack))
					sourceLocation = append(sourceLocation, fmt.Sprintf("%s:%d", stack, stack))
				}

				expectedFunctionNames := []string{
					"TestResolveError.func6",
				}
				for i, expectedfuncName := range expectedFunctionNames {
					assert.Equal(t, expectedfuncName, functionsName[i])
				}

				expectedSourceLocations := []string{
					"error_resolver_test.go:170",
				}
				for i, expectedSourceLoc := range expectedSourceLocations {
					assert.Equal(t, expectedSourceLoc, sourceLocation[i])
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(liberr.ResolveError(tc.in.err))
		})
	}
}
