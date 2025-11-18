package chat

import (
	"context"

	"github.com/zaguan/groq-go/option"
)

// Requester defines the interface for sending requests
type Requester interface {
	Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error
	// PostStream will be added later
}

// Completions handles chat completion requests
type Completions struct {
	requester Requester
}

// NewCompletions creates a new Completions service
func NewCompletions(requester Requester) *Completions {
	return &Completions{requester: requester}
}
