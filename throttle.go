package patterns

import (
	"context"
	"errors"
	"sync"
	"time"
)

func Throttle(e Effector, max uint, d time.Duration) (wrapped Effector, stop func()) {
	tokens := max
	var mu sync.Mutex

	stopc := make(chan struct{})
	ticker := time.NewTicker(d)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				if tokens < max {
					tokens++
				}
				mu.Unlock()
			case <-stopc:
				return
			}
		}
	}()

	return func(ctx context.Context) (string, error) {
		mu.Lock()
		defer mu.Unlock()

		if tokens > 0 {
			tokens--
		} else {
			return "", errors.New("request throttled")
		}

		return e(ctx)
	}, func() { close(stopc) }
}
