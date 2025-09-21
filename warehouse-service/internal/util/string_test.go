package util_test

import (
	"testing"
	"warehouse-service/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToIntWithDefault(t *testing.T) {
	type input struct {
		value        string
		defaultValue int
	}

	tests := []struct {
		name     string
		in       input
		assertFn func(int)
	}{
		{
			name: "Success Convert String",
			in: input{
				value:        "10",
				defaultValue: 0,
			},
			assertFn: func(result int) {
				assert.Equal(t, 10, result)
			},
		},
		{
			name: "Failed Convert Decimal Value",
			in: input{
				value:        "15.5",
				defaultValue: 10,
			},
			assertFn: func(result int) {
				assert.Equal(t, 10, result)
			},
		},
		{
			name: "Failed Convert Value",
			in: input{
				value:        "1abc",
				defaultValue: 10,
			},
			assertFn: func(result int) {
				assert.Equal(t, 10, result)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(util.ConvertStringToIntWithDefault(tc.in.value, tc.in.defaultValue))
		})
	}
}
