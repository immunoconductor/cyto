package fcs

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/immunoconductor/cyto/fcs/constants"
	"github.com/immunoconductor/cyto/fcs/parser"
	"github.com/immunoconductor/cyto/internal/csv_writer"
	"github.com/immunoconductor/cyto/internal/validator"
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
	DataString [][]string // Data is dtring format
}

func NewFCS(s string) (*FCS, error) {
	parser := parser.NewFCSParser(s)
	b, err := parser.Read()
	if err != nil {
		return nil, err
	}

	h, err := getHeader(b)
	if err != nil {
		return nil, err
	}

	t, err := getTextSegment(b, h)
	if err != nil {
		return nil, err
	}

	valid := validator.HasRequiredKeywords(t.Keywords)
	if !valid {
		return nil, errors.New("missing required keywords")
	}

	parameters, err := getParameterMetadata(t.Keywords)
	if err != nil {
		return nil, err
	}
	t.Parameters = parameters

	d, err := getDataSegment(t, b[h.Segments["DATA"].Start:h.Segments["DATA"].End+1])
	if err != nil {
		return nil, err
	}

	return &FCS{
		HEADER: *h,
		TEXT:   *t,
		DATA:   *d,
	}, nil
}

func (f *FCS) ToCSV(path string) {
	writer := csv_writer.NewCSVWriter(f.ToTibble(), path)
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

func getHeader(byteSlice []byte) (*FCSHeader, error) {
	versionByteOffset := constants.SegmentByteOffsets["Version"]
	version := strings.TrimSpace(string(byteSlice[versionByteOffset[0] : versionByteOffset[1]+1]))

	beginningOfTextSegmentOffset := constants.SegmentByteOffsets["FirstByteTEXTSegment"]
	beginningOfTextSegmentInt, err := getOffset(
		byteSlice, beginningOfTextSegmentOffset[0],
		beginningOfTextSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	endOfTextSegmentOffset := constants.SegmentByteOffsets["LastByteTEXTSegment"]
	endOfTextSegmentInt, err := getOffset(byteSlice, endOfTextSegmentOffset[0], endOfTextSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	beginningOfDataSegmentOffset := constants.SegmentByteOffsets["FirstByteDATASegment"]
	beginningOfDataSegmentInt, err := getOffset(byteSlice, beginningOfDataSegmentOffset[0], beginningOfDataSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	endOfDataSegmentOffset := constants.SegmentByteOffsets["LastByteDATASegment"]
	endOfDataSegmentInt, err := getOffset(byteSlice, endOfDataSegmentOffset[0], endOfDataSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	beginningOfAnalysisSegmentOffset := constants.SegmentByteOffsets["FirstByteANALYSISSegment"]
	beginningOfAnalysisSegmentInt, err := getOffset(byteSlice, beginningOfAnalysisSegmentOffset[0], beginningOfAnalysisSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	endOfAnalysisSegmentOffset := constants.SegmentByteOffsets["LastByteANALYSISSegment"]
	endOfAnalysisSegmentInt, err := getOffset(byteSlice, endOfAnalysisSegmentOffset[0], endOfAnalysisSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}

	var segments = map[constants.SegmentType]FCSSegment{
		constants.TEXT: {
			Type:  constants.TEXT,
			Start: *beginningOfTextSegmentInt,
			End:   *endOfTextSegmentInt,
		},
		constants.DATA: {
			Type:  constants.DATA,
			Start: *beginningOfDataSegmentInt,
			End:   *endOfDataSegmentInt,
		},
		constants.ANALYSIS: {
			Type:  constants.ANALYSIS,
			Start: *beginningOfAnalysisSegmentInt,
			End:   *endOfAnalysisSegmentInt,
		},
	}

	headerBytes := byteSlice[:endOfAnalysisSegmentOffset[1]+1] // up-to ANALYSIS segment

	userDefinedSegments := byteSlice[endOfAnalysisSegmentOffset[1]+1 : *beginningOfTextSegmentInt]
	if len(userDefinedSegments) > 0 {
		fmt.Println("user defined segments exist in file")
		fmt.Printf("offset to user defined OTHER segments: %s, length: %v\n", string(userDefinedSegments), len(userDefinedSegments)) // offset to user defined OTHER segments
		headerBytes = append(headerBytes, userDefinedSegments...)                                                                    // including any user defined segments
	}

	return &FCSHeader{
		Bytes:    headerBytes,
		Version:  version,
		Segments: segments,
	}, nil
}

func getTextSegment(byteSlice []byte, h *FCSHeader) (*FCSText, error) {
	textSegment := h.Segments["TEXT"]

	textSegmentBytes := byteSlice[textSegment.Start : textSegment.End+1]
	textSegmentString := string(textSegmentBytes)

	delimeter := textSegmentString[:1]
	r := csv.NewReader(strings.NewReader(string(textSegmentBytes)))
	r.Comma = []rune(delimeter)[0]

	segment, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	textSegmentSlice := segment[0][1 : len(segment[0])-1]
	var Keywords = make(map[string]string)
	for i := 0; i < len(textSegmentSlice); i = i + 2 {
		Keywords[strings.TrimSpace(textSegmentSlice[i])] = strings.TrimSpace(textSegmentSlice[i+1])
	}

	return &FCSText{
		Bytes:    textSegmentBytes,
		Keywords: Keywords,
	}, nil
}

func getOffset(b []byte, start int, end int) (*int, error) {
	ASCIIValue := strings.TrimSpace(string(b[start:end]))
	intValue, err := strconv.Atoi(ASCIIValue)
	if err != nil {
		return nil, err
	}
	return &intValue, nil
}

func getDataSegment(t *FCSText, byteSlice []byte) (*FCSData, error) {
	data := FCSData{
		Bytes:    byteSlice,
		Mode:     strings.TrimSpace(t.Keywords["$MODE"]),
		DataType: strings.TrimSpace(t.Keywords["$DATATYPE"]),
	}

	np, _ := strconv.Atoi(strings.TrimSpace(t.Keywords["$PAR"]))
	ne, _ := strconv.Atoi(strings.TrimSpace(t.Keywords["$TOT"]))
	byteOrder := strings.TrimSpace(t.Keywords["$BYTEORD"])

	order, err := determineByteOrder(byteOrder)
	if err != nil {
		return nil, err
	}

	fmt.Println(ne, " cells", " x ", np, " observations")

	float32Data := make([]float32, np*ne)
	r := bytes.NewReader(byteSlice)
	err = binary.Read(r, order, &float32Data) // determine endian
	if err != nil {
		log.Fatal("binary.Read failed ", err)
		return nil, err
	}

	rows := ne
	cols := np
	twoDimFloat32Data := make([][]float32, rows)
	twoDimString2Data := make([][]string, rows)

	// TODO: refactor
	for i := range twoDimFloat32Data {
		twoDimFloat32Data[i] = make([]float32, cols)
		twoDimString2Data[i] = make([]string, cols)

		for j := range twoDimFloat32Data[i] {
			twoDimFloat32Data[i][j] = float32Data[i*cols+j]
			twoDimString2Data[i][j] = fmt.Sprintf("%f", float32Data[i*cols+j])
		}
	}
	data.Data = twoDimFloat32Data
	data.DataString = twoDimString2Data
	return &data, nil
}

func determineByteOrder(order string) (binary.ByteOrder, error) {
	switch order {
	case "1,2,3,4":
		return binary.LittleEndian, nil
	case "4,3,2,1":
		return binary.BigEndian, nil
	default:
		return nil, fmt.Errorf("unknown byte order %s", order)
	}
}

func getParameterMetadata(keywords map[string]string) ([]FCSParameter, error) {
	var parameters []FCSParameter

	np, err := strconv.Atoi(strings.TrimSpace(keywords["$PAR"]))
	if err != nil {
		return nil, fmt.Errorf("could not convert %s to int", keywords["$PAR"])
	}

	for i := 1; i <= np; i++ {
		pnNValue, pnNExists := keywords[fmt.Sprintf("$P%dN", i)]
		pnBValue, pnBExists := keywords[fmt.Sprintf("$P%dB", i)]
		pnE, pnEExists := keywords[fmt.Sprintf("$P%dE", i)]
		pnRValue, pnRExists := keywords[fmt.Sprintf("$P%dR", i)]

		if !pnNExists || !pnBExists || !pnEExists || !pnRExists {
			return nil, fmt.Errorf("missing required parameter keyword: %s", fmt.Sprintf("$P%dN", i))
		}

		pnB, err := strconv.Atoi(pnBValue)
		if err != nil {
			return nil, err
		}
		pnR, err := strconv.Atoi(pnRValue)
		if err != nil {
			return nil, err
		}

		parameter := FCSParameter{
			ID:  i,
			PnN: pnNValue,
			PnS: keywords[fmt.Sprintf("$P%dS", i)],
			PnB: pnB,
			PnE: pnE,
			PnR: pnR,
		}
		parameters = append(parameters, parameter)

	}

	return parameters, nil
}
