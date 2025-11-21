# Document Content Example

This example demonstrates how to use the new **document content part** feature added in Groq Python SDK v0.36.0 and now available in groq-go.

## What is Document Content?

Document content allows you to send structured JSON data as part of your chat messages, alongside text and images. This is useful for:

- Sending tabular data for analysis
- Providing structured context (JSON objects)
- Passing configuration or metadata
- Sharing data that models can reference

## Usage

```bash
export GROQ_API_KEY=your_api_key_here
go run main.go
```

## Code Example

```go
docID := "sales-report-2025"

req := &types.CreateChatCompletionRequest{
    Model: "llama-3.3-70b-versatile",
    Messages: []types.ChatCompletionMessageParam{
        {
            Role: types.RoleUser,
            Content: []interface{}{
                types.ContentPartText{
                    Type: "text",
                    Text: "Analyze this sales data:",
                },
                types.ContentPartDocument{
                    Type: "document",
                    Document: types.ContentPartDocument_Document{
                        Data: map[string]interface{}{
                            "sales": []interface{}{
                                map[string]interface{}{
                                    "month":   "January",
                                    "revenue": 125000,
                                    "units":   450,
                                },
                            },
                            "region": "North America",
                        },
                        ID: &docID,
                    },
                },
            },
        },
    },
}
```

## Content Part Types

The SDK supports three types of content parts:

### 1. Text Content
```go
types.ContentPartText{
    Type: "text",
    Text: "Your text here",
}
```

### 2. Image Content
```go
types.ContentPartImage{
    Type: "image_url",
    ImageURL: types.ContentPartImage_ImageURL{
        URL:    "https://example.com/image.jpg",
        Detail: "high", // "auto", "low", or "high"
    },
}
```

### 3. Document Content (NEW)
```go
docID := "optional-id"
types.ContentPartDocument{
    Type: "document",
    Document: types.ContentPartDocument_Document{
        Data: map[string]interface{}{
            // Any JSON-serializable data
            "key": "value",
        },
        ID: &docID, // Optional
    },
}
```

## Backward Compatibility

Simple text messages still work as before:

```go
Messages: []types.ChatCompletionMessageParam{
    {
        Role:    types.RoleUser,
        Content: "Hello, world!", // String content
    },
}
```

## Model Support

Check the Groq API documentation for which models support document content parts. Not all models may support this feature.

## Related Examples

- `../simple_chat/` - Basic text chat
- `../vision/` - Image content (if available)
- `../streaming/` - Streaming responses

## API Reference

- [Groq API Documentation](https://console.groq.com/docs)
- [OpenAI Vision Guide](https://platform.openai.com/docs/guides/vision) (similar multimodal concepts)
