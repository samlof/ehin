package service

import "time"

// TimeProvider defines an interface for providing the current time.
// This allows for easy mocking in unit tests.
type TimeProvider interface {
	Now() time.Time
}

// DateService implements the TimeProvider interface using the standard time package.
type DateService struct{}

// NewDateService creates a new instance of DateService.
func NewDateService() *DateService {
	return &DateService{}
}

// Now returns the current local time.
func (s *DateService) Now() time.Time {
	return time.Now()
}
