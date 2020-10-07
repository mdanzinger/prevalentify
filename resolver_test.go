package prevalentify

import (
	"context"
	"github.com/stretchr/testify/require"
	"image"
	"os"
	"sync"
	"testing"

	_ "image/jpeg"
	_ "image/png"
)

func TestNewResolverPoolHappyPath(t *testing.T) {
	// generate fake input generator
	inputData := []string{
		"testdata/resolver/img1.jpg",
		"testdata/resolver/img2.png",
		"testdata/resolver/img3.png",
	}
	inputGen := mockInputter(inputData...)
	inputCh, _ := inputGen(context.Background())

	// create resolver
	resolver := NewResolverPool(2, mockResolveFn)
	imageCh, errCh := resolver(context.Background(), inputCh)

	var result []string
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for img := range imageCh {
			result = append(result, string(img.Source))
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

	require.ElementsMatch(t, result, inputData, "expected resolver pool to resolve %v images, but got :\n %v", inputData, result)
}

func mockInputter(data ...string) InputGenerator {
	return func(ctx context.Context) (<-chan ImageSource, <-chan error) {
		sourceCh := make(chan ImageSource)
		errCh := make(chan error)

		go func() {
			for _, inputData := range data {
				sourceCh <- ImageSource(inputData)
			}
			close(sourceCh)
			close(errCh)
		}()

		return sourceCh, errCh
	}
}

func mockResolveFn(ctx context.Context, source ImageSource) (Image, error) {
	f, err := os.Open(string(source))
	if err != nil {
		return Image{}, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return Image{}, err
	}

	return Image{
		Source: ImageSource(f.Name()),
		Image:  &img,
	}, nil
}
