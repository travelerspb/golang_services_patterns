package stability

import (
	"cloud_patters/stability/src"
	"context"
	"errors"
	"sync"
	"time"
)

func Breaker(circuit src.Circuit, failureThreshold uint) src.Circuit {
	var consecutiveFailures = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock() // Establish a "read lock"
		d := consecutiveFailures - int(failureThreshold)
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}
		m.RUnlock() // Release read lock

		response, err := circuit(ctx) // Issue request proper
		m.Lock()                      // Lock around shared resources
		defer m.Unlock()
		lastAttempt = time.Now() // Record time of attempt
		if err != nil {          // Circuit returned an error,
			consecutiveFailures++ // so we count the failure
			return response, err  // and return
		}
		consecutiveFailures = 0 // Reset failures counter

		return response, nil
	}
}
