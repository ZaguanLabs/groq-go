package sse

import (
	"bufio"
	"io"
	"strings"
)

// Event represents a Server-Sent Event
type Event struct {
	Event string
	Data  string
	ID    string
	Retry int
}

// Decoder handles SSE stream parsing
type Decoder struct {
	// buffer for current event
	event string
	data  []string
	id    string
	retry int
}

// Decode reads events from the reader
func (d *Decoder) Decode(r io.Reader) (<-chan Event, <-chan error) {
	events := make(chan Event)
	errors := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errors)

		scanner := bufio.NewScanner(r)
		// Set a large buffer size to handle large JSON chunks
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024) // 1MB max token size

		for scanner.Scan() {
			line := scanner.Text()

			// Empty line indicates end of event
			if line == "" {
				if d.hasEvent() {
					events <- d.flush()
				}
				continue
			}

			// Parse field
			colon := strings.IndexByte(line, ':')
			if colon == -1 {
				continue // malformed or comment
			}

			// Handle comments
			if colon == 0 {
				continue
			}

			field := line[:colon]
			value := ""
			if colon+1 < len(line) {
				value = strings.TrimPrefix(line[colon+1:], " ")
			}

			switch field {
			case "event":
				d.event = value
			case "data":
				d.data = append(d.data, value)
			case "id":
				if !strings.Contains(value, "\x00") {
					d.id = value
				}
			case "retry":
				// ignored for now
			}
		}

		if err := scanner.Err(); err != nil {
			errors <- err
		}
	}()

	return events, errors
}

func (d *Decoder) hasEvent() bool {
	return d.event != "" || len(d.data) > 0 || d.id != ""
}

func (d *Decoder) flush() Event {
	evt := Event{
		Event: d.event,
		Data:  strings.Join(d.data, "\n"),
		ID:    d.id,
		Retry: d.retry,
	}

	// Reset
	d.event = ""
	d.data = nil // Keep capacity?
	d.id = ""
	d.retry = 0

	return evt
}

// Custom split function to handle SSE double newlines if needed,
// but scanner with default SplitLines works if we treat empty line as separator.
// The spec says events are separated by pair of newlines.
// A single newline separates fields.
// An empty line (length 0 after trim?) is the separator.
