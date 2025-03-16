package fcs

import (
	"github.com/immunoconductor/cyto/fcs/constants"
	"github.com/immunoconductor/cyto/fcs/filereader"
	"github.com/immunoconductor/cyto/internal/csv_writer"
)

type FCS struct {
	HEADER FCSHeader
	TEXT   FCSText
	DATA   FCSData
}

type FCSHeader struct {
	Bytes    []byte
	Version  string
	Segments map[constants.SegmentType]FCSSegment
}

type FCSSegment struct {
	Type  constants.SegmentType
	Start int
	End   int
}

type FCSText struct {
	Bytes      []byte
	Keywords   map[string]string
	Parameters []FCSParameter
}

type FCSParameter struct {
	ID int

	// Required fields
	PnB int    // Number of bits reserved for parameter number n.
	PnE string // Amplification type for parameter n.
	PnN string // Short name for parameter n.
	PnR int    // Range for parameter number n.

	// Optional
	PnS string // name for parameter n.
}

type FCSData struct {
	Bytes      []byte
	Mode       string
	DataType   string
	Data       [][]float32
	DataString [][]string // Data is string format
}

func Read(inputFilePath string, transform bool) (*FCS, error) {
	f := filereader.NewFCSFileReader(inputFilePath)
	fcsFileBytes, err := f.Read()
	if err != nil {
		return nil, err
	}

	parser, err := NewFCSParser(fcsFileBytes, transform)
	if err != nil {
		return nil, err
	}

	return parser.Parse()
}

func (f *FCS) ToCSV(outputFilePath string) {
	writer := csv_writer.NewCSVWriter(f.ToTibble(), outputFilePath)
	writer.Write()
}

func (f *FCS) ToShortNameCSV(outputFilePath string) {
	writer := csv_writer.NewCSVWriter(f.ToShortNameTibble(), outputFilePath)
	writer.Write()
}

func (f *FCS) ToTibble() [][]string {
	var names []string
	for _, v := range f.TEXT.Parameters {
		names = append(names, v.PnN)
	}
	return append([][]string{names}, f.DATA.DataString...)
}

func (f *FCS) ToShortNameTibble() [][]string {
	var shortnames []string
	for _, v := range f.TEXT.Parameters {
		shortnames = append(shortnames, v.PnS)
	}
	return append([][]string{shortnames}, f.DATA.DataString...)
}

func (f *FCS) Names() []string {
	var names []string
	for _, v := range f.TEXT.Parameters {
		names = append(names, v.PnN)
	}
	return names
}

func (f *FCS) Keywords() map[string]string {
	return f.TEXT.Keywords
}
