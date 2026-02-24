package write

import (
	"context"
	"efrainpb/truefit-cashregister/internal/domain"
	"fmt"
	"io"
	"os"
)

type FileChangeWriter struct {
	path     string
	file     *os.File
	results  []domain.ChangeResult
	dest     io.Writer
	currency domain.Currency
}

func NewFileChangeWriter(path string, results []domain.ChangeResult, dest io.Writer, currency domain.Currency) *FileChangeWriter {
	return &FileChangeWriter{path: path, results: results, dest: dest, currency: currency}
}

func (w *FileChangeWriter) Writer(_ context.Context) error {
	if w.file == nil {
		f, err := os.Create(w.path)
		if err != nil {
			return fmt.Errorf("creating output file %q: %w", w.path, err)
		}
		w.file = f
	}

	for _, result := range w.results {
		line := w.currency.FormatChange(result.Items)
		if _, err := fmt.Fprintln(w.file, line); err != nil {
			return fmt.Errorf("writing result to file: %w", err)
		}
	}
	return nil
}

func (w *FileChangeWriter) Close() error {
	if w.file == nil {
		return nil
	}
	if err := w.file.Close(); err != nil {
		return fmt.Errorf("closing output file %q: %w", w.path, err)
	}
	return nil
}
