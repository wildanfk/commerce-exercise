package liberr_test

import (
	"product-service/internal/util/liberr"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewErrorDetails(t *testing.T) {
	type input struct {
		message string
		code    string
		field   string
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(*liberr.ErrorDetails)
	}{
		{
			name: "Success Create NewErrorDetails",
			in: input{
				message: "Error message",
				code:    "1001",
				field:   "name",
			},
			assertFn: func(result *liberr.ErrorDetails) {
				assert.Equal(t, &liberr.ErrorDetails{
					Message: "Error message",
					Code:    "1001",
					Field:   "name",
				}, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(liberr.NewErrorDetails(tc.in.message, tc.in.code, tc.in.field))
		})
	}
}

func TestErrorDetails_Error(t *testing.T) {
	type input struct {
		message string
		code    string
		field   string
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(string)
	}{
		{
			name: "Success Retrieve Error",
			in: input{
				message: "Error message",
				code:    "1001",
				field:   "name",
			},
			assertFn: func(result string) {
				assert.Equal(t, "Error message", result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(liberr.NewErrorDetails(tc.in.message, tc.in.code, tc.in.field).Error())
		})
	}
}

func TestErrorCodeEquals(t *testing.T) {
	type input struct {
		err  error
		code string
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(bool)
	}{
		{
			name: "Success With ErrorDetail",
			in: input{
				err:  liberr.NewErrorDetails("Error Message", "1001", ""),
				code: "1001",
			},
			assertFn: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "Failed With ErrorDetail",
			in: input{
				err:  liberr.NewErrorDetails("Error Message", "1001", ""),
				code: "1002",
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
		{
			name: "Failed With Unknown Error",
			in: input{
				err:  errors.New("error"),
				code: "1002",
			},
			assertFn: func(result bool) {
				assert.False(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(liberr.ErrorCodeEquals(tc.in.err, tc.in.code))
		})
	}
}
