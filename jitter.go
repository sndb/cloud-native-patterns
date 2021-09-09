package patterns

import (
	"math/rand"
	"time"
)

func ExponentialBackoffWithJitter(action func() error) {
	err := action()
	base, cap := time.Second, time.Minute

	for backoff := base; err != nil; backoff <<= 1 {
		if backoff > cap {
			backoff = cap
		}

		jitter := rand.Int63n(int64(backoff * 3))
		time.Sleep(base + time.Duration(jitter))

		err = action()
	}
}
