package processor

import (
	"context"
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"prevalentify"
	"testing"
)

func TestNewKMeansHappyPath(t *testing.T) {
	// load images
	images := loadImages(t, []string{
		"../../testdata/processor/img1.png",
		"../../testdata/processor/img2.png",
	})

	kmeanFn := NewKMeans()

	var result []prevalentify.ProcessResult
	for _, image := range images {
		processResult, err := kmeanFn(context.Background(), image)
		if err != nil {
			t.Errorf("unexpected error running kmean on image %s : %s", image.Source, err)
		}

		result = append(result, processResult)
	}

	expectedResults := []prevalentify.ProcessResult{
		{
			ImageSource: "../../testdata/processor/img1.png",
			Colors: []color.Color{
				color.RGBA{156, 193, 188, 255},
				color.RGBA{93, 88, 107, 255},
				color.RGBA{237, 106, 90, 255},
			},
		},
		{
			ImageSource: "../../testdata/processor/img2.png",
			Colors: []color.Color{
				color.RGBA{26, 169, 156, 255},
				color.RGBA{17, 24, 107, 255},
				color.RGBA{173, 221, 218, 255},
			},
		},
	}

	for resultIndex, r := range result {
		require.Equal(t, expectedResults[resultIndex].ImageSource, r.ImageSource)

		// TODO: the output of prominentcolor isn't deterministic, as there can be some variance in the returned colors;
		// we should probably be calculating the euclidian distance between the expected color, and the returned color and
		// assert it's within some predefined threshold.
	}
}

func loadImages(t *testing.T, imageSources []string) []prevalentify.Image {
	var images []prevalentify.Image

	for _, source := range imageSources {
		f, err := os.Open(source)
		if err != nil {
			t.Fatalf("error opening test image %s : %s", source, err)
		}
		img, _, err := image.Decode(f)
		f.Close()

		if err != nil {
			t.Fatalf("error decoding test image %s : %s", source, err)
		}
		images = append(images, prevalentify.Image{
			Image:  &img,
			Source: prevalentify.ImageSource(source),
		})
	}

	return images
}
