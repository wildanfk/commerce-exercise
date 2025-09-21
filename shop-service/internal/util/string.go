package util

import "strconv"

func ConvertStringToIntWithDefault(value string, defaultValue int) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
