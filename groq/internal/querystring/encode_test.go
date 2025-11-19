package querystring

import (
	"testing"
	"time"
)

func TestStringify(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil",
			input:    nil,
			expected: "",
		},
		{
			name: "simple map",
			input: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			expected: "baz=qux&foo=bar",
		},
		{
			name: "map with int",
			input: map[string]interface{}{
				"limit": 10,
				"page":  1,
			},
			expected: "limit=10&page=1",
		},
		{
			name: "map with bool",
			input: map[string]bool{
				"active": true,
			},
			expected: "active=true",
		},
		{
			name: "map with slice (comma array)",
			input: map[string][]string{
				"ids": {"1", "2", "3"},
			},
			expected: "ids=1%2C2%2C3", // Comma is encoded as %2C
		},
		{
			name: "map with interface slice",
			input: map[string]interface{}{
				"ids": []int{1, 2, 3},
			},
			expected: "ids=1%2C2%2C3",
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"filter": map[string]string{
					"name": "test",
				},
			},
			expected: "filter%5Bname%5D=test", // filter[name]=test
		},
		{
			name: "map with time",
			input: map[string]time.Time{
				"created": time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			expected: "created=2023-01-01T12%3A00%3A00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Stringify(tt.input)
			if err != nil {
				t.Fatalf("Stringify() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("Stringify() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestStringify_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   map[string]string{"key": ""},
			wantErr: false,
		},
		{
			name:    "zero int",
			input:   map[string]int{"count": 0},
			wantErr: false,
		},
		{
			name:    "false bool",
			input:   map[string]bool{"active": false},
			wantErr: false,
		},
		{
			name:    "empty slice",
			input:   map[string][]string{"ids": {}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Stringify(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stringify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringify_ComplexNesting(t *testing.T) {
	input := map[string]interface{}{
		"filter": map[string]interface{}{
			"user": map[string]string{
				"name": "john",
			},
		},
	}

	result, err := Stringify(input)
	if err != nil {
		t.Fatalf("Stringify() error = %v", err)
	}

	// Should encode nested structure
	if result == "" {
		t.Error("Expected non-empty result for nested map")
	}
}

func TestStringify_SpecialCharacters(t *testing.T) {
	input := map[string]string{
		"query": "hello world",
		"path":  "/api/v1",
		"email": "test@example.com",
	}

	result, err := Stringify(input)
	if err != nil {
		t.Fatalf("Stringify() error = %v", err)
	}

	// Spaces should be encoded
	if !contains(result, "hello+world") && !contains(result, "hello%20world") {
		t.Error("Expected spaces to be encoded")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
