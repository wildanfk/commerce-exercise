package entity

import "user-service/internal/util/libvalidate"

func init() {
	// Register Json Tag on Validator Error Field
	validate := libvalidate.Validator()
	libvalidate.RegisterJSONTagField(validate)
}
