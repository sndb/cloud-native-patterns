package patterns

import (
	"context"
	"fmt"
	"time"
)

func Retry(e Effector, max int) Effector {
	return func(ctx context.Context) (string, error) {
		var r string
		var err error
		for i := 0; i < max; i++ {
			r, err = e(ctx)
			if err != nil {
				select {
				case <-time.After(time.Second << i):
				case <-ctx.Done():
					return "", ctx.Err()
				}
			} else {
				return r, err
			}
		}
		return r, fmt.Errorf("maximum number of retries (%d) exceeded: %w", max, err)
	}
}
