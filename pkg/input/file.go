package input

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"prevalentify"
)

const defaultBuffer = 50

// FromFile returns an input generator that generates the input from a file.
func FromFile(file string) (prevalentify.InputGenerator, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening input file: %w", err)
	}

	return func(ctx context.Context) (<-chan prevalentify.ImageSource, <-chan error) {
		inputCh := make(chan prevalentify.ImageSource, defaultBuffer)
		errCh := make(chan error)

		go func() {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				select {
				case inputCh <- prevalentify.ImageSource(scanner.Text()):
				case <-ctx.Done():
					f.Close()
					close(inputCh)
					close(errCh)
					return
				}
			}

			f.Close()
			close(inputCh)
			close(errCh)
		}()

		return inputCh, errCh
	}, nil
}
