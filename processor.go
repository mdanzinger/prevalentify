package prevalentify

import (
	"context"
	"image/color"
	"sync"
)

// ProcessResult contains the result of processing an image. It contains a reference to the ImageSource, as well as
// the top N most prevalent colors.
type ProcessResult struct {
	ImageSource ImageSource
	Colors      []color.Color
}

// Processor reads off of the resolver image channel, processes the image, and sends it to the returned Result channel.
// Errors are sent to the returned error channel.
type Processor func(context.Context, <-chan Image) (<-chan ProcessResult, <-chan error)

// ProcessFn defines the function for actually performing the processing.
type ProcessFn func(context.Context, Image) (ProcessResult, error)

// NewProcessorPool returns a pool of processing workers in the form of a Processor.
func NewProcessorPool(workers int, processFn ProcessFn) Processor {
	return func(ctx context.Context, image <-chan Image) (<-chan ProcessResult, <-chan error) {
		// initialize channels
		resultCh := make(chan ProcessResult)
		errCh := make(chan error)

		var wg sync.WaitGroup
		wg.Add(workers)

		// spawn workers
		for i := 0; i < workers; i++ {
			go func() {
				process(ctx, image, resultCh, errCh, processFn)
				wg.Done()
			}()
		}

		// handle shutdown
		go func() {
			wg.Wait()
			close(resultCh)
			close(errCh)
		}()

		return resultCh, errCh
	}
}

func process(ctx context.Context, image <-chan Image, out chan ProcessResult, errCh chan error, processFn ProcessFn) {
	for {
		select {
		case img, ok := <-image:
			if !ok {
				return
			}

			result, err := processFn(ctx, img)
			if err != nil {
				errCh <- err
				continue
			}
			out <- result

		case <-ctx.Done():
			return
		}
	}
}
