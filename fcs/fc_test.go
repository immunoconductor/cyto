package fcs_test

import (
	"fmt"
	"testing"

	"github.com/immunoconductor/cyto/fcs"
)

func TestFCS3_1(t *testing.T) {
	fcs, err := fcs.Read("./filereader/test-data/fcs3.1.fcs", false)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(fcs.HEADER.Version)
	fmt.Println(fcs.HEADER.Segments)
	fmt.Println(fcs.TEXT.Keywords)

	data := fcs.ToTibble()
	fmt.Println(data[0])

	shortNamesData := fcs.ToShortNameTibble()
	fmt.Println(shortNamesData[0])

	fmt.Println(fcs.Names())

}
