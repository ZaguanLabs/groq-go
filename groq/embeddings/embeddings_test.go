package embeddings

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// mockRequester implements the Requester interface for testing
type mockRequester struct {
	postFunc func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error
}

func (m *mockRequester) Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
	if m.postFunc != nil {
		return m.postFunc(ctx, path, body, result, opts...)
	}
	return nil
}

func TestNew(t *testing.T) {
	mock := &mockRequester{}
	e := New(mock)

	if e == nil {
		t.Fatal("New returned nil")
	}
	if e.requester != mock {
		t.Error("requester not set correctly")
	}
}

func TestEmbeddings_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.CreateEmbeddingRequest
		mockResp    *types.CreateEmbeddingResponse
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "single string input",
			req: &types.CreateEmbeddingRequest{
				Input: "The quick brown fox",
				Model: "nomic-embed-text-v1.5",
			},
			mockResp: &types.CreateEmbeddingResponse{
				Object: "list",
				Model:  "nomic-embed-text-v1.5",
				Data: []types.Embedding{
					{
						Index:     0,
						Object:    "embedding",
						Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
					},
				},
				Usage: types.CompletionUsage{
					PromptTokens: 5,
					TotalTokens:  5,
				},
			},
			wantErr: false,
		},
		{
			name: "multiple string inputs",
			req: &types.CreateEmbeddingRequest{
				Input: []string{"First text", "Second text", "Third text"},
				Model: "nomic-embed-text-v1.5",
			},
			mockResp: &types.CreateEmbeddingResponse{
				Object: "list",
				Model:  "nomic-embed-text-v1.5",
				Data: []types.Embedding{
					{
						Index:     0,
						Object:    "embedding",
						Embedding: []float64{0.1, 0.2, 0.3},
					},
					{
						Index:     1,
						Object:    "embedding",
						Embedding: []float64{0.4, 0.5, 0.6},
					},
					{
						Index:     2,
						Object:    "embedding",
						Embedding: []float64{0.7, 0.8, 0.9},
					},
				},
				Usage: types.CompletionUsage{
					PromptTokens: 15,
					TotalTokens:  15,
				},
			},
			wantErr: false,
		},
		{
			name: "with encoding format",
			req: &types.CreateEmbeddingRequest{
				Input:          "Test text",
				Model:          "nomic-embed-text-v1.5",
				EncodingFormat: option.Ptr(option.Some("float")),
			},
			mockResp: &types.CreateEmbeddingResponse{
				Object: "list",
				Model:  "nomic-embed-text-v1.5",
				Data: []types.Embedding{
					{
						Index:     0,
						Object:    "embedding",
						Embedding: []float64{0.1, 0.2},
					},
				},
				Usage: types.CompletionUsage{
					PromptTokens: 2,
					TotalTokens:  2,
				},
			},
			wantErr: false,
		},
		{
			name: "with user identifier",
			req: &types.CreateEmbeddingRequest{
				Input: "User text",
				Model: "nomic-embed-text-v1.5",
				User:  "user-123",
			},
			mockResp: &types.CreateEmbeddingResponse{
				Object: "list",
				Model:  "nomic-embed-text-v1.5",
				Data: []types.Embedding{
					{
						Index:     0,
						Object:    "embedding",
						Embedding: []float64{0.5},
					},
				},
				Usage: types.CompletionUsage{
					PromptTokens: 2,
					TotalTokens:  2,
				},
			},
			wantErr: false,
		},
		{
			name: "error from requester",
			req: &types.CreateEmbeddingRequest{
				Input: "Test",
				Model: "nomic-embed-text-v1.5",
			},
			mockErr:     errors.New("network error"),
			wantErr:     true,
			errContains: "network error",
		},
		{
			name: "empty embedding response",
			req: &types.CreateEmbeddingRequest{
				Input: "",
				Model: "nomic-embed-text-v1.5",
			},
			mockResp: &types.CreateEmbeddingResponse{
				Object: "list",
				Model:  "nomic-embed-text-v1.5",
				Data:   []types.Embedding{},
				Usage: types.CompletionUsage{
					PromptTokens: 0,
					TotalTokens:  0,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
					if tt.mockErr != nil {
						return tt.mockErr
					}
					if tt.mockResp != nil {
						respBytes, _ := json.Marshal(tt.mockResp)
						json.Unmarshal(respBytes, result)
					}
					// Verify path
					if path != "/openai/v1/embeddings" {
						t.Errorf("unexpected path: %s", path)
					}
					// Verify request body
					if req, ok := body.(*types.CreateEmbeddingRequest); ok {
						if req.Model != tt.req.Model {
							t.Errorf("model = %q, want %q", req.Model, tt.req.Model)
						}
					}
					return nil
				},
			}

			e := New(mock)
			resp, err := e.Create(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp == nil {
				t.Fatal("response is nil")
			}

			if resp.Model != tt.mockResp.Model {
				t.Errorf("model = %q, want %q", resp.Model, tt.mockResp.Model)
			}

			if len(resp.Data) != len(tt.mockResp.Data) {
				t.Errorf("got %d embeddings, want %d", len(resp.Data), len(tt.mockResp.Data))
			}

			// Verify embedding dimensions
			for i, emb := range resp.Data {
				if len(emb.Embedding) != len(tt.mockResp.Data[i].Embedding) {
					t.Errorf("embedding[%d] has %d dimensions, want %d",
						i, len(emb.Embedding), len(tt.mockResp.Data[i].Embedding))
				}
			}
		})
	}
}

