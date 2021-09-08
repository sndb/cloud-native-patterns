package patterns

import (
	"context"
	"errors"
	"sync"
	"time"
)

func Breaker(c Circuit, threshold int) Circuit {
	remaining := threshold
	var last time.Time
	var mu sync.RWMutex

	return func(ctx context.Context) (string, error) {
		mu.RLock()
		if remaining <= 0 {
			if time.Now().Before(last.Add(time.Second * 2 << (-remaining))) {
				mu.RUnlock()
				return "", errors.New("open circuit")
			}
		}
		mu.RUnlock()

		r, err := c(ctx)

		mu.Lock()
		defer mu.Unlock()

		last = time.Now()
		if err != nil {
			remaining--
		} else {
			remaining = threshold
		}
		return r, err
	}
}
