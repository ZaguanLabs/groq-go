package batches

import (
	"context"
	"fmt"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// Requester defines the interface for sending requests
type Requester interface {
	Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error
	Get(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
}

// Batches handles batch operations
type Batches struct {
	requester Requester
}

// New creates a new Batches service
func New(requester Requester) *Batches {
	return &Batches{requester: requester}
}

// Create creates a batch
func (b *Batches) Create(ctx context.Context, req *types.CreateBatchRequest, opts ...option.RequestOption) (*types.Batch, error) {
	var result types.Batch
	err := b.requester.Post(ctx, "/openai/v1/batches", req, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Retrieve retrieves a batch
func (b *Batches) Retrieve(ctx context.Context, batchID string, opts ...option.RequestOption) (*types.Batch, error) {
	var result types.Batch
	err := b.requester.Get(ctx, fmt.Sprintf("/openai/v1/batches/%s", batchID), &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Cancel cancels a batch
func (b *Batches) Cancel(ctx context.Context, batchID string, opts ...option.RequestOption) (*types.Batch, error) {
	var result types.Batch
	err := b.requester.Post(ctx, fmt.Sprintf("/openai/v1/batches/%s/cancel", batchID), nil, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// List lists batches
func (b *Batches) List(ctx context.Context, req *types.ListBatchesRequest, opts ...option.RequestOption) (*types.BatchListResponse, error) {
	var result types.BatchListResponse
	// Merge query params from req into opts if needed, or pass req as body?
	// List requests usually use query params.
	// Client.Get uses opts for query params.
	// We need to convert req struct to query params.
	// We can use option.WithRequestQuery manually or helper.

	// The Client.Get signature: Get(ctx, path, result, opts...)
	// It doesn't take a "request struct" for query params automatically.
	// We should convert req to opts.

	if req != nil {
		if req.After != nil && req.After.IsSet() {
			opts = append(opts, option.WithRequestQuery("after", req.After.Value))
		}
		if req.Limit != nil && req.Limit.IsSet() {
			opts = append(opts, option.WithRequestQuery("limit", fmt.Sprintf("%d", req.Limit.Value)))
		}
	}

	err := b.requester.Get(ctx, "/openai/v1/batches", &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
