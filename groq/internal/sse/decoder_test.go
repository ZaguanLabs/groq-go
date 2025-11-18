package sse

import (
	"strings"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Event
	}{
		{
			name:  "single event",
			input: "data: hello\n\n",
			expected: []Event{
				{Data: "hello"},
			},
		},
		{
			name:  "multiple events",
			input: "data: first\n\ndata: second\n\n",
			expected: []Event{
				{Data: "first"},
				{Data: "second"},
			},
		},
		{
			name:  "event with id and type",
			input: "event: update\nid: 1\ndata: payload\n\n",
			expected: []Event{
				{Event: "update", ID: "1", Data: "payload"},
			},
		},
		{
			name:  "multi-line data",
			input: "data: line1\ndata: line2\n\n",
			expected: []Event{
				{Data: "line1\nline2"},
			},
		},
		{
			name:  "comments ignored",
			input: ": comment\ndata: real\n\n",
			expected: []Event{
				{Data: "real"},
			},
		},
		{
			name:  "no space after colon",
			input: "data:nospace\n\n",
			expected: []Event{
				{Data: "nospace"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{}
			r := strings.NewReader(tt.input)
			events, _ := d.Decode(r)

			var got []Event
			for evt := range events {
				got = append(got, evt)
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("Expected %d events, got %d", len(tt.expected), len(got))
			}

			for i, e := range tt.expected {
				g := got[i]
				if g.Data != e.Data || g.Event != e.Event || g.ID != e.ID {
					t.Errorf("Event %d mismatch: expected %+v, got %+v", i, e, g)
				}
			}
		})
	}
}
