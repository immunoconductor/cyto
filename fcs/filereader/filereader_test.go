package filereader_test

import (
	"testing"

	"github.com/immunoconductor/cyto/fcs/filereader"
)

func TestFCSParser(t *testing.T) {
	f := filereader.NewFCSFileReader("./test-data/test.fcs")
	_, err := f.Read()
	if err != nil {
		t.Errorf("expected Read() to return data")
	}
}
