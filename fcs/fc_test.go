package fcs_test

import (
	"fmt"
	"testing"

	"github.com/immunoconductor/cyto/fcs"
)

func TestFCS(t *testing.T) {
	fcs, err := fcs.NewFCS("./parser/test-data/test.fcs")
	if err != nil {
		t.Errorf("expected new FCS object to be created")
	}

	fmt.Println(fcs.HEADER.Version)
	fmt.Println(fcs.HEADER.Segments)
	fmt.Println(string(fcs.HEADER.Bytes))
}
