package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ZaguanLabs/groq-go/groq/internal/sse"
)

// Stream represents a streaming response iterator
type Stream[T any] struct {
	resp   *http.Response
	events <-chan sse.Event
	errors <-chan error
}

// NewStream creates a new stream
func NewStream[T any](resp *http.Response) *Stream[T] {
	d := &sse.Decoder{}
	events, errors := d.Decode(resp.Body)

	return &Stream[T]{
		resp:   resp,
		events: events,
		errors: errors,
	}
}

// Next returns the next item in the stream.
// Returns io.EOF when stream is complete.
func (s *Stream[T]) Next(ctx context.Context) (*T, error) {
	select {
	case event, ok := <-s.events:
		if !ok {
			return nil, io.EOF
		}

		if strings.HasPrefix(event.Data, "[DONE]") {
			// Should we close here? Or wait for channel close?
			// Usually [DONE] is the last message.
			return nil, io.EOF
		}

		var result T
		if err := json.Unmarshal([]byte(event.Data), &result); err != nil {
			return nil, fmt.Errorf("unmarshal SSE data: %w", err)
		}
		return &result, nil

	case err := <-s.errors:
		return nil, err

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Close closes the stream response body
func (s *Stream[T]) Close() error {
	return s.resp.Body.Close()
}
