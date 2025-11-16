package faulttolerance

import (
	"fmt"
	"time"
)

func Timeout(operation func() error, timeout time.Duration) error {
	done := make(chan error, 1)

	go func() {
		err := operation()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("faulttolerance.Timeout: the function failed with an error: %w", err)
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("faulttolerance.Timeout: The function did not respond in time %v", timeout)
	}
}
