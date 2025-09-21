package libvalidate_test

import (
	"errors"
	"testing"
	"user-service/internal/util/liberr"
	"user-service/internal/util/libvalidate"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	testCases := []struct {
		name     string
		assertFn func(*validator.Validate)
	}{
		{
			name: "Success Init Validator",
			assertFn: func(result *validator.Validate) {
				assert.NotNil(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(libvalidate.Validator())
		})
	}
}

func TestResolveError(t *testing.T) {
	libvalidate.RegisterJSONTagField(libvalidate.Validator())

	type leaf2 struct {
		Numvar       int      `json:"numvar" validate:"required,gte=5,lte=10"`
		Stringvar    string   `json:"stringvar" validate:"required,min=5,max=10"`
		Listvar      []string `json:"listvar" validate:"required,min=2,max=4"`
		Stripjsonvar string   `json:"-" validate:"required,min=5,max=10"`
		Nonjsonvar   string   `validate:"required,min=5,max=10"`
	}

	type leaf1 struct {
		Numvar       int      `json:"numvar" validate:"required,gte=5,lte=10"`
		Stringvar    string   `json:"stringvar" validate:"required,min=5,max=10"`
		Listvar      []string `json:"listvar" validate:"required,min=2,max=4"`
		Stripjsonvar string   `json:"-" validate:"required,min=5,max=10"`
		Nonjsonvar   string   `validate:"required,min=5,max=10"`
		Objectleaf2  *leaf2   `json:"objectleaf2" validate:"required"`
	}

	type root struct {
		Numvar       int      `json:"numvar" validate:"required,gte=5,lte=10"`
		Stringvar    string   `json:"stringvar" validate:"required,min=5,max=10"`
		Listvar      []string `json:"listvar" validate:"required,min=2,max=4"`
		Stripjsonvar string   `json:"-" validate:"required,min=5,max=10"`
		Nonjsonvar   string   `validate:"required,min=5,max=10"`
		Objectleaf1  *leaf1   `json:"objectleaf1" validate:"required"`
	}

	type input struct {
		err error
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func(result error)
	}{
		{
			name: "Success Resolve Error With Empty Value Of Root",
			in: input{
				err: libvalidate.Validator().Struct(&root{}),
			},
			assertFn: func(result error) {
				assert.NotNil(t, result)

				berr, ok := result.(*liberr.BaseError)
				assert.True(t, ok)

				errorCode := "ERROR_CODE"
				detailErrs := []*liberr.ErrorDetails{
					liberr.NewErrorDetails("numvar is a required field", errorCode, "numvar"),
					liberr.NewErrorDetails("stringvar is a required field", errorCode, "stringvar"),
					liberr.NewErrorDetails("listvar is a required field", errorCode, "listvar"),
					liberr.NewErrorDetails("Stripjsonvar is a required field", errorCode, "Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar is a required field", errorCode, "Nonjsonvar"),
					liberr.NewErrorDetails("objectleaf1 is a required field", errorCode, "objectleaf1"),
				}
				assert.ElementsMatch(t, detailErrs, berr.GetDetails())
			},
		},
		{
			name: "Success Resolve Error With Empty Value Of Leaf 1",
			in: input{
				err: libvalidate.Validator().Struct(&root{
					Objectleaf1: &leaf1{},
				}),
			},
			assertFn: func(result error) {
				assert.NotNil(t, result)

				berr, ok := result.(*liberr.BaseError)
				assert.True(t, ok)

				errorCode := "ERROR_CODE"
				detailErrs := []*liberr.ErrorDetails{
					liberr.NewErrorDetails("numvar is a required field", errorCode, "numvar"),
					liberr.NewErrorDetails("stringvar is a required field", errorCode, "stringvar"),
					liberr.NewErrorDetails("listvar is a required field", errorCode, "listvar"),
					liberr.NewErrorDetails("Stripjsonvar is a required field", errorCode, "Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar is a required field", errorCode, "Nonjsonvar"),

					liberr.NewErrorDetails("numvar is a required field", errorCode, "objectleaf1.numvar"),
					liberr.NewErrorDetails("stringvar is a required field", errorCode, "objectleaf1.stringvar"),
					liberr.NewErrorDetails("listvar is a required field", errorCode, "objectleaf1.listvar"),
					liberr.NewErrorDetails("Stripjsonvar is a required field", errorCode, "objectleaf1.Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar is a required field", errorCode, "objectleaf1.Nonjsonvar"),
					liberr.NewErrorDetails("objectleaf2 is a required field", errorCode, "objectleaf1.objectleaf2"),
				}
				assert.ElementsMatch(t, detailErrs, berr.GetDetails())
			},
		},
		{
			name: "Success Resolve Error With Empty Value Of Leaf 2",
			in: input{
				err: libvalidate.Validator().Struct(&root{
					Objectleaf1: &leaf1{
						Objectleaf2: &leaf2{},
					},
				}),
			},
			assertFn: func(result error) {
				assert.NotNil(t, result)

				berr, ok := result.(*liberr.BaseError)
				assert.True(t, ok)

				errorCode := "ERROR_CODE"
				detailErrs := []*liberr.ErrorDetails{
					liberr.NewErrorDetails("numvar is a required field", errorCode, "numvar"),
					liberr.NewErrorDetails("stringvar is a required field", errorCode, "stringvar"),
					liberr.NewErrorDetails("listvar is a required field", errorCode, "listvar"),
					liberr.NewErrorDetails("Stripjsonvar is a required field", errorCode, "Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar is a required field", errorCode, "Nonjsonvar"),

					liberr.NewErrorDetails("numvar is a required field", errorCode, "objectleaf1.numvar"),
					liberr.NewErrorDetails("stringvar is a required field", errorCode, "objectleaf1.stringvar"),
					liberr.NewErrorDetails("listvar is a required field", errorCode, "objectleaf1.listvar"),
					liberr.NewErrorDetails("Stripjsonvar is a required field", errorCode, "objectleaf1.Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar is a required field", errorCode, "objectleaf1.Nonjsonvar"),

					liberr.NewErrorDetails("numvar is a required field", errorCode, "objectleaf1.objectleaf2.numvar"),
					liberr.NewErrorDetails("stringvar is a required field", errorCode, "objectleaf1.objectleaf2.stringvar"),
					liberr.NewErrorDetails("listvar is a required field", errorCode, "objectleaf1.objectleaf2.listvar"),
					liberr.NewErrorDetails("Stripjsonvar is a required field", errorCode, "objectleaf1.objectleaf2.Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar is a required field", errorCode, "objectleaf1.objectleaf2.Nonjsonvar"),
				}
				assert.ElementsMatch(t, detailErrs, berr.GetDetails())
			},
		},
		{
			name: "Success Resolve Error With Minimal Value",
			in: input{
				err: libvalidate.Validator().Struct(&root{
					Numvar:       1,
					Stringvar:    "ab",
					Listvar:      []string{"1"},
					Stripjsonvar: "abc",
					Nonjsonvar:   "abcd",
					Objectleaf1: &leaf1{
						Numvar:       1,
						Stringvar:    "ab",
						Listvar:      []string{"1"},
						Stripjsonvar: "abc",
						Nonjsonvar:   "abcd",
						Objectleaf2: &leaf2{
							Numvar:       1,
							Stringvar:    "ab",
							Listvar:      []string{"1"},
							Stripjsonvar: "abc",
							Nonjsonvar:   "abcd",
						},
					},
				}),
			},
			assertFn: func(result error) {
				assert.NotNil(t, result)

				berr, ok := result.(*liberr.BaseError)
				assert.True(t, ok)

				errorCode := "ERROR_CODE"
				detailErrs := []*liberr.ErrorDetails{
					liberr.NewErrorDetails("numvar must be 5 or greater", errorCode, "numvar"),
					liberr.NewErrorDetails("stringvar must be at least 5 characters in length", errorCode, "stringvar"),
					liberr.NewErrorDetails("listvar must contain at least 2 items", errorCode, "listvar"),
					liberr.NewErrorDetails("Stripjsonvar must be at least 5 characters in length", errorCode, "Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar must be at least 5 characters in length", errorCode, "Nonjsonvar"),

					liberr.NewErrorDetails("numvar must be 5 or greater", errorCode, "objectleaf1.numvar"),
					liberr.NewErrorDetails("stringvar must be at least 5 characters in length", errorCode, "objectleaf1.stringvar"),
					liberr.NewErrorDetails("listvar must contain at least 2 items", errorCode, "objectleaf1.listvar"),
					liberr.NewErrorDetails("Stripjsonvar must be at least 5 characters in length", errorCode, "objectleaf1.Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar must be at least 5 characters in length", errorCode, "objectleaf1.Nonjsonvar"),

					liberr.NewErrorDetails("numvar must be 5 or greater", errorCode, "objectleaf1.objectleaf2.numvar"),
					liberr.NewErrorDetails("stringvar must be at least 5 characters in length", errorCode, "objectleaf1.objectleaf2.stringvar"),
					liberr.NewErrorDetails("listvar must contain at least 2 items", errorCode, "objectleaf1.objectleaf2.listvar"),
					liberr.NewErrorDetails("Stripjsonvar must be at least 5 characters in length", errorCode, "objectleaf1.objectleaf2.Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar must be at least 5 characters in length", errorCode, "objectleaf1.objectleaf2.Nonjsonvar"),
				}
				assert.ElementsMatch(t, detailErrs, berr.GetDetails())
			},
		},
		{
			name: "Success Resolve Error With Maximal Value",
			in: input{
				err: libvalidate.Validator().Struct(&root{
					Numvar:       11,
					Stringvar:    "abcdefghijk",
					Listvar:      []string{"1", "2", "3", "4", "5"},
					Stripjsonvar: "abcdefghijkl",
					Nonjsonvar:   "abcdefghijklm",
					Objectleaf1: &leaf1{
						Numvar:       11,
						Stringvar:    "abcdefghijk",
						Listvar:      []string{"1", "2", "3", "4", "5"},
						Stripjsonvar: "abcdefghijkl",
						Nonjsonvar:   "abcdefghijklm",
						Objectleaf2: &leaf2{
							Numvar:       11,
							Stringvar:    "abcdefghijk",
							Listvar:      []string{"1", "2", "3", "4", "5"},
							Stripjsonvar: "abcdefghijkl",
							Nonjsonvar:   "abcdefghijklm",
						},
					},
				}),
			},
			assertFn: func(result error) {
				assert.NotNil(t, result)

				berr, ok := result.(*liberr.BaseError)
				assert.True(t, ok)

				errorCode := "ERROR_CODE"
				detailErrs := []*liberr.ErrorDetails{
					liberr.NewErrorDetails("numvar must be 10 or less", errorCode, "numvar"),
					liberr.NewErrorDetails("stringvar must be a maximum of 10 characters in length", errorCode, "stringvar"),
					liberr.NewErrorDetails("listvar must contain at maximum 4 items", errorCode, "listvar"),
					liberr.NewErrorDetails("Stripjsonvar must be a maximum of 10 characters in length", errorCode, "Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar must be a maximum of 10 characters in length", errorCode, "Nonjsonvar"),

					liberr.NewErrorDetails("numvar must be 10 or less", errorCode, "objectleaf1.numvar"),
					liberr.NewErrorDetails("stringvar must be a maximum of 10 characters in length", errorCode, "objectleaf1.stringvar"),
					liberr.NewErrorDetails("listvar must contain at maximum 4 items", errorCode, "objectleaf1.listvar"),
					liberr.NewErrorDetails("Stripjsonvar must be a maximum of 10 characters in length", errorCode, "objectleaf1.Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar must be a maximum of 10 characters in length", errorCode, "objectleaf1.Nonjsonvar"),

					liberr.NewErrorDetails("numvar must be 10 or less", errorCode, "objectleaf1.objectleaf2.numvar"),
					liberr.NewErrorDetails("stringvar must be a maximum of 10 characters in length", errorCode, "objectleaf1.objectleaf2.stringvar"),
					liberr.NewErrorDetails("listvar must contain at maximum 4 items", errorCode, "objectleaf1.objectleaf2.listvar"),
					liberr.NewErrorDetails("Stripjsonvar must be a maximum of 10 characters in length", errorCode, "objectleaf1.objectleaf2.Stripjsonvar"),
					liberr.NewErrorDetails("Nonjsonvar must be a maximum of 10 characters in length", errorCode, "objectleaf1.objectleaf2.Nonjsonvar"),
				}
				assert.ElementsMatch(t, detailErrs, berr.GetDetails())
			},
		},
		{
			name: "Success Resolve Error Without Error",
			in: input{
				err: libvalidate.Validator().Struct(&root{
					Numvar:       9,
					Stringvar:    "abcdefghij",
					Listvar:      []string{"1", "2", "3", "4"},
					Stripjsonvar: "abcdefg",
					Nonjsonvar:   "abcdefgh",
					Objectleaf1: &leaf1{
						Numvar:       9,
						Stringvar:    "abcdefghij",
						Listvar:      []string{"1", "2", "3", "4"},
						Stripjsonvar: "abcdefg",
						Nonjsonvar:   "abcdefgh",
						Objectleaf2: &leaf2{
							Numvar:       9,
							Stringvar:    "abcdefghij",
							Listvar:      []string{"1", "2", "3", "4"},
							Stripjsonvar: "abcdefg",
							Nonjsonvar:   "abcdefgh",
						},
					},
				}),
			},
			assertFn: func(result error) {
				assert.Nil(t, result)
			},
		},
		{
			name: "Success Resolve Error With Any Error",
			in: input{
				err: errors.New("Error Happened"),
			},
			assertFn: func(result error) {
				assert.NotNil(t, result)

				_, ok := result.(*liberr.ErrorTracer)
				assert.True(t, ok)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(libvalidate.ResolveError(tc.in.err, "ERROR_CODE"))
		})
	}
}
