package patterns

import "context"

func Timeout(f SlowFunction) WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		rc := make(chan string)
		errc := make(chan error)

		go func() {
			r, err := f(arg)
			rc <- r
			errc <- err
		}()

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case r := <-rc:
			return r, <-errc
		}
	}
}
