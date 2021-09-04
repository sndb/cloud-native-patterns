package patterns

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func getCounter() Effector {
	var mu sync.Mutex
	var total int

	return func(ctx context.Context) (string, error) {
		mu.Lock()
		defer mu.Unlock()
		total++
		return fmt.Sprint(total), nil
	}
}

func TestDebounce(t *testing.T) {
	c := Debounce(getCounter(), time.Second)
	var r string
	for i := 0; i < 100; i++ {
		r, _ = c(context.Background())
	}
	if r != "1" {
		t.Errorf("r == %s, want %s", r, "1")
	}

	time.Sleep(time.Second)

	for i := 0; i < 5; i++ {
		r, _ = c(context.Background())
	}
	if r != "2" {
		t.Errorf("r == %s, want %s", r, "2")
	}
}
