package processor

import (
	"context"
	extractor "github.com/marekm4/color-extractor"
	"image/color"
	"prevalentify"
)

var defaultSimpleConfig = extractor.Config{
	DownSizeTo:  224,
	SmallBucket: 0.001,
}

// NewSimple returns a simple implementation of an image processor. Under the hood this is leveraging the
// github.com/marekm4/color-extractor package to extract prevalent colors.
func NewSimple(maxColors int) prevalentify.ProcessFn {
	return func(ctx context.Context, image prevalentify.Image) (prevalentify.ProcessResult, error) {
		colors := extractor.ExtractColorsWithConfig(*image.Image, defaultSimpleConfig)
		return prevalentify.ProcessResult{
			ImageSource: image.Source,
			Colors:      limit(colors, maxColors),
		}, nil
	}
}

func limit(colors []color.Color, max int) []color.Color {
	if len(colors) <= max {
		return colors
	}
	return colors[:max]
}
