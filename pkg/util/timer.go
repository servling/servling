package util

import "time"

// Timer holds the starting time of an operation.
type Timer struct {
	startTime time.Time
}

// NewTimer creates and returns a new Timer instance.
// Use this when you need to measure multiple, independent operations.
func NewTimer() *Timer {
	return &Timer{}
}

// Start records the current time as the beginning of an operation.
func (t *Timer) Start() {
	t.startTime = time.Now()
}

// Stop returns the duration elapsed since Start was called.
// If Start was never called, it returns zero.
func (t *Timer) Stop() time.Duration {
	if t.startTime.IsZero() {
		return 0
	}
	return time.Since(t.startTime)
}

// --- Global Default Timer ---

// defaultTimer provides a convenient global instance for measuring a single,
// primary operation without needing to create a Timer instance.
var defaultTimer = NewTimer()

// StartTimer begins the global default timer.
func StartTimer() {
	defaultTimer.Start()
}

// StopTimer returns the duration elapsed since the global Start was called.
func StopTimer() time.Duration {
	return defaultTimer.Stop()
}
