package processor

import (
	"context"
	"github.com/stretchr/testify/require"
	"image/color"
	"prevalentify"
	"testing"
)

func TestNewSimpleHappyPath(t *testing.T) {
	// load images
	images := loadImages(t, []string{
		"../../testdata/processor/img1.png",
		"../../testdata/processor/img2.png",
	})

	simpleFn := NewSimple(3)

	var result []prevalentify.ProcessResult
	for _, image := range images {
		processResult, err := simpleFn(context.Background(), image)
		if err != nil {
			t.Errorf("unexpected error running kmean on image %s : %s", image.Source, err)
		}

		result = append(result, processResult)
	}

	expectedResults := []prevalentify.ProcessResult{
		{
			ImageSource: "../../testdata/processor/img1.png",
			Colors: []color.Color{
				color.RGBA{155, 193, 188, 255},
				color.RGBA{93, 87, 107, 255},
				color.RGBA{237, 106, 90, 255},
			},
		},
		{
			ImageSource: "../../testdata/processor/img2.png",
			Colors: []color.Color{
				color.RGBA{26, 169, 156, 255},
				color.RGBA{17, 23, 107, 255},
				color.RGBA{173, 221, 218, 255},
			},
		},
	}

	for resultIndex, r := range result {
		require.Equal(t, expectedResults[resultIndex].ImageSource, r.ImageSource)
		require.ElementsMatch(t, expectedResults[resultIndex].Colors, r.Colors)
	}
}
