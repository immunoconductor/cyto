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
	"github.com/immunoconductor/cyto/internal/transform/cytof"
	"github.com/immunoconductor/cyto/internal/validator"
)

type FCSParser struct {
	FileBytes []byte
	Transform bool
}

func NewFCSParser(b []byte, transform bool) (*FCSParser, error) {
	return &FCSParser{b, transform}, nil
}

func (p *FCSParser) Parse() (*FCS, error) {
	transform := p.Transform
	fcsFileBytes := p.FileBytes
	h, err := getHeader(p.FileBytes)
	if err != nil {
		return nil, err
	}

	t, err := getTextSegment(fcsFileBytes, h)
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

	err = h.Sanitize(t)
	if err != nil {
		return nil, err
	}

	d, err := getDataSegment(t,
		fcsFileBytes[h.Segments["DATA"].Start:h.Segments["DATA"].End+1],
		transform)
	if err != nil {
		return nil, err
	}

	return &FCS{
		HEADER: *h,
		TEXT:   *t,
		DATA:   *d,
	}, nil
}

func getHeader(byteSlice []byte) (*FCSHeader, error) {
	versionByteOffset := constants.SegmentByteOffsets["Version"]
	version := strings.TrimSpace(string(byteSlice[versionByteOffset[0] : versionByteOffset[1]+1]))

	beginningOfTextSegmentOffset := constants.SegmentByteOffsets["FirstByteTEXTSegment"]
	beginningOfTextSegmentInt, err := getOffsetAndConvertToInt(
		byteSlice, beginningOfTextSegmentOffset[0],
		beginningOfTextSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	endOfTextSegmentOffset := constants.SegmentByteOffsets["LastByteTEXTSegment"]
	endOfTextSegmentInt, err := getOffsetAndConvertToInt(byteSlice, endOfTextSegmentOffset[0], endOfTextSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	beginningOfDataSegmentOffset := constants.SegmentByteOffsets["FirstByteDATASegment"]
	beginningOfDataSegmentInt, err := getOffsetAndConvertToInt(byteSlice, beginningOfDataSegmentOffset[0], beginningOfDataSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	endOfDataSegmentOffset := constants.SegmentByteOffsets["LastByteDATASegment"]
	endOfDataSegmentInt, err := getOffsetAndConvertToInt(byteSlice, endOfDataSegmentOffset[0], endOfDataSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	beginningOfAnalysisSegmentOffset := constants.SegmentByteOffsets["FirstByteANALYSISSegment"]
	beginningOfAnalysisSegmentInt, err := getOffsetAndConvertToInt(byteSlice, beginningOfAnalysisSegmentOffset[0], beginningOfAnalysisSegmentOffset[1]+1)
	if err != nil {
		return nil, err
	}
	endOfAnalysisSegmentOffset := constants.SegmentByteOffsets["LastByteANALYSISSegment"]
	endOfAnalysisSegmentInt, err := getOffsetAndConvertToInt(byteSlice, endOfAnalysisSegmentOffset[0], endOfAnalysisSegmentOffset[1]+1)
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
		fmt.Println()
		fmt.Println("[ User defined segments exist in file ]")
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

func getOffsetAndConvertToInt(b []byte, start int, end int) (*int, error) {
	ASCIIValue := strings.TrimSpace(string(b[start:end]))
	intValue, err := strconv.Atoi(ASCIIValue)
	if err != nil {
		return nil, err
	}
	return &intValue, nil
}

func getDataSegment(t *FCSText, byteSlice []byte, transform bool) (*FCSData, error) {
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

	fmt.Println("[ ", ne, " cells", " x ", np, " observables ]")

	if len(byteSlice) < 4 {
		return nil, fmt.Errorf("not enough bytes for float32: got %d, need 4", len(byteSlice))
	}
	float32Data := make([]float32, np*ne)
	r := bytes.NewReader(byteSlice)
	err = binary.Read(r, order, &float32Data)
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
			dataValue := float32Data[i*cols+j]
			if transform {
				dataValue = float32(cytof.Arcsinh(float64(dataValue), 0, 0.2))
			}
			twoDimFloat32Data[i][j] = dataValue
			twoDimString2Data[i][j] = fmt.Sprintf("%f", dataValue)
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
		return nil, fmt.Errorf("unsupported byte order %s", order)
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
