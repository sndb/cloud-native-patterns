package patterns

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func getCounter() Circuit {
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

func TestDebounceDataRace(t *testing.T) {
	ctx := context.Background()
	circuit := failAfter(1)
	debounce := Debounce(circuit, time.Second)
	var wg sync.WaitGroup

	for count := 1; count <= 10; count++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			_, err := debounce(ctx)
			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	time.Sleep(time.Second * 2)

	for count := 1; count <= 10; count++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			_, err := debounce(ctx)
			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	wg.Wait()
}
