package prevalentify

import (
	"context"
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	"os"
	"sync"
	"testing"
)

// hard-coded mock result
var prevalentColors = []color.Color{
	color.RGBA{233, 100, 82, 1},
	color.RGBA{133, 50, 12, 1},
	color.RGBA{33, 80, 42, 1},
}

func TestNewProcessorPoolHappyPath(t *testing.T) {
	// generate fake input generator
	inputData := []string{
		"testdata/processor/img1.png",
		"testdata/processor/img2.png",
	}
	resolver := mockResolver(t, inputData...)
	resolveCh, _ := resolver(context.Background(), nil)

	// create resolver
	processor := NewProcessorPool(2, mockProcessFn)
	resultCh, errCh := processor(context.Background(), resolveCh)

	var result []ProcessResult
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for pr := range resultCh {
			result = append(result, pr)
		}
		wg.Done()
	}()

	go func() {
		for err := range errCh {
			t.Fatalf("unexpected error in resolver: %s", err)
		}
		wg.Done()
	}()

	// Wait for channels to close
	wg.Wait()

	// TODO: Revist assertions here, they feel pretty weak
	require.Len(t, result, 2)
	for _, r := range result {
		require.Len(t, r.Colors, 3)
	}
}

func mockResolver(t *testing.T, imageSource ...string) Resolver {
	return func(ctx context.Context, source <-chan ImageSource) (<-chan Image, <-chan error) {
		resolveCh := make(chan Image)
		errCh := make(chan error)

		go func() {
			for _, imgSrc := range imageSource {
				f, err := os.Open(imgSrc)
				if err != nil {
					t.Fatalf("error opening %s : %s", imgSrc, err)
				}

				img, _, err := image.Decode(f)
				if err != nil {
					t.Fatalf("error decoding %s : %s", imgSrc, err)
				}

				resolveCh <- Image{
					Image:  &img,
					Source: ImageSource(imgSrc),
				}

			}
			close(resolveCh)
			close(errCh)
		}()

		return resolveCh, errCh
	}
}

func mockProcessFn(ctx context.Context, image Image) (ProcessResult, error) {
	return ProcessResult{
		ImageSource: image.Source,
		Colors:      prevalentColors,
	}, nil
}
