package faulttolerance

import (
	"fmt"
	"math"
	"time"
)

func Retry(operation func() error, maxRetries int, baseDelay time.Duration) error {
	var (
		err error
		maxDelay time.Duration = 5 * time.Second
	)

	for retry := 0; retry < maxRetries; retry++ {
		err = operation()
		if err == nil {
			return nil
		}

		backoffTime := baseDelay * time.Duration(math.Pow(2, float64(retry)))
		if backoffTime > maxDelay {
			backoffTime = maxDelay
		}
		
		time.Sleep(backoffTime)
	}

	return fmt.Errorf("faulttolerance.Retry: operation failed after %d retries", maxRetries)
}