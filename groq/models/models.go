package models

import (
	"context"
	"fmt"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// Requester defines the interface for sending requests
type Requester interface {
	Get(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
	Delete(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
}

// Models handles model requests
type Models struct {
	requester Requester
}

// New creates a new Models service
func New(requester Requester) *Models {
	return &Models{requester: requester}
}

// List lists the currently available models
func (m *Models) List(ctx context.Context, opts ...option.RequestOption) (*types.ModelListResponse, error) {
	var result types.ModelListResponse
	err := m.requester.Get(ctx, "/openai/v1/models", &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Retrieve retrieves a model instance, providing basic information about the model such as the owner and permissioning
func (m *Models) Retrieve(ctx context.Context, modelID string, opts ...option.RequestOption) (*types.Model, error) {
	var result types.Model
	err := m.requester.Get(ctx, fmt.Sprintf("/openai/v1/models/%s", modelID), &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a model
func (m *Models) Delete(ctx context.Context, modelID string, opts ...option.RequestOption) (*types.ModelDeleted, error) {
	var result types.ModelDeleted
	err := m.requester.Delete(ctx, fmt.Sprintf("/openai/v1/models/%s", modelID), &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
