package prevalentify

import "context"

// OutputWriter reads from the ProcessResult channel and outputs the result.
type OutputWriter func(context.Context, <-chan ProcessResult) <-chan error
