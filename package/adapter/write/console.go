package write

import (
	"context"
	"efrainpb/truefit-cashregister/internal/domain"
	"fmt"
	"io"
)

type ConsoleChangeWriter struct {
	results  []domain.ChangeResult
	dest     io.Writer
	currency domain.Currency
}

func NewConsoleChangeWriter(results []domain.ChangeResult, dest io.Writer, currency domain.Currency) *ConsoleChangeWriter {
	return &ConsoleChangeWriter{results: results, dest: dest, currency: currency}
}

func (w *ConsoleChangeWriter) Writer(_ context.Context) error {
	for _, r := range w.results {
		line := w.currency.FormatChange(r.Items)
		if _, err := fmt.Fprintln(w.dest, line); err != nil {
			return fmt.Errorf("failed to write line: %v", err)
		}
	}
	return nil
}
