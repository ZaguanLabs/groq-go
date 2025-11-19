package chat

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/types"
)

func TestNewStream(t *testing.T) {
	body := strings.NewReader(`data: {"id":"test"}

`)
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(body),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "text/event-stream")

	stream := NewStream[types.ChatCompletionChunk](resp)
	if stream == nil {
		t.Fatal("NewStream returned nil")
	}
	if stream.resp != resp {
		t.Error("response not set correctly")
	}

	stream.Close()
}

func TestStream_Next(t *testing.T) {
	tests := []struct {
		name        string
		sseData     string
		wantChunks  int
		wantErr     error
		errContains string
	}{
		{
			name: "single chunk",
			sseData: `data: {"id":"test","choices":[{"delta":{"content":"Hi"},"index":0}]}

data: [DONE]

`,
			wantChunks: 1,
			wantErr:    io.EOF,
		},
		{
			name: "multiple chunks",
			sseData: `data: {"id":"test","choices":[{"delta":{"content":"Hi"},"index":0}]}

data: {"id":"test","choices":[{"delta":{"content":" there"},"index":0}]}

data: [DONE]

`,
			wantChunks: 2,
			wantErr:    io.EOF,
		},
		{
			name: "chunk with [DONE]",
			sseData: `data: {"id":"test","choices":[{"delta":{"content":"Hi"},"index":0}]}

data: [DONE]

`,
			wantChunks: 1,
			wantErr:    io.EOF,
		},
		{
			name: "[DONE] only",
			sseData: `data: [DONE]

`,
			wantChunks: 0,
			wantErr:    io.EOF,
		},
		{
			name: "invalid JSON",
			sseData: `data: {invalid json}

`,
			wantChunks:  0,
			errContains: "unmarshal",
		},
		{
			name:       "empty stream",
			sseData:    ``,
			wantChunks: 0,
			wantErr:    io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.sseData)
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(body),
				Header:     make(http.Header),
			}
			resp.Header.Set("Content-Type", "text/event-stream")

			stream := NewStream[types.ChatCompletionChunk](resp)
			defer stream.Close()

			chunks := 0
			var lastErr error
			for {
				chunk, err := stream.Next(context.Background())

				if err != nil {
					lastErr = err
					if tt.wantErr != nil && errors.Is(err, tt.wantErr) {
						break
					}
					if tt.errContains != "" && strings.Contains(err.Error(), tt.errContains) {
						break
					}
					if tt.wantErr == nil && tt.errContains == "" {
						t.Fatalf("unexpected error: %v", err)
					}
					break
				}

				if chunk != nil {
					chunks++
				}
			}

			// Verify we got the expected error if specified
			if tt.wantErr != nil && !errors.Is(lastErr, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, lastErr)
			}

			if chunks != tt.wantChunks {
				t.Errorf("got %d chunks, want %d", chunks, tt.wantChunks)
			}
		})
	}
}

func TestStream_NextWithContext(t *testing.T) {
	t.Run("context cancellation", func(t *testing.T) {
		// Create a stream that blocks
		body := io.NopCloser(strings.NewReader(""))
		resp := &http.Response{
			StatusCode: 200,
			Body:       body,
			Header:     make(http.Header),
		}

		stream := NewStream[types.ChatCompletionChunk](resp)
		defer stream.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := stream.Next(ctx)
		if err == nil {
			t.Fatal("expected context cancellation error")
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}

func TestStream_Close(t *testing.T) {
	body := io.NopCloser(strings.NewReader(""))
	resp := &http.Response{
		StatusCode: 200,
		Body:       body,
		Header:     make(http.Header),
	}

	stream := NewStream[types.ChatCompletionChunk](resp)

	err := stream.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Closing again should not panic
	err = stream.Close()
	if err != nil {
		t.Errorf("second Close() error = %v", err)
	}
}

func TestStream_CompleteFlow(t *testing.T) {
	sseData := `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"llama-3.1-8b-instant","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"llama-3.1-8b-instant","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"llama-3.1-8b-instant","choices":[{"index":0,"delta":{"content":" there"},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"llama-3.1-8b-instant","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":"stop"}]}

data: [DONE]

`

	body := strings.NewReader(sseData)
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(body),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "text/event-stream")

	stream := NewStream[types.ChatCompletionChunk](resp)
	defer stream.Close()

	var content strings.Builder
	chunkCount := 0

	for {
		chunk, err := stream.Next(context.Background())
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		chunkCount++
		if len(chunk.Choices) > 0 {
			content.WriteString(chunk.Choices[0].Delta.Content)
		}
	}

	if chunkCount != 4 {
		t.Errorf("got %d chunks, want 4", chunkCount)
	}

	expectedContent := "Hello there!"
	if content.String() != expectedContent {
		t.Errorf("content = %q, want %q", content.String(), expectedContent)
	}
}

func TestStream_WithUsageInfo(t *testing.T) {
	sseData := `data: {"id":"test","choices":[{"delta":{"content":"Hi"},"index":0}]}

data: {"id":"test","choices":[{"delta":{},"index":0,"finish_reason":"stop"}],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}

data: [DONE]

`

	body := strings.NewReader(sseData)
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(body),
		Header:     make(http.Header),
	}

	stream := NewStream[types.ChatCompletionChunk](resp)
	defer stream.Close()

	var lastChunk *types.ChatCompletionChunk
	for {
		chunk, err := stream.Next(context.Background())
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		lastChunk = chunk
	}

	if lastChunk == nil {
		t.Fatal("no chunks received")
	}

	if lastChunk.Usage == nil {
		t.Fatal("expected usage info in final chunk")
	}

	if lastChunk.Usage.TotalTokens != 15 {
		t.Errorf("total tokens = %d, want 15", lastChunk.Usage.TotalTokens)
	}
}
