package resolver

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"prevalentify"
)

const defaultMaxFileSizeMB = 10

type httpConfig struct {
	client *http.Client

	maxFileSizeMB int64
}

// HTTPOpt is a functional option for the HTTP resolver
type HTTPOpt func(c *httpConfig)

// WithHTTPClient injects a custom http client into the http resolver
func WithHTTPClient(client *http.Client) HTTPOpt {
	return func(c *httpConfig) {
		c.client = client
	}
}

// NewHTTP returns an http resolver
func NewHTTPFn(opts ...HTTPOpt) prevalentify.ResolveFn {
	cfg := &httpConfig{
		client:        http.DefaultClient,
		maxFileSizeMB: defaultMaxFileSizeMB,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return func(ctx context.Context, source prevalentify.ImageSource) (prevalentify.Image, error) {
		// prepare request
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, string(source), nil)
		if err != nil {
			return prevalentify.Image{}, fmt.Errorf("error creating request: %w", err)
		}

		// make request
		resp, err := cfg.client.Do(req)
		if err != nil {
			return prevalentify.Image{}, fmt.Errorf("error performing request: %w", err)
		}
		defer resp.Body.Close()

		// Limit size to avoid exhausting memory
		if resp.ContentLength > cfg.maxFileSizeMB*1024*1024 {
			return prevalentify.Image{}, fmt.Errorf("error resolving image, image exceeds max file limit of %vMB", cfg.maxFileSizeMB)
		}

		// As an extra precaution, since we can't trust the content length, use a limited reader to avoid reading large bodys.
		img, _, err := image.Decode(io.LimitReader(resp.Body, cfg.maxFileSizeMB*1024*1024))
		if err != nil {
			return prevalentify.Image{}, fmt.Errorf("error decoding image: %w", err)
		}

		return prevalentify.Image{
			Image:  &img,
			Source: prevalentify.ImageSource(resp.Request.URL.String()),
		}, nil
	}
}
