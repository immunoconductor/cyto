package csv_writer

import (
	"encoding/csv"
	"os"

	"github.com/immunoconductor/cyto/internal/writer"
)

type CSVWriter struct {
	Data       [][]string
	OutputPath string
}

func NewCSVWriter(data [][]string, path string) writer.Writer {
	return &CSVWriter{data, path}
}

func (w *CSVWriter) Write() error {
	// Create a file
	file, err := os.Create(w.OutputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Initialize csv writer
	writer := csv.NewWriter(file)

	defer writer.Flush()

	// Write all rows at once
	err = writer.WriteAll(w.Data)
	if err != nil {
		return err
	}
	return nil
}
