package embeddings

import (
	"context"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// Requester defines the interface for sending requests
type Requester interface {
	Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error
}

// Embeddings handles embedding requests
type Embeddings struct {
	requester Requester
}

// New creates a new Embeddings service
func New(requester Requester) *Embeddings {
	return &Embeddings{requester: requester}
}

// Create creates an embedding vector representing the input text
func (e *Embeddings) Create(ctx context.Context, req *types.CreateEmbeddingRequest, opts ...option.RequestOption) (*types.CreateEmbeddingResponse, error) {
	var result types.CreateEmbeddingResponse
	err := e.requester.Post(ctx, "/openai/v1/embeddings", req, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
