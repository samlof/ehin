package utils

import (
	"time"
)

// SecondsUntil12 calculates the number of seconds until 12:00:00 UTC.
// If the current time is after 12:00:00 UTC, it calculates until 12:00:00 UTC the next day.
func SecondsUntil12(now time.Time) int {
	utcNow := now.UTC()
	target := time.Date(utcNow.Year(), utcNow.Month(), utcNow.Day(), 12, 0, 0, 0, time.UTC)

	if utcNow.After(target) || utcNow.Equal(target) {
		target = target.AddDate(0, 0, 1)
	}

	return int(target.Sub(utcNow).Seconds())
}

// RFC1123GMT is the format for GMT times in HTTP headers.
const RFC1123GMT = "Mon, 02 Jan 2006 15:04:05 GMT"

// GetGmtStringForCache formats a time.Time into an RFC 1123 formatted string with GMT suffix.
// Example: "Sat, 29 Mar 2025 11:57:00 GMT"
func GetGmtStringForCache(t time.Time) string {
	return t.UTC().Format(RFC1123GMT)
}

// GetGmtStringForCacheFromParts creates an RFC 1123 formatted string for a specific date
// at 11:57:00 UTC plus the specified extra seconds.
func GetGmtStringForCacheFromParts(date time.Time, extraSeconds int) string {
	t := time.Date(date.Year(), date.Month(), date.Day(), 11, 57, extraSeconds, 0, time.UTC)
	return GetGmtStringForCache(t)
}
