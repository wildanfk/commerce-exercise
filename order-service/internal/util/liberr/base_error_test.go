package liberr_test

import (
	"fmt"
	"order-service/internal/util/liberr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseError_GetDetails(t *testing.T) {
	errDetail1 := liberr.NewErrorDetails("Error Happened On Request 1", "ERROR_1", "request1")
	errDetail2 := liberr.NewErrorDetails("Error Happened On Request 2", "ERROR_2", "request2")

	type input struct {
		err *liberr.BaseError
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func([]*liberr.ErrorDetails)
	}{
		{
			name: "Success Retrieve Error Detail",
			in: input{
				err: liberr.NewBaseError(errDetail1),
			},
			assertFn: func(result []*liberr.ErrorDetails) {
				expected := []*liberr.ErrorDetails{
					errDetail1,
				}
				assert.ElementsMatch(t, expected, result)
			},
		},
		{
			name: "Success Retrieve Multi Error Detail",
			in: input{
				err: liberr.NewBaseError(errDetail1, errDetail2),
			},
			assertFn: func(result []*liberr.ErrorDetails) {
				expected := []*liberr.ErrorDetails{
					errDetail1,
					errDetail2,
				}
				assert.ElementsMatch(t, expected, result)
			},
		},
		{
			name: "Success Retrieve Error Detail With Addition",
			in: input{
				err: func() *liberr.BaseError {
					err := liberr.NewBaseError(errDetail1)
					err.AddErrorDetails(errDetail2)

					return err
				}(),
			},
			assertFn: func(result []*liberr.ErrorDetails) {
				expected := []*liberr.ErrorDetails{
					errDetail1,
					errDetail2,
				}
				assert.ElementsMatch(t, expected, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.err.GetDetails())
		})
	}
}

func TestBaseError_StackTrace(t *testing.T) {
	type input struct {
		err *liberr.BaseError
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(*liberr.BaseError)
	}{
		{
			name: "Success Retrieve Stacktrace",
			in: input{
				err: liberr.NewBaseError(liberr.NewErrorDetails("Error Happened On Request", "ERROR", "request")),
			},
			assertFn: func(in *liberr.BaseError) {
				stacktrace := in.StackTrace()

				functionsName := make([]string, 0)
				sourceLocation := make([]string, 0)
				for _, stack := range stacktrace {
					functionsName = append(functionsName, fmt.Sprintf("%n", stack))
					sourceLocation = append(sourceLocation, fmt.Sprintf("%s:%d", stack, stack))
				}

				expectedFunctionNames := []string{
					"TestBaseError_StackTrace",
				}
				for i, expectedfuncName := range expectedFunctionNames {
					assert.Equal(t, expectedfuncName, functionsName[i])
				}

				expectedSourceLocations := []string{
					"base_error_test.go:88",
				}
				for i, expectedSourceLoc := range expectedSourceLocations {
					assert.Equal(t, expectedSourceLoc, sourceLocation[i])
				}
			},
		},
		{
			name: "Success Retrieve Stacktrace From Multi Error",
			in: input{
				err: func() *liberr.BaseError {
					errDetail1 := liberr.NewErrorDetails("Error Happened On Request 1", "ERROR_1", "request1")
					errDetail2 := liberr.NewErrorDetails("Error Happened On Request 2", "ERROR_2", "request2")
					errDetail3 := liberr.NewErrorDetails("Error Happened On Request 3", "ERROR_3", "request3")

					return liberr.NewBaseError(errDetail1, errDetail2, errDetail3)
				}(),
			},
			assertFn: func(in *liberr.BaseError) {
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
					"TestBaseError_StackTrace.func2",
					"TestBaseError_StackTrace",
				}
				for i, expectedfuncName := range expectedFunctionNames {
					assert.Equal(t, expectedfuncName, functionsName[i])
				}

				expectedSourceLocations := []string{
					"base_error_test.go:123",
					"base_error_test.go:124",
				}
				for i, expectedSourceLoc := range expectedSourceLocations {
					assert.Equal(t, expectedSourceLoc, sourceLocation[i])
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.err)
		})
	}
}

func TestBaseError_Error(t *testing.T) {
	errDetail1 := liberr.NewErrorDetails("Error Happened On Request 1", "ERROR_1", "request1")
	errDetail2 := liberr.NewErrorDetails("Error Happened On Request 2", "ERROR_2", "request2")

	type input struct {
		err *liberr.BaseError
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(string)
	}{
		{
			name: "Success Retrieve Error",
			in: input{
				err: liberr.NewBaseError(errDetail1),
			},
			assertFn: func(result string) {
				expected := "Error on\ncode: ERROR_1; error: Error Happened On Request 1; field: request1"
				assert.Equal(t, expected, result)
			},
		},
		{
			name: "Success Retrieve Multi Error",
			in: input{
				err: liberr.NewBaseError(errDetail1, errDetail2),
			},
			assertFn: func(result string) {
				expected := "Error on\ncode: ERROR_1; error: Error Happened On Request 1; field: request1\ncode: ERROR_2; error: Error Happened On Request 2; field: request2"
				assert.Equal(t, expected, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.err.Error())
		})
	}
}

