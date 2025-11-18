package chat

import (
	"context"
	"errors"
	"net/http"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// Requester defines the interface for sending requests
type Requester interface {
	Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error
	PostStream(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error)
}

// Completions handles chat completion requests
type Completions struct {
	requester Requester
}

// NewCompletions creates a new Completions service
func NewCompletions(requester Requester) *Completions {
	return &Completions{requester: requester}
}

// Create sends a new chat completion request
func (c *Completions) Create(ctx context.Context, req *types.CreateChatCompletionRequest, opts ...option.RequestOption) (*types.ChatCompletion, error) {
	if req.Stream != nil && req.Stream.IsSet() && req.Stream.Value {
		return nil, errors.New("use CreateStream for streaming requests")
	}

	var result types.ChatCompletion
	err := c.requester.Post(ctx, "/openai/v1/chat/completions", req, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateStream sends a new streaming chat completion request
func (c *Completions) CreateStream(ctx context.Context, req *types.CreateChatCompletionRequest, opts ...option.RequestOption) (*Stream[types.ChatCompletionChunk], error) {
	req.Stream = option.Ptr(option.Some(true))

	resp, err := c.requester.PostStream(ctx, "/openai/v1/chat/completions", req, opts...)
	if err != nil {
		return nil, err
	}

	return NewStream[types.ChatCompletionChunk](resp), nil
}
