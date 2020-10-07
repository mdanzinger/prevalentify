package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"prevalentify"
	"prevalentify/pkg/input"
	"prevalentify/pkg/output"
	"prevalentify/pkg/processor"
	"prevalentify/pkg/resolver"
)

const (
	defaultWorkerCount = 30
	defaultInputFile   = "input.txt"
	defaultOutputfile  = "output.csv"
)

// TODO: considering defaulting output to stdout, generating files automatically generally isn't good UX
var (
	// TODO: consider splitting out the resolver count and processor count into separate flags
	workers    = flag.Int("workers", defaultWorkerCount, "number of concurrent image resolvers / processor workers")
	inputSrc   = flag.String("input", defaultInput(), "input file to process prevalent colors from")
	outputDest = flag.String("output", defaultOutput(), "prevalentify output destination")
)

func main() {
	flag.Parse()

	logger := logrus.New()

	// setup steps for the pipeline
	input, err := input.FromFile(*inputSrc)
	if err != nil {
		logger.Fatal(err)
	}

	resolveFn := resolver.NewHTTPFn()
	resolverPool := prevalentify.NewResolverPool(*workers, resolveFn)

	processFn := processor.NewSimple(3)
	processorPool := prevalentify.NewProcessorPool(*workers, processFn)

	writer, err := output.NewCSVFileWriter(*outputDest)
	if err != nil {
		logger.Fatal(err)
	}

	// build pipeline
	pipeline := prevalentify.NewPipeline(input, resolverPool, processorPool, writer, logger)

	ctx := context.Background()

	log.Printf("starting to process %s", *inputSrc)
	pipeline.Start(ctx)
	log.Printf("processing complete, see output results in %s", *outputDest)
}

func defaultOutput() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting working directory: %w", err)
	}

	return fmt.Sprintf("%s/%s", dir, defaultOutputfile)
}

func defaultInput() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting working directory: %w", err)
	}

	return fmt.Sprintf("%s/%s", dir, defaultInputFile)
}
