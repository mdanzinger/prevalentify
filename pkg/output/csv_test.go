package output

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"

	//"github.com/stretchr/testify/require"
	"image/color"
	"prevalentify"
	"sync"
	"testing"
)

func TestNewCSVHappyPath(t *testing.T) {
	// create result Ch
	resultCh := mockResultCh([]prevalentify.ProcessResult{
		{
			ImageSource: "http://imgur.com/test_image.png",
			Colors: []color.Color{
				color.RGBA{123, 231, 152, 1},
				color.RGBA{92, 21, 250, 1},
				color.RGBA{50, 50, 33, 1},
			},
		},
		{
			ImageSource: "http://imgur.com/test_image_2.png",
			Colors: []color.Color{
				color.RGBA{255, 255, 255, 1},
				color.RGBA{100, 100, 100, 1},
				color.RGBA{5, 50, 33, 1},
			},
		},
	})

	buffer := &bytes.Buffer{}
	writer := NewCSV(buffer)

	var wg sync.WaitGroup
	wg.Add(1)

	// start writer
	errCh := writer(context.Background(), resultCh)

	// drain error channel, ensure there are no unexpected errors
	go func() {
		for err := range errCh {
			t.Fatalf("unexpected error in csv writer: %s", err)
		}
		wg.Done()
	}()

	// wait for error channel to close / signal writing is complete.
	wg.Wait()

	expectedResult := `http://imgur.com/test_image.png,#7be798,#5c15fa,#323221
http://imgur.com/test_image_2.png,#ffffff,#646464,#053221
`

	require.Equal(t, expectedResult, buffer.String())
}

func mockResultCh(results []prevalentify.ProcessResult) chan prevalentify.ProcessResult {
	resultCh := make(chan prevalentify.ProcessResult)

	go func() {
		for _, result := range results {
			resultCh <- result
		}
		close(resultCh)
	}()

	return resultCh
}
