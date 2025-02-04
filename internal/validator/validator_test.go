package validator_test

import (
	"testing"

	"github.com/immunoconductor/cyto/internal/validator"
)

func TestIsValidKeyword(t *testing.T) {
	// valid print character
	s := "Hello\x7EWorld"
	result := validator.IsValidKeyword(s)
	if result != true {
		t.Errorf("expected keyword to be valid")
	}

	// DEL character - outside range (print character)
	s2 := "Hello\x7FWorld"
	result2 := validator.IsValidKeyword(s2)
	if result2 != false {
		t.Errorf("expected keyword to be invalid")
	}

	// control character
	s3 := "Hello\x1FWorld"
	result3 := validator.IsValidKeyword(s3)
	if result3 != false {
		t.Errorf("expected keyword to be invalid")
	}

}
