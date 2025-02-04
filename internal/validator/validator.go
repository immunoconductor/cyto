package validator

import "unicode"

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
