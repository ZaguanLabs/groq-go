package chat

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// mockRequester implements the Requester interface for testing
type mockRequester struct {
	postFunc       func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error
	postStreamFunc func(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error)
}

func (m *mockRequester) Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
	if m.postFunc != nil {
		return m.postFunc(ctx, path, body, result, opts...)
	}
	return nil
}

func (m *mockRequester) PostStream(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error) {
	if m.postStreamFunc != nil {
		return m.postStreamFunc(ctx, path, body, opts...)
	}
	return nil, nil
}

func TestNewCompletions(t *testing.T) {
	mock := &mockRequester{}
	c := NewCompletions(mock)

	if c == nil {
		t.Fatal("NewCompletions returned nil")
	}
	if c.requester != mock {
		t.Error("requester not set correctly")
	}
}

func TestCompletions_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.CreateChatCompletionRequest
		mockResp    *types.ChatCompletion
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "successful request",
			req: &types.CreateChatCompletionRequest{
				Model: "llama-3.1-8b-instant",
				Messages: []types.ChatCompletionMessageParam{
					{Role: types.RoleUser, Content: "Hello"},
				},
			},
			mockResp: &types.ChatCompletion{
				ID:    "test-id",
				Model: "llama-3.1-8b-instant",
				Choices: []types.ChatCompletionChoice{
					{
						Index: 0,
						Message: types.ChatCompletionMessage{
							Role:    types.RoleAssistant,
							Content: "Hi there!",
						},
						FinishReason: types.FinishReasonStop,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error from requester",
			req: &types.CreateChatCompletionRequest{
				Model: "llama-3.1-8b-instant",
				Messages: []types.ChatCompletionMessageParam{
					{Role: types.RoleUser, Content: "Hello"},
				},
			},
			mockErr:     errors.New("network error"),
			wantErr:     true,
			errContains: "network error",
		},
		{
			name: "stream flag set - should error",
			req: &types.CreateChatCompletionRequest{
				Model: "llama-3.1-8b-instant",
				Messages: []types.ChatCompletionMessageParam{
					{Role: types.RoleUser, Content: "Hello"},
				},
				Stream: option.Ptr(option.Some(true)),
			},
			wantErr:     true,
			errContains: "use CreateStream for streaming requests",
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
						// Copy mock response to result
						respBytes, _ := json.Marshal(tt.mockResp)
						json.Unmarshal(respBytes, result)
					}
					return nil
				},
			}

			c := NewCompletions(mock)
			resp, err := c.Create(context.Background(), tt.req)

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

			if resp.ID != tt.mockResp.ID {
				t.Errorf("ID = %q, want %q", resp.ID, tt.mockResp.ID)
			}
		})
	}
}

func TestCompletions_CreateStream(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.CreateChatCompletionRequest
		mockResp    string
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "successful streaming request",
			req: &types.CreateChatCompletionRequest{
				Model: "llama-3.1-8b-instant",
				Messages: []types.ChatCompletionMessageParam{
					{Role: types.RoleUser, Content: "Hello"},
				},
			},
			mockResp: `data: {"id":"test","choices":[{"delta":{"content":"Hi"},"index":0}]}

data: [DONE]

`,
			wantErr: false,
		},
		{
			name: "error from requester",
			req: &types.CreateChatCompletionRequest{
				Model: "llama-3.1-8b-instant",
				Messages: []types.ChatCompletionMessageParam{
					{Role: types.RoleUser, Content: "Hello"},
				},
			},
			mockErr:     errors.New("connection error"),
			wantErr:     true,
			errContains: "connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				postStreamFunc: func(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error) {
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}

					// Create mock HTTP response
					resp := &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(tt.mockResp)),
						Header:     make(http.Header),
					}
					resp.Header.Set("Content-Type", "text/event-stream")
					return resp, nil
				},
			}

			c := NewCompletions(mock)
			stream, err := c.CreateStream(context.Background(), tt.req)

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

			if stream == nil {
				t.Fatal("stream is nil")
			}

			// Verify stream flag was set
			if tt.req.Stream == nil || !tt.req.Stream.IsSet() || !tt.req.Stream.Value {
				t.Error("Stream flag not set to true")
			}

			// Clean up
			stream.Close()
		})
	}
}

