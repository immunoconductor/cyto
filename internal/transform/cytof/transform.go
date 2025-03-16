package cytof

import (
	"math"
)

// Arcsinh applies an arcsinh transformation to CyTOF data
// a is the offset parameter (typically 0)
// b is the cofactor (typically 1/5 or 0.2 for CyTOF data)
func Arcsinh(value float64, a, b float64) float64 {
	return math.Asinh((value - a) * b)
}
