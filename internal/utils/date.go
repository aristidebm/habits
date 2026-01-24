package utils

import (
	"time"
)

func IsToday(date time.Time) bool {
	today := time.Now()
	return date.Year() == today.Year() &&
		date.Month() == today.Month() &&
		date.Day() == today.Day()
}
