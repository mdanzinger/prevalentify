package prevalentify

import "context"

// InputGenerator generates input for the prevalentify pipeline.
type InputGenerator func(ctx context.Context) (<-chan ImageSource, <-chan error)
