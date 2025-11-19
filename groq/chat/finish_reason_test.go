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

// TestStreamFinishReason_MockAPI tests that the SDK properly receives finish_reason
// This test uses mock data to verify the streaming behavior matches expectations.
func TestStreamFinishReason_MockAPI(t *testing.T) {
	// Mock SSE data that simulates what Groq API should send
	sseData := `data: {"id":"test","choices":[{"delta":{"content":"Hello"},"index":0}]}

data: {"id":"test","choices":[{"delta":{"content":" world"},"index":0}]}

data: {"id":"test","choices":[{"delta":{},"index":0,"finish_reason":"stop"}]}

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

	sawFinishReason := false
	chunkCount := 0
	var lastChunk *types.ChatCompletionChunk

	t.Log("Starting stream iteration...")
	for {
		chunk, err := stream.Next(context.Background())
		if err != nil {
			if errors.Is(err, io.EOF) {
				t.Logf("Received EOF after %d chunks", chunkCount)
				break
			}
			t.Fatalf("unexpected error: %v", err)
		}

		chunkCount++
		lastChunk = chunk

		t.Logf("Chunk %d: model=%q, choices=%d", chunkCount, chunk.Model, len(chunk.Choices))

		for i, choice := range chunk.Choices {
			if choice.FinishReason != "" {
				sawFinishReason = true
				t.Logf("  Choice %d: finish_reason=%q, content=%q",
					i, choice.FinishReason, choice.Delta.Content)
			} else {
				t.Logf("  Choice %d: content=%q", i, choice.Delta.Content)
			}
		}
	}

	if !sawFinishReason {
		t.Errorf("FAIL: Stream ended without receiving finish_reason")
		t.Logf("Total chunks received: %d", chunkCount)
		if lastChunk != nil {
			t.Logf("Last chunk: model=%q, choices=%d", lastChunk.Model, len(lastChunk.Choices))
			for i, choice := range lastChunk.Choices {
				t.Logf("  Last choice %d: finish_reason=%q, content=%q",
					i, choice.FinishReason, choice.Delta.Content)
			}
		}
	} else {
		t.Logf("SUCCESS: Received finish_reason before EOF")
	}
}
