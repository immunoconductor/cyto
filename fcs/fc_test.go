package fcs_test

import (
	"fmt"
	"testing"

	"github.com/immunoconductor/cyto/fcs"
)

func TestFCS(t *testing.T) {
	fcs, err := fcs.NewFCS("./parser/test-data/fcs3.0.fcs")
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(fcs.HEADER.Version)
	fmt.Println(fcs.HEADER.Segments)
	fmt.Println(string(fcs.HEADER.Bytes))

	fmt.Println(string(fcs.TEXT))

}