func TestCompletions_CreateWithOptions(t *testing.T) {
	mock := &mockRequester{
		postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
			// Verify options were passed
			if len(opts) == 0 {
				t.Error("expected options to be passed")
			}
			return nil
		},
	}

	c := NewCompletions(mock)
	req := &types.CreateChatCompletionRequest{
		Model: "llama-3.1-8b-instant",
		Messages: []types.ChatCompletionMessageParam{
			{Role: types.RoleUser, Content: "Hello"},
		},
	}

	_, err := c.Create(
		context.Background(),
		req,
		option.WithRequestHeader("X-Custom", "test"),
		option.WithIdempotencyKey("test-key"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompletions_CreateWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	mock := &mockRequester{
		postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	c := NewCompletions(mock)
	req := &types.CreateChatCompletionRequest{
		Model: "llama-3.1-8b-instant",
		Messages: []types.ChatCompletionMessageParam{
			{Role: types.RoleUser, Content: "Hello"},
		},
	}

	_, err := c.Create(ctx, req)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestCompletions_Integration(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/openai/v1/chat/completions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Parse request body
		var req types.CreateChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		// Check if streaming
		if req.Stream != nil && req.Stream.IsSet() && req.Stream.Value {
			// Return SSE stream
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`data: {"id":"test","choices":[{"delta":{"content":"Hi"},"index":0}]}

data: [DONE]

`))
			return
		}

		// Return normal response
		resp := types.ChatCompletion{
			ID:    "test-id",
			Model: req.Model,
			Choices: []types.ChatCompletionChoice{
				{
					Index: 0,
					Message: types.ChatCompletionMessage{
						Role:    types.RoleAssistant,
						Content: "Test response",
					},
					FinishReason: types.FinishReasonStop,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create mock requester that uses the test server
	mock := &mockRequester{
		postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
			// Make actual HTTP request to test server
			reqBytes, _ := json.Marshal(body)
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+path, strings.NewReader(string(reqBytes)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			return json.NewDecoder(resp.Body).Decode(result)
		},
		postStreamFunc: func(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error) {
			reqBytes, _ := json.Marshal(body)
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+path, strings.NewReader(string(reqBytes)))
			req.Header.Set("Content-Type", "application/json")
			return http.DefaultClient.Do(req)
		},
	}

	c := NewCompletions(mock)

	t.Run("non-streaming", func(t *testing.T) {
		req := &types.CreateChatCompletionRequest{
			Model: "llama-3.1-8b-instant",
			Messages: []types.ChatCompletionMessageParam{
				{Role: types.RoleUser, Content: "Hello"},
			},
		}

		resp, err := c.Create(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.ID != "test-id" {
			t.Errorf("ID = %q, want %q", resp.ID, "test-id")
		}
		if len(resp.Choices) != 1 {
			t.Errorf("got %d choices, want 1", len(resp.Choices))
		}
	})

	t.Run("streaming", func(t *testing.T) {
		req := &types.CreateChatCompletionRequest{
			Model: "llama-3.1-8b-instant",
			Messages: []types.ChatCompletionMessageParam{
				{Role: types.RoleUser, Content: "Hello"},
			},
		}

		stream, err := c.CreateStream(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer stream.Close()

		// Read first chunk
		chunk, err := stream.Next(context.Background())
		if err != nil {
			t.Fatalf("unexpected error reading chunk: %v", err)
		}

		if chunk.ID != "test" {
			t.Errorf("chunk ID = %q, want %q", chunk.ID, "test")
		}

		// Read [DONE]
		_, err = stream.Next(context.Background())
		if !errors.Is(err, io.EOF) {
			t.Errorf("expected EOF, got %v", err)
		}
	})
}
