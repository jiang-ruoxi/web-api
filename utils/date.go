package utils

import "time"

func Validate(date string) bool {
	if _, err := time.ParseInLocation("2006/01/02", date, time.Local); err != nil {
		if _, err := time.ParseInLocation("2006/01", date, time.Local); err != nil {
			return false
		}
	}

	return true
}
