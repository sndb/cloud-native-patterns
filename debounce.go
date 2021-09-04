package patterns

import (
	"context"
	"sync"
	"time"
)

type Effector func(context.Context) (string, error)

func Debounce(e Effector, threshold time.Duration) Effector {
	var r string
	var err error
	var last time.Time
	var mu sync.Mutex

	return func(ctx context.Context) (string, error) {
		mu.Lock()
		defer mu.Unlock()
		if last.Add(threshold).Before(time.Now()) {
			r, err = e(ctx)
			last = time.Now()
		}
		return r, err
	}
}
