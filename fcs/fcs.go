package fcs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/immunoconductor/cyto/fcs/constants"
	"github.com/immunoconductor/cyto/fcs/parser"
	"github.com/immunoconductor/cyto/internal/validator"
)

type FCS struct {
	HEADER FCSHeader
	TEXT   []byte
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

	_ = validator.HasRequiredKeywords(t)

	return &FCS{
		HEADER: *h,
		TEXT:   t,
	}, nil
}

func getHeader(byteSlice []byte) (*FCSHeader, error) {
	// fmt.Printf("version identifier: %s, length: %v\n", string(byteSlice[0:6]), len(byteSlice[0:6]))                           // version identifier
	// fmt.Printf("space characters: %s, length: %v\n", string(byteSlice[6:10]), len(byteSlice[6:10]))                           // space characters
	// fmt.Printf("offset to first byte of TEXT segment: %s, length: %v\n", string(byteSlice[10:18]), len(byteSlice[10:18]))     // offset to first byte of TEXT segment
	// fmt.Printf("offset to last byte of TEXT segment: %s, length: %v\n", string(byteSlice[18:26]), len(byteSlice[18:26]))      // offset to last byte of TEXT segment
	// fmt.Printf("offset to first byte of DATA segment: %s, length: %v\n", string(byteSlice[26:34]), len(byteSlice[26:34]))     // offset to first byte of DATA segment
	// fmt.Printf("offset to last byte of DATA segment: %s, length: %v\n", string(byteSlice[34:42]), len(byteSlice[34:42]))      // offset to last byte of DATA segment
	// fmt.Printf("offset to first byte of ANALYSIS segment: %s, length: %v\n", string(byteSlice[42:50]), len(byteSlice[42:50])) // offset to first byte of ANALYSIS segment
	// fmt.Printf("offset to last byte of ANALYSIS segment: %s, length: %v\n", string(byteSlice[50:58]), len(byteSlice[50:58]))  // offset to last byte of ANALYSIS segment

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

	headerBytes := byteSlice[:58] // up-to ANALYSIS segment

	userDefinedSegments := byteSlice[58:*beginningOfTextSegmentInt]
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

func getTextSegment(byteSlice []byte, h *FCSHeader) ([]byte, error) {
	textSegment := h.Segments["TEXT"]

	return byteSlice[textSegment.Start : textSegment.End+1], nil
}

func getOffset(b []byte, start int, end int) (*int, error) {
	ASCIIValue := strings.TrimSpace(string(b[start:end]))
	intValue, err := strconv.Atoi(ASCIIValue)
	if err != nil {
		return nil, err
	}
	return &intValue, nil
}
