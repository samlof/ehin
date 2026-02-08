package service

import (
	"testing"
	"time"
)

func TestDateService_Now(t *testing.T) {
	service := NewDateService()

	now := service.Now()
	systemNow := time.Now()

	// Check if the time returned is close to the system time (within 1 second)
	diff := systemNow.Sub(now)
	if diff < 0 {
		diff = -diff
	}

	if diff > time.Second {
		t.Errorf("Expected time to be close to system time, got difference of %v", diff)
	}
}

// TestMockTimeProvider demonstrates how the TimeProvider interface can be mocked.
type mockTimeProvider struct {
	mockTime time.Time
}

func (m *mockTimeProvider) Now() time.Time {
	return m.mockTime
}

func TestMocking(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	mock := &mockTimeProvider{mockTime: fixedTime}

	if mock.Now() != fixedTime {
		t.Errorf("Expected %v, got %v", fixedTime, mock.Now())
	}
}
