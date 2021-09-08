package patterns

import (
	"context"
	"errors"
	"testing"
)

var errIntentional = errors.New("intentional fail!")

func failNTimes(n int) Effector {
	var count int

	return func(ctx context.Context) (string, error) {
		count++
		if count <= n {
			return "", errIntentional
		}
		return "success", nil
	}
}

func TestRetry(t *testing.T) {
	f1 := Retry(failNTimes(2), 3)
	f2 := Retry(failNTimes(3), 3)

	if _, err := f1(context.Background()); errors.Is(err, errIntentional) {
		t.Errorf("got %v, want maximum number of retries exceeding", err)
	}
	if _, err := f2(context.Background()); !errors.Is(err, errIntentional) {
		t.Errorf("got %v, want %v", err, errIntentional)
	}
}