func TestBaseError_IsAllCodeEqual(t *testing.T) {
	type input struct {
		err  *liberr.BaseError
		code string
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(bool)
	}{
		{
			name: "Success IsAllCodeEqual",
			in: input{
				err:  liberr.NewBaseError(liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "")),
				code: "ERROR_CODE",
			},
			assertFn: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "Success IsAllCodeEqual On Multiple Error Detail",
			in: input{
				err: liberr.NewBaseError(
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE", ""),
				),
				code: "ERROR_CODE",
			},
			assertFn: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "Failed IsAllCodeEqual",
			in: input{
				err:  liberr.NewBaseError(liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "")),
				code: "WRONG_CODE",
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
		{
			name: "Failed IsAllCodeEqual On Multiple Error Detail",
			in: input{
				err: liberr.NewBaseError(
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE1", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE2", ""),
				),
				code: "ERROR_CODE1",
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
		{
			name: "Failed IsAllCodeEqual On Empty Error Detail",
			in: input{
				err:  &liberr.BaseError{},
				code: "WRONG_CODE",
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.err.IsAllCodeEqual(tc.in.code))
		})
	}
}

func TestBaseError_IsAllExpectedCodes(t *testing.T) {
	type input struct {
		err   *liberr.BaseError
		codes []string
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(bool)
	}{
		{
			name: "Success IsAllExpectedCodes",
			in: input{
				err:   liberr.NewBaseError(liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "")),
				codes: []string{"ERROR_CODE"},
			},
			assertFn: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "Success IsAllExpectedCodes On Multiple Error Detail",
			in: input{
				err: liberr.NewBaseError(
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE1", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE2", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE1", ""),
				),
				codes: []string{"ERROR_CODE1", "ERROR_CODE2"},
			},
			assertFn: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "Failed IsAllExpectedCodes",
			in: input{
				err:   liberr.NewBaseError(liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "")),
				codes: []string{"WRONG_CODE"},
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
		{
			name: "Failed IsAllExpectedCodes On Multiple Error Detail 1",
			in: input{
				err: liberr.NewBaseError(
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE1", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE2", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE3", ""),
				),
				codes: []string{"ERROR_CODE1", "ERROR_CODE2"},
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
		{
			name: "Failed IsAllExpectedCodes On Multiple Error Detail 2",
			in: input{
				err: liberr.NewBaseError(
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE1", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE2", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE4", ""),
				),
				codes: []string{"ERROR_CODE1", "ERROR_CODE2"},
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
		{
			name: "Failed IsAllExpectedCodes On Empty Error Detail",
			in: input{
				err:   &liberr.BaseError{},
				codes: []string{"ERROR_CODE1"},
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.err.IsAllExpectedCodes(tc.in.codes...))
		})
	}
}

func TestBaseError_IsAnyCodeEqual(t *testing.T) {
	type input struct {
		err  *liberr.BaseError
		code string
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(bool)
	}{
		{
			name: "Success IsAnyCodeEqual",
			in: input{
				err:  liberr.NewBaseError(liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "")),
				code: "ERROR_CODE",
			},
			assertFn: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "Success IsAnyCodeEqual On Multiple Error Detail",
			in: input{
				err: liberr.NewBaseError(
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE1", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE2", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE", ""),
				),
				code: "ERROR_CODE",
			},
			assertFn: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "Failed IsAnyCodeEqual",
			in: input{
				err:  liberr.NewBaseError(liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "")),
				code: "WRONG_CODE",
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
		{
			name: "Failed IsAnyCodeEqual On Empty Error Detail",
			in: input{
				err:  &liberr.BaseError{},
				code: "WRONG_CODE",
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(tc.in.err.IsAnyCodeEqual(tc.in.code))
		})
	}
}

func TestBaseError_ReplaceFieldWithMapFieldByCode(t *testing.T) {
	type input struct {
		err            *liberr.BaseError
		mapFieldByCode map[string]string
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func([]*liberr.ErrorDetails)
	}{
		{
			name: "Success ReplaceFieldWithMapFieldByCode",
			in: input{
				err: liberr.NewBaseError(liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "field")),
				mapFieldByCode: map[string]string{
					"ERROR_CODE": "field_update",
				},
			},
			assertFn: func(result []*liberr.ErrorDetails) {
				expected := []*liberr.ErrorDetails{
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "field_update"),
				}

				assert.ElementsMatch(t, expected, result)
			},
		},
		{
			name: "Success ReplaceFieldWhenCodeEqual On Multiple Error Detail",
			in: input{
				err: liberr.NewBaseError(
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE1", "field1"),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE", ""),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE2", "field2"),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE", ""),
				),
				mapFieldByCode: map[string]string{
					"ERROR_CODE1": "new_field1",
					"ERROR_CODE":  "new_field",
				},
			},
			assertFn: func(result []*liberr.ErrorDetails) {
				expected := []*liberr.ErrorDetails{
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE1", "new_field1"),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "new_field"),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE2", "field2"),
					liberr.NewErrorDetails("Error Happened", "ERROR_CODE", "new_field"),
				}

				assert.ElementsMatch(t, expected, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.in.err.ReplaceFieldWithMapFieldByCode(tc.in.mapFieldByCode)
			tc.assertFn(tc.in.err.GetDetails())
		})
	}
}
