package faulttolerance

import (
	"fmt"
	"math"
	"time"
)

func Retry(operation func() error, maxRetries int, baseDelay time.Duration) error {
	var (
		err                         error
		maxDelay                            = 5 * time.Second
		defaultNumberForBackoffTime float64 = 2
	)

	for retry := range maxRetries {
		err = operation()
		if err == nil {
			return nil
		}

		backoffTime := baseDelay * time.Duration(math.Pow(defaultNumberForBackoffTime, float64(retry)))
		if backoffTime > maxDelay {
			backoffTime = maxDelay
		}

		time.Sleep(backoffTime)
	}

	return fmt.Errorf("faulttolerance.Retry: operation failed after %d retries", maxRetries)
}
