package option

import (
	"encoding/json"
	"testing"
)

type TestStruct struct {
	Field *Optional[string] `json:"field,omitempty"`
	Int   *Optional[int]    `json:"int,omitempty"`
}

// Helper to create pointer to Optional
func Ref[T any](o Optional[T]) *Optional[T] {
	return &o
}

func TestOptional_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    TestStruct
		expected string
	}{
		{
			name:     "both set",
			input:    TestStruct{Field: Ref(Some("foo")), Int: Ref(Some(42))},
			expected: `{"field":"foo","int":42}`,
		},
		{
			name:     "none set (nil pointers)",
			input:    TestStruct{Field: nil, Int: nil},
			expected: `{}`,
		},
		{
			name:     "explicit null (unset optional)",
			input:    TestStruct{Field: Ref(None[string]()), Int: nil},
			expected: `{"field":null}`,
		},
		{
			name:     "mixed",
			input:    TestStruct{Field: Ref(Some("foo")), Int: nil},
			expected: `{"field":"foo"}`,
		},
		{
			name:     "zero values",
			input:    TestStruct{Field: Ref(Some("")), Int: Ref(Some(0))},
			expected: `{"field":"","int":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(b) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(b))
			}
		})
	}
}

func TestOptional_Unmarshal(t *testing.T) {
	input := `{"field":"foo","int":42}`
	var s TestStruct
	if err := json.Unmarshal([]byte(input), &s); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if s.Field == nil || !s.Field.IsSet() || s.Field.Value != "foo" {
		t.Errorf("Expected Field foo, got %v", s.Field)
	}
	if s.Int == nil || !s.Int.IsSet() || s.Int.Value != 42 {
		t.Errorf("Expected Int 42, got %v", s.Int)
	}

	inputNull := `{"field":null}`
	var sNull TestStruct
	if err := json.Unmarshal([]byte(inputNull), &sNull); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if sNull.Field != nil {
		t.Error("Expected Field to be nil for null")
	}
}
