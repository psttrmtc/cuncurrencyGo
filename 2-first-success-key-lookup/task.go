package main

import (
	"context"
)

type Getter interface {
	Get(ctx context.Context, address, key string) (string, error)
}

// Call `Getter.Get()` for each address in parallel.
// Returns the first successful response.
// If all requests fail, returns an error.
func Get(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {
	if len(addresses) <= 0 {
		return "", nil
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	resCh := make(chan string)
	errorsCh := make(chan error, len(addresses))
	for _, address := range addresses {
		go func() {
			if val, err := getter.Get(ctx, address, key); err != nil {
				errorsCh <- err
			} else {
				resCh <- val
			}
		}()
	}
	var errCount int

	for {
		select {
		case err := <-errorsCh:
			errCount++
			if errCount == len(addresses) {
				return "", err
			}
		case val := <-resCh:
			return val, nil
		case <-ctx.Done():
			return "", context.Canceled
		}
	}
}
