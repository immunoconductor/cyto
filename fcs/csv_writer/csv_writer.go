package csv_writer

import "github.com/immunoconductor/cyto/internal/writer"

type CSVWriter struct {
	OutputPath string
}

func NewCSVWriter(path string) writer.Writer {
	return &CSVWriter{OutputPath: path}
}

func (w *CSVWriter) Write() error {
	return nil
}
