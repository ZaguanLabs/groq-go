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
