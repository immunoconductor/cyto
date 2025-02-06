package validator

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/immunoconductor/cyto/fcs/constants"
)

// Only printable ASCII characters in the range: 32-126 (20-7E hex), are valid
func IsValidKeyword(s string) bool {
	const MinPrintableASCII = '\u0020'

	for i := 0; i < len(s); i++ {
		if s[i] < MinPrintableASCII || s[i] >= unicode.MaxASCII {
			return false
		}
	}
	return true
}

func HasRequiredKeywords(t []byte) bool {
	bString := string(t)

	for _, keyword := range constants.TextSegmentKeywords {
		if !strings.Contains(bString, string(keyword)) {
			fmt.Println("missing keyword: ", keyword)
			return false
		}
	}
	return true
}
