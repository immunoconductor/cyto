package fcs

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/immunoconductor/cyto/fcs/constants"
	"github.com/immunoconductor/cyto/fcs/models"
	"github.com/immunoconductor/cyto/fcs/parser"
	"github.com/immunoconductor/cyto/internal/validator"
)

func NewFCS(s string) (*models.FCS, error) {
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

	_ = validator.HasRequiredKeywords(t.Keywords)

	d, err := getDataSegment(t, b[h.Segments["DATA"].Start:h.Segments["DATA"].End+1])
	if err != nil {
		return nil, err
	}

	return &models.FCS{
		HEADER: *h,
		TEXT:   *t,
		DATA:   *d,
	}, nil
}

func getHeader(byteSlice []byte) (*models.FCSHeader, error) {
	// fmt.Printf("version identifier: %s, length: %v\n", string(byteSlice[0:6]), len(byteSlice[0:6]))                           // version identifier
	// fmt.Printf("space characters: %s, length: %v\n", string(byteSlice[6:10]), len(byteSlice[6:10]))                           // space characters
	// fmt.Printf("offset to first byte of TEXT segment: %s, length: %v\n", string(byteSlice[10:18]), len(byteSlice[10:18]))     // offset to first byte of TEXT segment
	// fmt.Printf("offset to last byte of TEXT segment: %s, length: %v\n", string(byteSlice[18:26]), len(byteSlice[18:26]))      // offset to last byte of TEXT segment
	// fmt.Printf("offset to first byte of DATA segment: %s, length: %v\n", string(byteSlice[26:34]), len(byteSlice[26:34]))     // offset to first byte of DATA segment
	// fmt.Printf("offset to last byte of DATA segment: %s, length: %v\n", string(byteSlice[34:42]), len(byteSlice[34:42]))      // offset to last byte of DATA segment
	// fmt.Printf("offset to first byte of ANALYSIS segment: %s, length: %v\n", string(byteSlice[42:50]), len(byteSlice[42:50])) // offset to first byte of ANALYSIS segment
	// fmt.Printf("offset to last byte of ANALYSIS segment: %s, length: %v\n", string(byteSlice[50:58]), len(byteSlice[50:58]))  // offset to last byte of ANALYSIS segment

	// refactor to use constants for offsets
	version := strings.TrimSpace(string(byteSlice[0:6]))

	beginningOfTextSegmentInt, err := getOffset(byteSlice, 10, 18)
	if err != nil {
		return nil, err
	}
	endOfTextSegmentInt, err := getOffset(byteSlice, 18, 26)
	if err != nil {
		return nil, err
	}
	beginningOfDataSegmentInt, err := getOffset(byteSlice, 26, 34)
	if err != nil {
		return nil, err
	}
	endOfDataSegmentInt, err := getOffset(byteSlice, 34, 42)
	if err != nil {
		return nil, err
	}
	beginningOfAnalysisSegmentInt, err := getOffset(byteSlice, 42, 50)
	if err != nil {
		return nil, err
	}
	endOfAnalysisSegmentInt, err := getOffset(byteSlice, 50, 58)
	if err != nil {
		return nil, err
	}

	var segments = map[constants.SegmentType]models.FCSSegment{
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

	headerBytes := byteSlice[:58] // up-to ANALYSIS segment

	userDefinedSegments := byteSlice[58:*beginningOfTextSegmentInt]
	if len(userDefinedSegments) > 0 {
		fmt.Println("user defined segments exist in file")
		fmt.Printf("offset to user defined OTHER segments: %s, length: %v\n", string(userDefinedSegments), len(userDefinedSegments)) // offset to user defined OTHER segments
		headerBytes = append(headerBytes, userDefinedSegments...)                                                                    // including any user defined segments
	}

	return &models.FCSHeader{
		Bytes:    headerBytes,
		Version:  version,
		Segments: segments,
	}, nil
}

func getTextSegment(byteSlice []byte, h *models.FCSHeader) (*models.FCSText, error) {
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
		// fmt.Printf("%s =>  %s\n", recSlice[i], recSlice[i+1])
		Keywords[strings.TrimSpace(textSegmentSlice[i])] = strings.TrimSpace(textSegmentSlice[i+1])
	}

	return &models.FCSText{
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

func getDataSegment(t *models.FCSText, byteSlice []byte) (*models.FCSData, error) {
	data := models.FCSData{
		Bytes:    byteSlice,
		Mode:     strings.TrimSpace(t.Keywords["$MODE"]),
		DataType: strings.TrimSpace(t.Keywords["$DATATYPE"]),
	}

	np, _ := strconv.Atoi(strings.TrimSpace(t.Keywords["$PAR"]))
	ne, _ := strconv.Atoi(strings.TrimSpace(t.Keywords["$TOT"]))
	byteOrder := strings.TrimSpace(t.Keywords["$BYTEORD"])

	// fmt.Println(byteOrder)
	order, err := determineByteOrder(byteOrder)
	if err != nil {
		return nil, err
	}

	fmt.Println(ne, " x ", np)

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

	// TODO: refactor
	for i := range twoDimFloat32Data {
		twoDimFloat32Data[i] = make([]float32, cols)
		for j := range twoDimFloat32Data[i] {
			twoDimFloat32Data[i][j] = float32Data[i*cols+j]
		}
	}

	data.Data = twoDimFloat32Data
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
