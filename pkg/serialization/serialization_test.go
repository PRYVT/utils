package serialization

import (
	"encoding/json"
	"testing"
)

type TestStruct struct {
	Name string
	Age  int
}

func TestUnsafeSerializeAny(t *testing.T) {
	original := TestStruct{Name: "John", Age: 30}
	data := UnsafeSerializeAny(original)

	var result TestStruct
	err := json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	if result.Name != original.Name || result.Age != original.Age {
		t.Errorf("Expected %v, got %v", original, result)
	}
}

func TestUnsafeDeserializeAny(t *testing.T) {
	data := []byte(`{"Name":"John","Age":30}`)
	result := UnsafeDeserializeAny[TestStruct](data)

	if result.Name != "John" || result.Age != 30 {
		t.Errorf("Expected Name: John, Age: 30, got Name: %v, Age: %v", result.Name, result.Age)
	}
}