func TestEmbeddings_CreateWithOptions(t *testing.T) {
	mock := &mockRequester{
		postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
			if len(opts) == 0 {
				t.Error("expected options to be passed")
			}
			return nil
		},
	}

	e := New(mock)
	req := &types.CreateEmbeddingRequest{
		Input: "Test",
		Model: "nomic-embed-text-v1.5",
	}

	_, err := e.Create(
		context.Background(),
		req,
		option.WithRequestHeader("X-Custom", "test"),
		option.WithIdempotencyKey("test-key"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmbeddings_CreateWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	mock := &mockRequester{
		postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	e := New(mock)
	req := &types.CreateEmbeddingRequest{
		Input: "Test",
		Model: "nomic-embed-text-v1.5",
	}

	_, err := e.Create(ctx, req)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestEmbeddings_LargeEmbedding(t *testing.T) {
	// Test with large embedding vector (e.g., 768 dimensions)
	dimensions := 768
	embedding := make([]float64, dimensions)
	for i := 0; i < dimensions; i++ {
		embedding[i] = float64(i) / float64(dimensions)
	}

	mock := &mockRequester{
		postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
			resp := &types.CreateEmbeddingResponse{
				Object: "list",
				Model:  "nomic-embed-text-v1.5",
				Data: []types.Embedding{
					{
						Index:     0,
						Object:    "embedding",
						Embedding: embedding,
					},
				},
				Usage: types.CompletionUsage{
					PromptTokens: 10,
					TotalTokens:  10,
				},
			}
			respBytes, _ := json.Marshal(resp)
			json.Unmarshal(respBytes, result)
			return nil
		},
	}

	e := New(mock)
	req := &types.CreateEmbeddingRequest{
		Input: "Test with large embedding",
		Model: "nomic-embed-text-v1.5",
	}

	resp, err := e.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 embedding, got %d", len(resp.Data))
	}

	if len(resp.Data[0].Embedding) != dimensions {
		t.Errorf("embedding has %d dimensions, want %d", len(resp.Data[0].Embedding), dimensions)
	}
}

func TestEmbeddings_BatchProcessing(t *testing.T) {
	// Test with batch of inputs
	inputs := make([]string, 100)
	for i := 0; i < 100; i++ {
		inputs[i] = "Text " + string(rune(i))
	}

	embeddings := make([]types.Embedding, 100)
	for i := 0; i < 100; i++ {
		embeddings[i] = types.Embedding{
			Index:     i,
			Object:    "embedding",
			Embedding: []float64{float64(i) * 0.1},
		}
	}

	mock := &mockRequester{
		postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
			resp := &types.CreateEmbeddingResponse{
				Object: "list",
				Model:  "nomic-embed-text-v1.5",
				Data:   embeddings,
				Usage: types.CompletionUsage{
					PromptTokens: 200,
					TotalTokens:  200,
				},
			}
			respBytes, _ := json.Marshal(resp)
			json.Unmarshal(respBytes, result)
			return nil
		},
	}

	e := New(mock)
	req := &types.CreateEmbeddingRequest{
		Input: inputs,
		Model: "nomic-embed-text-v1.5",
	}

	resp, err := e.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 100 {
		t.Errorf("got %d embeddings, want 100", len(resp.Data))
	}

	// Verify indices are correct
	for i, emb := range resp.Data {
		if emb.Index != i {
			t.Errorf("embedding[%d] has index %d, want %d", i, emb.Index, i)
		}
	}
}
