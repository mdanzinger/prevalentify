package prevalentify

import (
	"context"
	"sync"
)

// Resolver reads input from the ImageSource, resolves the images and sends them to the returned Image channel.
// Errors are sent to the returned error channel.
type Resolver func(context.Context, <-chan ImageSource) (<-chan Image, <-chan error)

// ResolveFn is a function used to resolve a given ImageSource.
type ResolveFn func(context.Context, ImageSource) (Image, error)

// NewResolverPool returns a pool of resolver workers in the form of a Resolver.
func NewResolverPool(workers int, resolveFn ResolveFn) Resolver {
	return func(ctx context.Context, source <-chan ImageSource) (<-chan Image, <-chan error) {
		// initialize channels
		imageCh := make(chan Image)
		errCh := make(chan error)

		var wg sync.WaitGroup
		wg.Add(workers)

		// spawn workers
		for i := 0; i < workers; i++ {
			go func() {
				resolve(ctx, source, imageCh, errCh, resolveFn)
				wg.Done()
			}()
		}

		// handle shutdown
		go func() {
			wg.Wait()
			close(imageCh)
			close(errCh)
		}()

		return imageCh, errCh
	}
}

func resolve(ctx context.Context, source <-chan ImageSource, out chan Image, errCh chan error, resolveFn ResolveFn) {
	for {
		select {
		case s, ok := <-source:
			if !ok {
				return
			}

			img, err := resolveFn(ctx, s)
			if err != nil {
				errCh <- err
				continue
			}
			out <- img

		case <-ctx.Done():
			return
		}
	}
}
