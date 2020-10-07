package output

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"prevalentify"
	"sync"
)

// NewCSVFileWriter returns a new CSV prevalentify.OutputWriter writing to the passes in file
func NewCSVFileWriter(file string) (prevalentify.OutputWriter, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, fmt.Errorf("error creating file for CSV OutputWriter: %w", err)
	}

	return NewCSV(f), nil
}

// NewCSV returns a CSV prevalentify.OutputWriter, writing to the injected io.Writer. It writes in the form of:
// imageSource, color1, color2, colorN...
func NewCSV(w io.Writer) prevalentify.OutputWriter {
	return func(ctx context.Context, results <-chan prevalentify.ProcessResult) <-chan error {
		errCh := make(chan error)
		csvWriter := csv.NewWriter(w)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			csvWrite(ctx, csvWriter, results, errCh)
			wg.Done()
		}()

		go func() {
			wg.Wait()
			close(errCh)

			// it's possible we are writing to a closable writer, like a file, we need to close the writer if it's closable.
			if closer, ok := w.(io.WriteCloser); ok {
				closer.Close()
			}
		}()

		return errCh

	}

}

func csvWrite(ctx context.Context, writer *csv.Writer, resultCh <-chan prevalentify.ProcessResult, errCh chan<- error) {
	for {
		select {
		case result, ok := <-resultCh:
			if !ok {
				return
			}
			if err := writer.Write(toRecord(result)); err != nil {
				errCh <- fmt.Errorf("error writing to csv file: %w", err)
			}

		case <-ctx.Done():
			writer.Flush()
			return
		}

		writer.Flush()
	}
}

func toRecord(result prevalentify.ProcessResult) []string {
	var record []string
	record = append(record, string(result.ImageSource))

	for _, c := range result.Colors {
		r, g, b, _ := c.RGBA()
		record = append(record, fmt.Sprintf("#%02x%02x%02x", uint8(r), uint8(g), uint8(b)))
	}

	return record
}
