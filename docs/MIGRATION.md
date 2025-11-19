# Migration Guide

## Migrating from Python SDK to Go SDK

This guide outlines the key differences and mappings between the official Groq Python SDK and this Go SDK.

### Client Initialization

**Python:**
```python
from groq import Groq

client = Groq(api_key="...")
```

**Go:**
```go
import "github.com/ZaguanLabs/groq-go/groq"

client, err := groq.NewClient(groq.WithAPIKey("..."))
```

### Chat Completions

**Python:**
```python
completion = client.chat.completions.create(
    model="llama3-8b-8192",
    messages=[
        {"role": "user", "content": "Hello"}
    ],
    temperature=0.5
)
print(completion.choices[0].message.content)
```

**Go:**
```go
import (
    "github.com/ZaguanLabs/groq-go/groq/types"
    "github.com/ZaguanLabs/groq-go/groq/option"
)

resp, err := client.Chat.Create(ctx, &types.CreateChatCompletionRequest{
    Model: "llama3-8b-8192",
    Messages: []types.ChatCompletionMessageParam{
        {Role: types.RoleUser, Content: "Hello"},
    },
    Temperature: option.Ptr(option.Some(0.5)),
})
fmt.Println(resp.Choices[0].Message.Content)
```

### Streaming

**Python:**
```python
stream = client.chat.completions.create(..., stream=True)
for chunk in stream:
    print(chunk.choices[0].delta.content or "", end="")
```

**Go:**
```go
stream, err := client.Chat.CreateStream(ctx, &types.CreateChatCompletionRequest{...})
defer stream.Close()

for {
    chunk, err := stream.Next(ctx)
    if err == io.EOF {
        break
    }
    fmt.Print(chunk.Choices[0].Delta.Content)
}
```

### Key Differences

1.  **Context Support**: All Go methods accept a `context.Context` as the first argument for cancellation and timeout control.
2.  **Optional Fields**: Python uses `NotGiven` sentinel or `None`. Go uses `option.Optional[T]` or pointers to handle optional/nullable fields.
3.  **Error Handling**: Go returns errors as the second return value, which must be checked.
4.  **Streaming**: Go uses a separate `CreateStream` method and a `Stream` iterator, whereas Python uses the same `create` method with `stream=True` returning a different type/generator.
