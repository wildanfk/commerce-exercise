package entity

import "order-service/internal/util/libvalidate"

func init() {
	// Register Json Tag on Validator Error Field
	validate := libvalidate.Validator()
	libvalidate.RegisterJSONTagField(validate)
}
