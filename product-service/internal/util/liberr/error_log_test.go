package liberr_test

import (
	"product-service/internal/util/liberr"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAppendErrorLogField(t *testing.T) {
	type input struct {
		fields []zap.Field
		err    error
	}

	testCases := []struct {
		name     string
		in       input
		assertFn func([]zap.Field)
	}{
		{
			name: "Success AppendErrorLogField With Base Error",
			in: input{
				fields: []zap.Field{zap.String("fieldKey", "fieldValue")},
				err:    liberr.NewBaseError(liberr.NewErrorDetails("Error happend", "ERROR_CODE", "field")),
			},
			assertFn: func(result []zap.Field) {
				mapFields := map[string]any{
					"fieldKey":   true,
					"error":      true,
					"stacktrace": true,
				}

				for _, field := range result {
					_, ok := mapFields[field.Key]
					assert.True(t, ok)
				}
			},
		},
		{
			name: "Success AppendErrorLogField With Any Error",
			in: input{
				fields: []zap.Field{zap.String("fieldKey", "fieldValue")},
				err:    errors.New("error happened"),
			},
			assertFn: func(result []zap.Field) {
				mapFields := map[string]any{
					"fieldKey":   true,
					"error":      true,
					"stacktrace": true,
				}

				for _, field := range result {
					_, ok := mapFields[field.Key]
					assert.True(t, ok)
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn(liberr.AppendErrorLogField(tc.in.fields, tc.in.err))
		})
	}
}
