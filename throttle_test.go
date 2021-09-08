package patterns

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	c, stop := Throttle(Effector(getCounter()), 5, time.Second)
	var count string
	var err error

	for i := 0; i < 5; i++ {
		_, err = c(context.Background())
	}
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}

	_, err = c(context.Background())
	if !strings.Contains(err.Error(), "throttled") {
		t.Errorf("want %v to contain %q", err, "throttled")
	}

	time.Sleep(time.Second*2 + time.Millisecond*50)

	for i := 0; i < 2; i++ {
		count, err = c(context.Background())
	}
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}
	if count != "7" {
		t.Errorf("count == %v, want %v", count, "7")
	}

	_, err = c(context.Background())
	if err == nil || !strings.Contains(err.Error(), "throttled") {
		t.Errorf("want %v to contain %q", err, "throttled")
	}

	stop()

	time.Sleep(time.Second*2 + time.Millisecond*50)

	_, err = c(context.Background())
	if !strings.Contains(err.Error(), "throttled") {
		t.Errorf("want %v to contain %q", err, "throttled")
	}
}
