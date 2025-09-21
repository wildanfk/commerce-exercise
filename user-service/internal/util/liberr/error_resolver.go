package liberr

func ResolveError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*BaseError); ok {
		return err
	}

	if _, ok := err.(*ErrorTracer); ok {
		return err
	}

	if detail, ok := err.(*ErrorDetails); ok {
		berr := NewBaseError(detail)

		// Remove this method (ResolveError) from the stacktrace
		berr.stackTrace = berr.stackTrace[1:]

		return berr
	}

	return NewTracerFromError(err)
}
