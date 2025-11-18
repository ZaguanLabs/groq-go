package groq

import (
	"context"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

func BenchmarkClient_BuildRequest(b *testing.B) {
	c, _ := NewClient(WithAPIKey("bench-key"))
	ctx := context.Background()
	req := &types.CreateChatCompletionRequest{
		Model: "llama3-8b",
		Messages: []types.ChatCompletionMessageParam{
			{Role: types.RoleUser, Content: "Hello world"},
		},
	}
	opts := &option.RequestOptions{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.buildRequest(ctx, "POST", "/openai/v1/chat/completions", req, opts)
	}
}

func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(WithAPIKey("bench-key"))
	}
}
