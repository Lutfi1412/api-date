package utils

func NullIfEmpty(value string) interface{} {
	if value == "" {
		return nil
	}

	return value
}
