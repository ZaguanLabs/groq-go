# Groq Go SDK 🚀

[![Go Reference](https://pkg.go.dev/badge/github.com/ZaguanLabs/groq-go.svg)](https://pkg.go.dev/github.com/ZaguanLabs/groq-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZaguanLabs/groq-go)](https://goreportcard.com/report/github.com/ZaguanLabs/groq-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Welcome to the **unofficial** Go client library for the [Groq Cloud API](https://console.groq.com/docs/api-reference)!

This SDK is designed to provide a **strongly-typed**, **idiomatic**, and **robust** experience for Go developers building next-generation AI applications on Groq's blazing-fast LPU™ Inference Engine.

## ✨ v0.3.0 Release - Python SDK 0.37.0 Parity!

**Latest Version:** v0.3.0 (Stable)  
**Status:** ✅ Production Ready  
**Test Coverage:** 73.5%+  
**Quality Grade:** A- (91%)

This release includes:
- 🔄 **Python SDK 0.37.0 Sync** - Full API parity with official Groq Python SDK
- 📊 **Enhanced Usage Stats** - Queue time, reasoning tokens, cached tokens
- 🔧 **MCP Tool Discovery** - Model Context Protocol tool support
- 📝 **JSON Schema Output** - Structured output with JSON Schema validation
- 🎵 **Audio Enhancements** - Sample rate control, URL-based transcription
- 💾 **Cache Statistics** - DRAM/SRAM cache metrics
- ✅ **238 Tests** - Comprehensive test suite with 73.5%+ coverage
- 🏆 **Production Quality** - A- audit grade (91%)

## 🏭 Production Use

This SDK is **running in production** at [Zaguán](https://zaguanai.com), powering real-world AI applications with proven reliability:

- ✅ **OpenWebUI Integration** - Seamless chat functionality with no issues
- ✅ **Qwen-Code with Tools** - Full function and tool calling support working flawlessly
- ✅ **Battle-Tested** - Handling production workloads with confidence

---

## 🌟 Why use this SDK?

- **🎯 Production Ready**: 73.5% test coverage with 140+ comprehensive tests, A- audit grade
- **🤖 Advanced AI**: Compound AI, reasoning models, RAG with citations, web search
- **⚡ Idiomatic Go**: Built with `context`, functional options, and strict typing
- **📦 Complete Coverage**: Chat, Audio, Embeddings, Models, Files, Batches APIs
- **🔄 Streaming First**: Native SSE support with easy-to-use iterators
- **🎛️ Precise Control**: Generic `Optional[T]` types for zero-values vs omitted fields
- **🛡️ Robust**: Exponential backoff retries, rate limit handling, safe error types

## 📦 Installation

```bash
go get github.com/ZaguanLabs/groq-go@v0.3.0
```

**Requirements:**
- Go 1.21 or higher
- Valid Groq API key ([Get one here](https://console.groq.com/keys))

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

### Core APIs
- **Chat Completions** ✅ - Standard and streaming support via SSE (100% tested)
- **Audio** ✅ - Speech generation, transcription, translation (100% tested)
- **Embeddings** ✅ - Vector generation for text (100% tested)
- **Models** ✅ - List and retrieve available models (100% tested)
- **Files** ✅ - Upload, list, and retrieve files (100% tested)
- **Batches** ✅ - Batch processing operations (100% tested)

### Advanced Features
- **Compound AI** - Multi-model orchestration with custom tools
- **Reasoning Models** - Advanced reasoning with configurable effort levels
- **Documents & RAG** - Document context with automatic citations
- **Web Search** - Fine-grained search control with domain filtering
- **Enhanced Streaming** - Complete metadata and usage breakdown
- **Multimodal Content** - Text, image, and document content parts
- **JSON Schema Output** 🆕 - Structured output with JSON Schema validation (v0.37.0)
- **MCP Tool Discovery** 🆕 - Model Context Protocol tool support (v0.37.0)
- **Cache Statistics** 🆕 - DRAM/SRAM cache metrics (v0.37.0)

## 💡 Examples

Check the [groq/examples/](groq/examples/) directory for runnable examples:

- [Chat Completion](groq/examples/chat_completion/main.go) - Basic chat completion
- [Streaming Chat](groq/examples/streaming/main.go) - Streaming responses with SSE
- [Compound AI](groq/examples/compound_ai/main.go) 🆕 - Multi-model workflows with tools
- [Documents & RAG](groq/examples/documents_rag/main.go) 🆕 - Document-based Q&A with citations
- [Reasoning Models](groq/examples/reasoning/main.go) 🆕 - Advanced reasoning capabilities
- [Document Content](groq/examples/document_content/main.go) 🆕 - Multimodal messages with JSON documents

### Quick Example: Compound AI

```go
resp, err := client.Chat.Create(ctx, &types.CreateChatCompletionRequest{
    Model: types.ModelCompoundCustom,
    Messages: []types.ChatCompletionMessage{
        {Role: "user", Content: option.Some("Explain quantum computing")},
    },
    CompoundCustom: option.Ptr(option.Some(types.CompoundCustom{
        AnsweringModel: option.Some("llama-3.3-70b-versatile"),
        ReasoningModel: option.Some("llama-3.1-70b-versatile"),
    })),
    CompoundCustomTools: option.Ptr(option.Some(types.CompoundCustomTools{
        WebSearch: option.Some(true),
    })),
})
```

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

## 📊 Quality & Testing

**v0.3.0 Quality Metrics:**
- ✅ **73.5%+ Test Coverage**
- ✅ **238 Comprehensive Tests**
- ✅ **100% Coverage** on all 6 resource packages
- ✅ **A- Audit Grade** (91%)
- ✅ **Zero Race Conditions**
- ✅ **100% Test Pass Rate**

See [docs/80_PERCENT_FINAL_STATUS.md](docs/80_PERCENT_FINAL_STATUS.md) for detailed coverage report.

## 📖 Documentation

- [Release Notes](docs/RELEASE_v0.3.0.md) - Full v0.3.0 release documentation
- [Audit Report](docs/GROQ_GO_AUDIT_REPORT.md) - Comprehensive code audit
- [Coverage Report](docs/80_PERCENT_FINAL_STATUS.md) - Test coverage details
- [API Reference](https://pkg.go.dev/github.com/ZaguanLabs/groq-go) - Go package documentation
- [Finish Reason Investigation](INVESTIGATION_SUMMARY.md) - Analysis of streaming finish_reason behavior

## 🤝 Contributing

Contributions are strictly encouraged! We love the open source community.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

**Testing Requirements:**
- All new features must include tests
- Maintain or improve test coverage
- All tests must pass (`go test ./...`)
- No race conditions (`go test -race ./...`)

## 📄 License

Distributed under the Apache 2.0 License. See `LICENSE` for more information.
