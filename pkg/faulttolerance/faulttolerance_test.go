package faulttolerance

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	tests := []struct {
		name       string
		operation  func() error
		maxRetries int
		baseDelay  time.Duration
		expected   error
	}{
		{
			name: "succeeds on first attempt",
			operation: func() error {
				return nil
			},
			maxRetries: 3,
			baseDelay:  10 * time.Millisecond,
			expected:   nil,
		},
		{
			name: "succeeds after second attempt",
			operation: func() error {
				staticAttemptCount := 0
				return func() error {
					staticAttemptCount++
					if staticAttemptCount == 1 {
						return errors.New("temporary failure")
					}
					return nil
				}()
			},
			maxRetries: 3,
			baseDelay:  10 * time.Millisecond,
			expected:   nil,
		},
		{
			name: "fails after max retries",
			operation: func() error {
				return errors.New("permanent failure")
			},
			maxRetries: 3,
			baseDelay:  10 * time.Millisecond,
			expected:   errors.New("faulttolerance.Retry: operation failed after 3 retries"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "succeeds after second attempt" {
				attemptCount := 0
				operation := func() error {
					attemptCount++
					if attemptCount == 1 {
						return errors.New("temporary failure")
					}
					return nil
				}
				err := Retry(operation, tt.maxRetries, tt.baseDelay)
				if (err != nil) != (tt.expected != nil) {
					t.Errorf("Retry() error = %v, expected = %v", err, tt.expected)
				}
				return
			}

			err := Retry(tt.operation, tt.maxRetries, tt.baseDelay)
			if (err != nil) != (tt.expected != nil) {
				t.Errorf("Retry() error = %v, expected = %v", err, tt.expected)
			}
		})
	}
}

func TestTimeout(t *testing.T) {
	tests := []struct {
		name      string
		operation func() error
		timeout   time.Duration
		expected  error
	}{
		{
			name: "succeeds within timeout",
			operation: func() error {
				return nil
			},
			timeout:  100 * time.Millisecond,
			expected: nil,
		},
		{
			name: "fails with error within timeout",
			operation: func() error {
				return errors.New("operation failed")
			},
			timeout:  100 * time.Millisecond,
			expected: errors.New("faulttolerance.Timeout: the function failed with an error: operation failed"),
		},
		{
			name: "times out",
			operation: func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			},
			timeout:  10 * time.Millisecond,
			expected: errors.New("faulttolerance.Timeout: The function did not respond in time 10ms"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Timeout(tt.operation, tt.timeout)
			if (err != nil) != (tt.expected != nil) {
				t.Errorf("Timeout() error = %v, expected = %v", err, tt.expected)
			}
		})
	}
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestDeadLetterQueue(t *testing.T) {
	tests := []struct {
		name     string
		messages []string
		handler  func(string) error
		expected []string
	}{
		{
			name:     "no failures",
			messages: []string{"msg1", "msg2", "msg3"},
			handler: func(msg string) error {
				return nil
			},
			expected: []string{},
		},
		{
			name:     "some failures",
			messages: []string{"msg1", "msg2", "msg3", "msg4"},
			handler: func(msg string) error {
				if msg == "msg2" || msg == "msg4" {
					return errors.New("processing failed")
				}
				return nil
			},
			expected: []string{"msg2", "msg4"},
		},
		{
			name:     "all failures",
			messages: []string{"msg1", "msg2", "msg3"},
			handler: func(msg string) error {
				return errors.New("processing failed")
			},
			expected: []string{"msg1", "msg2", "msg3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dlq := NewDeadLetterQueue()
			ProcessWithDLQ(tt.messages, tt.handler, dlq)
			result := dlq.GetMessages()
			if !stringsEqual(result, tt.expected) {
				t.Errorf("ProcessWithDLQ() got = %v, want %v", result, tt.expected)
			}
		})
	}
}