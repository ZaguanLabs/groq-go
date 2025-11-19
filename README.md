# Groq Go SDK 🚀

[![Go Reference](https://pkg.go.dev/badge/github.com/ZaguanLabs/groq-go.svg)](https://pkg.go.dev/github.com/ZaguanLabs/groq-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZaguanLabs/groq-go)](https://goreportcard.com/report/github.com/ZaguanLabs/groq-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Welcome to the **unofficial** Go client library for the [Groq Cloud API](https://console.groq.com/docs/api-reference)!

This SDK is designed to provide a **strongly-typed**, **idiomatic**, and **robust** experience for Go developers building next-generation AI applications on Groq's blazing-fast LPU™ Inference Engine.

> ⚠️ **Note**: This project is currently in **Beta (v0.2.0-beta)**. The SDK is feature-complete with full Python SDK v0.35.0 parity. APIs are stabilizing as we approach 1.0.

---

## 🌟 Why use this SDK?

- **Idiomatic Go**: Built with `context`, functional options, and strict typing in mind.
- **Complete Coverage**: Supports Chat, Audio, Embeddings, Models, Files, and Batches.
- **Production Ready**: Built-in exponential backoff retries, rate limit handling, and safe error types.
- **Streaming First**: Native support for Server-Sent Events (SSE) with easy-to-use iterators.
- **Precise Control**: Generic `Optional[T]` types allow you to distinguish between zero-values and omitted fields.

## 📦 Installation

```bash
go get github.com/ZaguanLabs/groq-go/groq
```

## 🚀 Quickstart

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

## 📚 Supported APIs

- **Chat Completions**: Standard and streaming support via Server-Sent Events (SSE).
- **Compound AI**: Multi-model orchestration with custom tool configuration.
- **Documents & RAG**: Provide documents as context with citation support.
- **Reasoning Models**: Advanced reasoning with configurable output formats.
- **Audio**: Speech generation, transcription, and translation.
- **Embeddings**: Vector generation for text.
- **Models**: List and retrieve available models.
- **Files**: Upload, list, and retrieve files.
- **Batches**: Batch processing operations.

## 💡 Examples

Check the [groq/examples/](groq/examples/) directory for runnable examples:

- [Chat Completion](groq/examples/chat_completion/main.go) - Basic chat completion
- [Streaming Chat](groq/examples/streaming/main.go) - Streaming responses
- [Compound AI](groq/examples/compound_ai/main.go) - Multi-model workflows with tools
- [Documents & RAG](groq/examples/documents_rag/main.go) - Document-based context with citations
- [Reasoning Models](groq/examples/reasoning/main.go) - Advanced reasoning capabilities

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

## 🤝 Contributing

Contributions are strictly encouraged! We love the open source community.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

Distributed under the Apache 2.0 License. See `LICENSE` for more information.
