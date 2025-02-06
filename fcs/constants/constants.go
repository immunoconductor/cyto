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
	PnB           TextSegmentKeyword = "$PnB"
	PnE           TextSegmentKeyword = "$PnE"
	PnN           TextSegmentKeyword = "$PnN"
	PnR           TextSegmentKeyword = "$PnR"
	TOT           TextSegmentKeyword = "$TOT"
)

type SegmentType string

const (
	TEXT     SegmentType = "TEXT"
	DATA     SegmentType = "DATA"
	ANALYSIS SegmentType = "ANALYSIS"
	OTHER    SegmentType = "OTHER"
)
