package processor

import (
	"context"
	"github.com/EdlinOrg/prominentcolor"
	"image/color"
	"prevalentify"
)

const (
	defaultK = 3

	// TODO: Make the arguments configurable
	defaultArgs = prominentcolor.ArgumentNoCropping
)

type kmeanCfg struct {
	k     int
	masks []prominentcolor.ColorBackgroundMask
}

// KMeanOpt represents a optional param for the KMean processFn
type KMeanOpt func(c *kmeanCfg)

// TopNColors configures the KMeans processor to return the N most prevalent colors.
func TopNColors(n int) KMeanOpt {
	return func(c *kmeanCfg) {
		c.k = n
	}
}

// WithMasks configures the processor to mask out colors from the image.
func WithMasks(masks []prominentcolor.ColorBackgroundMask) KMeanOpt {
	return func(c *kmeanCfg) {
		c.masks = masks
	}
}

// NewKMeans returns an image processFn that utilizes the github.com/EdlinOrg/prominentcolor library's k-means implementation
// to extract the most dominant colors in a given image.
func NewKMeans(opts ...KMeanOpt) prevalentify.ProcessFn {
	cfg := &kmeanCfg{
		k: defaultK,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return func(_ context.Context, image prevalentify.Image) (prevalentify.ProcessResult, error) {
		colors, err := prominentcolor.KmeansWithAll(3, *image.Image, defaultArgs, 224, cfg.masks)
		if err != nil {
			return prevalentify.ProcessResult{}, err
		}

		return prevalentify.ProcessResult{
			ImageSource: image.Source,
			Colors:      toColors(colors),
		}, nil
	}
}

func toColors(colors []prominentcolor.ColorItem) []color.Color {
	var result []color.Color
	for _, c := range colors {
		result = append(result, color.RGBA{R: uint8(c.Color.R), G: uint8(c.Color.G), B: uint8(c.Color.B), A: 1})
	}
	return result
}
