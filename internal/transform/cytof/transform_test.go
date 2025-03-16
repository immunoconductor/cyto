package cytof_test

import (
	"testing"

	"github.com/immunoconductor/cyto/internal/transform/cytof"
)

func TestTransform(t *testing.T) {

	// Apply the arcsinh transformation with standard CyTOF parameters
	transformedData := cytof.Arcsinh(100, 0, 0.2)
	expected := 3.6895038689889055

	if transformedData != expected {
		t.Errorf("expected %v, got %v", expected, transformedData)
	}

}
