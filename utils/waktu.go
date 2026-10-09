package utils

import (
	"fmt"
	"time"
)

func ParseJam(value string, loc *time.Location) (time.Time, error) {
	formats := []string{
		"15:04:05",
		"15:04",
	}

	for _, format := range formats {
		t, err := time.ParseInLocation(format, value, loc)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf(
		"format jam harus HH:MM atau HH:MM:SS",
	)
}
