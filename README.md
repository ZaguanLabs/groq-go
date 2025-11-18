# Groq Go SDK

Unfficial Go client library for accessing the [Groq API](https://console.groq.com/docs/api-reference). This SDK provides a strongly-typed, idiomatic Go experience for interacting with Groq's LPU™ Inference Engine.

## Installation

```bash
go get github.com/ZaguanLabs/groq-go/groq
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ZaguanLabs/groq-go/groq"
	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

func main() {
	// Initialize the client
	client, err := groq.NewClient(
		groq.WithAPIKey(os.Getenv("GROQ_API_KEY")),
	)
	if err != nil {
		panic(err)
	}

	// Create a chat completion
	resp, err := client.Chat.Create(context.Background(), &types.CreateChatCompletionRequest{
		Model: "llama3-8b-8192",
		Messages: []types.ChatCompletionMessageParam{
			{
				Role:    types.RoleUser,
				Content: "Explain quantum computing in one sentence.",
			},
		},
		Temperature: option.Ptr(option.Some(0.7)),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
```

## Features

- **Chat Completions**: Standard and streaming support via Server-Sent Events (SSE).
- **Audio**: Speech generation, transcription, and translation.
- **Embeddings**: Vector generation for text.
- **Models**: List and retrieve available models.
- **Files**: Upload, list, and retrieve files.
- **Batches**: Batch processing operations.
- **Robustness**: Automatic retries with exponential backoff, context cancellation support.
- **Type Safety**: Comprehensive types for requests and responses, including `Optional[T]` for handling null/omitted fields.

## Examples

Check the [groq/examples/](groq/examples/) directory for runnable examples:

- [Chat Completion](groq/examples/chat_completion/main.go)
- [Streaming Chat](groq/examples/streaming/main.go)

## Advanced Usage

### Streaming

```go
stream, err := client.Chat.CreateStream(ctx, &types.CreateChatCompletionRequest{...})
defer stream.Close()

for {
    chunk, err := stream.Next(ctx)
    if errors.Is(err, io.EOF) {
        break
    }
    if err != nil {
        return err
    }
    fmt.Print(chunk.Choices[0].Delta.Content)
}
```

### Optional Fields

This SDK uses `option.Optional[T]` to distinguish between zero values (e.g., `0`, `""`, `false`) and unset values.

- Use `option.Some(value)` to set a value.
- Use `option.None[T]()` to explicitly send `null` (if supported by API).
- Omit the field to exclude it from the request.
- Use `option.Ptr(option.Some(v))` for pointer fields in request structs.

### Request Options

You can pass per-request options to any method:

```go
client.Chat.Create(ctx, req, 
    option.WithHeader("X-Custom-Header", "value"),
    option.WithRequestQuery("verbose", "true"),
    option.WithIdempotencyKey("unique-key-123"),
)
```

## Project Structure

- `groq/`: Main SDK source code
- `groq/types/`: Request/Response definitions
- `groq/option/`: Functional options and Optional type
- `groq/chat/`, `groq/audio/`, etc.: Resource-specific packages

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
