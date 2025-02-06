package parser_test

import (
	"testing"

	"github.com/immunoconductor/cyto/fcs/parser"
)

func TestFCSParser(t *testing.T) {
	parser := parser.NewFCSParser("./test-data/test.fcs")
	_, err := parser.Read()
	if err != nil {
		t.Errorf("expected Read() to return data")
	}
}
