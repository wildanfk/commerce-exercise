package liberr

import (
	"fmt"

	"go.uber.org/zap"
)

func AppendErrorLogField(fields []zap.Field, err error) []zap.Field {
	if err != nil {
		fields = append(fields, zap.String("error", err.Error()))

		if stracer, ok := err.(StackTracer); ok {
			stack := stracer.StackTrace()

			// Format the stack trace
			var formattedStack string
			for _, f := range stack {
				formattedStack += fmt.Sprintf("%+v\n", f)
			}
			fields = append(fields, zap.String("stacktrace", formattedStack))
		}
	}

	return fields
}
