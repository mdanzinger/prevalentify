package prevalentify

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

// Pipeline handles the orchestration / lifecycle of the pipeline.
type Pipeline struct {
	input     InputGenerator
	resolver  Resolver
	processor Processor
	writer    OutputWriter

	l *logrus.Logger
}

// NewPipeline returns an instantiated pipeline
func NewPipeline(input InputGenerator, resolver Resolver, processor Processor, writer OutputWriter, l *logrus.Logger) *Pipeline {
	p := &Pipeline{
		input:     input,
		resolver:  resolver,
		processor: processor,
		writer:    writer,
		l:         l,
	}

	return p
}

// Start starts running the pipeline and is a blocking call.
func (p *Pipeline) Start(ctx context.Context) {
	// begin generating input.
	sourceCh, inputErrCh := p.input(ctx)

	// resolve images from the ImageSource
	imageCh, resolveErrCh := p.resolver(ctx, sourceCh)

	// process resolved images
	resultCh, processErrCh := p.processor(ctx, imageCh)

	// write results to the writer
	writerErrCh := p.writer(ctx, resultCh)

	p.handleErrors(ctx, inputErrCh, resolveErrCh, processErrCh, writerErrCh)
}

// handleErrors merges the error channels and logs any errors
func (p *Pipeline) handleErrors(ctx context.Context, errChs ...<-chan error) {
	var wg sync.WaitGroup
	wg.Add(len(errChs))

	log := func(errCh <-chan error) {
		for err := range errCh {
			p.l.Error(err)
		}
		wg.Done()
	}

	for _, errCh := range errChs {
		go log(errCh)
	}

	wg.Wait()
}
