package patterns

import (
	"context"
	"sync"
	"time"
)

type Future interface {
	Result() (string, error)
}

type InnerFuture struct {
	once sync.Once
	done chan struct{}

	r    string
	err  error
	rc   <-chan string
	errc <-chan error
}

func (f *InnerFuture) Result() (string, error) {
	f.once.Do(func() {
		f.r, f.err = <-f.rc, <-f.errc
		close(f.done)
	})

	<-f.done
	return f.r, f.err
}

func Promiser(ctx context.Context) Future {
	rc := make(chan string)
	errc := make(chan error)

	go func() {
		select {
		case <-time.After(time.Second * 2):
			rc <- "I slept for 2 seconds"
			errc <- nil
		case <-ctx.Done():
			rc <- ""
			errc <- ctx.Err()
		}
	}()

	return &InnerFuture{rc: rc, errc: errc}
}
