package parser

import (
	"os"

	"github.com/immunoconductor/cyto/internal/reader"
)

type FCSParser struct {
	FilePath string
	Contents []byte
}

func NewFCSParser(path string) reader.Reader {
	return &FCSParser{FilePath: path}
}

func (p *FCSParser) Read() ([]byte, error) {
	byteSlice, err := os.ReadFile(p.FilePath)
	if err != nil {
		return nil, err
	}

	return byteSlice, nil

}
