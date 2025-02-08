package models

import "github.com/immunoconductor/cyto/fcs/constants"

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
	Bytes    []byte
	Keywords map[string]string
}

type FCSData struct {
	Bytes    []byte
	Mode     string
	DataType string
	Data     [][]float32
}
