package liberr

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"
)

type BaseError struct {
	details    []*ErrorDetails
	stackTrace errors.StackTrace
}

func NewBaseError(details ...*ErrorDetails) *BaseError {
	berr := &BaseError{details: details}

	// initiate pkg/errors to build stacktrace
	err := errors.New("")
	stracer, ok := err.(StackTracer)
	if ok {
		// exclude the first 1 stacktrace
		// because we want to remove this method (NewBaseError) from the stacktrace
		berr.stackTrace = stracer.StackTrace()[1:]
	}
	return berr
}

func (b *BaseError) StackTrace() errors.StackTrace {
	return b.stackTrace
}

func (b *BaseError) AddErrorDetails(errors ...*ErrorDetails) {
	b.details = append(b.details, errors...)
}

func (b *BaseError) GetDetails() []*ErrorDetails {
	return b.details
}

func (b *BaseError) Error() string {
	buff := bytes.NewBufferString("")

	buff.WriteString("Error on\n")
	for _, err := range b.details {
		buff.WriteString("code: ")
		buff.WriteString(err.Code)
		buff.WriteString("; error: ")
		buff.WriteString(err.Error())
		buff.WriteString("; field: ")
		buff.WriteString(err.Field)
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func (b *BaseError) IsAllCodeEqual(code string) bool {
	if len(b.details) == 0 {
		return false
	}

	for _, d := range b.GetDetails() {
		if d.Code != code {
			return false
		}
	}
	return true
}

func (b *BaseError) IsAllExpectedCodes(codes ...string) bool {
	if len(b.details) == 0 {
		return false
	}

	expectedCodes := map[string]bool{}
	for _, code := range codes {
		expectedCodes[code] = true
	}

	for _, d := range b.GetDetails() {
		if !expectedCodes[d.Code] {
			return false
		}
	}
	return true
}

func (b *BaseError) IsAnyCodeEqual(code string) bool {
	for _, d := range b.GetDetails() {
		if d.Code == code {
			return true
		}
	}
	return false
}

func (b *BaseError) ReplaceFieldWithMapFieldByCode(mapFieldByCode map[string]string) {
	for i, d := range b.GetDetails() {
		if field, ok := mapFieldByCode[d.Code]; ok {
			b.details[i].Field = field
		}
	}
}
