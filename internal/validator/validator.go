package validator

import (
	"fmt"
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
	// also add logic to check that it begins with '$'
	return true
}

func HasRequiredKeywords(keywords map[string]string) bool {
	for _, keyword := range constants.TextSegmentKeywords {
		_, exists := keywords[string(keyword)]
		if !exists {
			fmt.Println("missing keyword: ", keyword)
			return false
		}
	}
	return true
}
