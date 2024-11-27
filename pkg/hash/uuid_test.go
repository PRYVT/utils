package hash

import (
	"testing"
)

func TestGenerateGUID(t *testing.T) {
	input := "test-input"
	guid := GenerateGUID(input)
	expectedGuid := "ed5d3bd6-3c53-5059-be16-b5dd4b170406"
	if guid.String() != expectedGuid {
		t.Errorf("Expected %v, but got %v", expectedGuid, guid)
	}
}
