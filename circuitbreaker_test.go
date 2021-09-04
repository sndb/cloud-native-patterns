package patterns

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func getResettingCounter() Circuit {
	var timer *time.Timer
	var total, i int
	var mu sync.Mutex

	return func(ctx context.Context) (string, error) {
		mu.Lock()
		defer mu.Unlock()

		total++
		if i >= 3 {
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(time.Second*1, func() {
				mu.Lock()
				defer mu.Unlock()
				i = 0
			})
			return "", errors.New("please wait")
		}
		i++
		return fmt.Sprint(total), nil
	}
}

func TestCircuitBreaker(t *testing.T) {
	e := Breaker(getResettingCounter(), 3)
	for i := 0; i < 100; i++ {
		e(context.Background())
	}
	time.Sleep(time.Second * 3)

	total, _ := e(context.Background())
	if total != "7" {
		t.Errorf("total == %s, want %s", total, "7")
	}
	total, _ = e(context.Background())
	if total != "8" {
		t.Errorf("total == %s, want %s", total, "8")
	}
}

func failAfter(threshold int) Circuit {
	count := 0

	return func(ctx context.Context) (string, error) {
		count++
		if count > threshold {
			return "", errors.New("intentional fail!")
		}
		return "success", nil
	}
}

func waitAndContinue() Circuit {
	return func(ctx context.Context) (string, error) {
		time.Sleep(time.Second)
		if rand.Int()%2 == 0 {
			return "success", nil
		}
		return "", fmt.Errorf("forced failure")
	}
}

func TestCircuitBreakerWithFailAfter(t *testing.T) {
	circuit := failAfter(5)
	breaker := Breaker(circuit, 1)

	circuitOpen := false
	doesCircuitOpen := false
	doesCircuitReclose := false

	count := 0
	for range time.NewTicker(time.Second).C {
		_, err := breaker(context.Background())

		if err != nil {
			if err.Error() == "open circuit" {
				if !circuitOpen {
					circuitOpen = true
					doesCircuitOpen = true

					t.Log("circuit has opened")
				}
			} else {
				if circuitOpen {
					circuitOpen = false
					doesCircuitReclose = true

					t.Log("circuit has automatically closed")
				}
			}
		} else {
			t.Log("circuit closed and operational")
		}

		count++
		if count >= 10 {
			break
		}
	}

	if !doesCircuitOpen {
		t.Error("circuit didn't appear to open")
	}
	if !doesCircuitReclose {
		t.Error("circuit didn't appear to close after time")
	}
}

func TestCircuitBreakerDataRace(t *testing.T) {
	ctx := context.Background()

	circuit := waitAndContinue()
	breaker := Breaker(circuit, 1)

	var wg sync.WaitGroup
	for count := 1; count <= 20; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := breaker(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}
	wg.Wait()
}
