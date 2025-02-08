package constants

// TEXT segment required keywords
type TextSegmentKeyword string

const (
	BEGINANALYSIS TextSegmentKeyword = "$BEGINANALYSIS"
	QUANTITY      TextSegmentKeyword = "$BEGINDATA"
	BEGINSTEXT    TextSegmentKeyword = "$BEGINSTEXT"
	BYTEORD       TextSegmentKeyword = "$BYTEORD"
	DATATYPE      TextSegmentKeyword = "$DATATYPE"
	ENDANALYSIS   TextSegmentKeyword = "$ENDANALYSIS"
	ENDDATA       TextSegmentKeyword = "$ENDDATA"
	ENDSTEXT      TextSegmentKeyword = "$ENDSTEXT"
	MODE          TextSegmentKeyword = "$MODE"
	NEXTDATA      TextSegmentKeyword = "$NEXTDATA"
	PAR           TextSegmentKeyword = "$PAR"
	TOT           TextSegmentKeyword = "$TOT"
	PnB           TextSegmentKeyword = "$PnB"
	PnE           TextSegmentKeyword = "$PnE"
	PnN           TextSegmentKeyword = "$PnN"
	PnR           TextSegmentKeyword = "$PnR"
)

var SegmentByteOffsets = map[string][]int{
	"Version":                  {0, 5},
	"SpaceCharacters":          {6, 9},
	"FirstByteTEXTSegment":     {10, 17},
	"LastByteTEXTSegment":      {18, 25},
	"FirstByteDATASegment":     {26, 33},
	"LastByteDATASegment":      {34, 41},
	"FirstByteANALYSISSegment": {42, 49},
	"LastByteANALYSISSegment":  {50, 57},
}

type SegmentType string

const (
	TEXT     SegmentType = "TEXT"
	DATA     SegmentType = "DATA"
	ANALYSIS SegmentType = "ANALYSIS"
	OTHER    SegmentType = "OTHER"
)

var TextSegmentRequiredParameterKeywords = []string{
	"$P%dB",
	"$P%dE",
	"$P%dN",
	"$P%dR",
}

var TextSegmentRequiredKeywords = []TextSegmentKeyword{
	"$BEGINANALYSIS",
	"$BEGINDATA",
	"$BEGINSTEXT",
	"$BYTEORD",
	"$DATATYPE",
	"$ENDANALYSIS",
	"$ENDDATA",
	"$ENDSTEXT",
	"$MODE",
	"$NEXTDATA",
	"$PAR",
	"$TOT",
}
