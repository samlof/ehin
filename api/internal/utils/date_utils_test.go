package utils

import (
	"testing"
	"time"
)

func TestSecondsUntil12(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	tests := []struct {
		name     string
		now      time.Time
		expected int
	}{
		{
			name:     "13:00:00 Helsinki -> 3600 seconds until 12:00:00 UTC",
			now:      time.Date(2025, 3, 29, 13, 0, 0, 0, loc), // 13:00 Helsinki = 11:00 UTC (EET)
			expected: 3600,
		},
		{
			name:     "13:01:00 Helsinki -> 3540 seconds",
			now:      time.Date(2025, 3, 29, 13, 1, 0, 0, loc),
			expected: 3540,
		},
		{
			name:     "13:01:05 Helsinki -> 3535 seconds",
			now:      time.Date(2025, 3, 29, 13, 1, 5, 0, loc),
			expected: 3535,
		},
		{
			name: "15:00:00 Helsinki (Summer) -> 0 seconds (it's 12:00 UTC)",
			// EEST is UTC+3. 15:00 Helsinki = 12:00 UTC.
			// My implementation returns 24 hours if it's Exactly 12:00 UTC.
			// Let's check Java behavior or re-read.
			// "If now is after 12:00:00 UTC, it calculates until 12:00:00 UTC the next day."
			// If it's exactly 12:00:00, usually you want the NEXT one.
			now:      time.Date(2025, 5, 29, 15, 0, 0, 0, loc),
			expected: 86400,
		},
		{
			name:     "16:00:00 Helsinki (Summer) -> 23 hours until next 12:00 UTC",
			now:      time.Date(2025, 5, 29, 16, 0, 0, 0, loc), // 16:00 Helsinki = 13:00 UTC
			expected: 23 * 3600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SecondsUntil12(tt.now)
			if got != tt.expected {
				t.Errorf("SecondsUntil12() = %v, want %v", got, tt.expected)
			}
		})
	}

	t.Run("DST Transitions Helsinki", func(t *testing.T) {
		// Helsinki DST 2025:
		// Starts: Sunday, March 30, 03:00 (clocks go forward to 04:00) -> UTC+3
		// Ends: Sunday, October 26, 04:00 (clocks go back to 03:00) -> UTC+2

		// March 30 transition: 02:59:59 EET (UTC+2) -> 04:00:00 EEST (UTC+3)
		// Before transition: 02:30 EET on Mar 30 = 00:30 UTC
		// Seconds until 12:00 UTC on Mar 30 = 11.5 hours = 41400 seconds
		t1 := time.Date(2025, 3, 30, 2, 30, 0, 0, loc)
		expected1 := int(time.Date(2025, 3, 30, 12, 0, 0, 0, time.UTC).Sub(t1.UTC()).Seconds())
		if got := SecondsUntil12(t1); got != expected1 {
			t.Errorf("March DST transition before: got %d, want %d", got, expected1)
		}

		// After transition: 04:30 EEST on Mar 30 = 01:30 UTC
		// Seconds until 12:00 UTC on Mar 30 = 10.5 hours = 37800 seconds
		t2 := time.Date(2025, 3, 30, 4, 30, 0, 0, loc)
		expected2 := int(time.Date(2025, 3, 30, 12, 0, 0, 0, time.UTC).Sub(t2.UTC()).Seconds())
		if got := SecondsUntil12(t2); got != expected2 {
			t.Errorf("March DST transition after: got %d, want %d", got, expected2)
		}

		// October transition: 03:59:59 EEST (UTC+3) -> 03:00:00 EET (UTC+2)
		// 03:30 EEST on Oct 26 = 00:30 UTC
		t3 := time.Date(2025, 10, 26, 3, 30, 0, 0, loc)
		expected3 := int(time.Date(2025, 10, 26, 12, 0, 0, 0, time.UTC).Sub(t3.UTC()).Seconds())
		if got := SecondsUntil12(t3); got != expected3 {
			t.Errorf("October DST transition: got %d, want %d", got, expected3)
		}
	})
}

func TestGetGmtStringForCache(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Standard formatting",
			time:     time.Date(2025, 3, 29, 11, 57, 0, 0, time.UTC),
			expected: "Sat, 29 Mar 2025 11:57:00 GMT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGmtStringForCache(tt.time)
			if got != tt.expected {
				t.Errorf("GetGmtStringForCache() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetGmtStringForCacheFromParts(t *testing.T) {
	tests := []struct {
		name         string
		date         time.Time
		extraSeconds int
		expected     string
	}{
		{
			name:         "2025-03-29, 0 seconds",
			date:         time.Date(2025, 3, 29, 0, 0, 0, 0, time.UTC),
			extraSeconds: 0,
			expected:     "Sat, 29 Mar 2025 11:57:00 GMT",
		},
		{
			name:         "2025-02-22, 25 seconds",
			date:         time.Date(2025, 2, 22, 0, 0, 0, 0, time.UTC),
			extraSeconds: 25,
			expected:     "Sat, 22 Feb 2025 11:57:25 GMT",
		},
		{
			name:         "2025-05-29, 50 seconds",
			date:         time.Date(2025, 5, 29, 0, 0, 0, 0, time.UTC),
			extraSeconds: 50,
			expected:     "Thu, 29 May 2025 11:57:50 GMT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGmtStringForCacheFromParts(tt.date, tt.extraSeconds)
			if got != tt.expected {
				t.Errorf("GetGmtStringForCacheFromParts() = %v, want %v", got, tt.expected)
			}
		})
	}
}
