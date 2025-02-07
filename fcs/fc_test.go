package fcs_test

import (
	"fmt"
	"testing"

	"github.com/immunoconductor/cyto/fcs"
)

func TestFCS3_0(t *testing.T) {
	fcs, err := fcs.NewFCS("./parser/test-data/fcs3.0_2.fcs")
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(fcs.HEADER.Version)
	// fmt.Println(fcs.HEADER.Segments)
	// fmt.Println(fcs.TEXT.Keywords)

}

func TestFCS3_1(t *testing.T) {
	fcs, err := fcs.NewFCS("./parser/test-data/fcs3.1.fcs")
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(fcs.HEADER.Version)
	// fmt.Println(fcs.HEADER.Segments)
	// fmt.Println(fcs.TEXT.Keywords)

}
