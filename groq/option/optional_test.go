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

func TestOptional_IsSet(t *testing.T) {
	tests := []struct {
		name    string
		opt     Optional[string]
		wantSet bool
	}{
		{"Some value", Some("test"), true},
		{"None", None[string](), false},
		{"Zero value", Optional[string]{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.IsSet(); got != tt.wantSet {
				t.Errorf("IsSet() = %v, want %v", got, tt.wantSet)
			}
		})
	}
}

func TestOptional_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		opt      Optional[int]
		wantZero bool
	}{
		{"Some value", Some(42), false},
		{"Some zero", Some(0), false},
		{"None", None[int](), true},
		{"Zero value", Optional[int]{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.IsZero(); got != tt.wantZero {
				t.Errorf("IsZero() = %v, want %v", got, tt.wantZero)
			}
		})
	}
}

func TestOptional_Ptr(t *testing.T) {
	opt := Some("test")
	ptr := Ptr(opt)

	if ptr == nil {
		t.Fatal("Ptr returned nil")
	}
	if !ptr.IsSet() {
		t.Error("Ptr value not set")
	}
	if ptr.Value != "test" {
		t.Errorf("Ptr value = %q, want %q", ptr.Value, "test")
	}
}

func TestOptional_MarshalTypes(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"string", Some("hello"), `"hello"`},
		{"int", Some(42), `42`},
		{"float", Some(3.14), `3.14`},
		{"bool true", Some(true), `true`},
		{"bool false", Some(false), `false`},
		{"null string", None[string](), `null`},
		{"null int", None[int](), `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(b) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(b))
			}
		})
	}
}

func TestOptional_UnmarshalTypes(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantSet   bool
		wantValue interface{}
	}{
		{"string", `"hello"`, true, "hello"},
		{"int", `42`, true, 42},
		{"float", `3.14`, true, 3.14},
		{"bool", `true`, true, true},
		{"null", `null`, false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.wantValue.(type) {
			case string:
				var opt Optional[string]
				if err := json.Unmarshal([]byte(tt.input), &opt); err != nil {
					t.Fatalf("Unmarshal error: %v", err)
				}
				if opt.IsSet() != tt.wantSet {
					t.Errorf("IsSet() = %v, want %v", opt.IsSet(), tt.wantSet)
				}
				if tt.wantSet && opt.Value != v {
					t.Errorf("Value = %v, want %v", opt.Value, v)
				}
			case int:
				var opt Optional[int]
				if err := json.Unmarshal([]byte(tt.input), &opt); err != nil {
					t.Fatalf("Unmarshal error: %v", err)
				}
				if opt.IsSet() != tt.wantSet {
					t.Errorf("IsSet() = %v, want %v", opt.IsSet(), tt.wantSet)
				}
				if tt.wantSet && opt.Value != v {
					t.Errorf("Value = %v, want %v", opt.Value, v)
				}
			case float64:
				var opt Optional[float64]
				if err := json.Unmarshal([]byte(tt.input), &opt); err != nil {
					t.Fatalf("Unmarshal error: %v", err)
				}
				if opt.IsSet() != tt.wantSet {
					t.Errorf("IsSet() = %v, want %v", opt.IsSet(), tt.wantSet)
				}
				if tt.wantSet && opt.Value != v {
					t.Errorf("Value = %v, want %v", opt.Value, v)
				}
			case bool:
				var opt Optional[bool]
				if err := json.Unmarshal([]byte(tt.input), &opt); err != nil {
					t.Fatalf("Unmarshal error: %v", err)
				}
				if opt.IsSet() != tt.wantSet {
					t.Errorf("IsSet() = %v, want %v", opt.IsSet(), tt.wantSet)
				}
				if tt.wantSet && opt.Value != v {
					t.Errorf("Value = %v, want %v", opt.Value, v)
				}
			case nil:
				var opt Optional[string]
				if err := json.Unmarshal([]byte(tt.input), &opt); err != nil {
					t.Fatalf("Unmarshal error: %v", err)
				}
				if opt.IsSet() {
					t.Error("Expected IsSet() = false for null")
				}
			}
		})
	}
}
