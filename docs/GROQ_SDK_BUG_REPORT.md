# Groq SDK Bug Report: Missing finish_reason in Streaming Responses

**Repository**: https://github.com/ZaguanLabs/groq-go  
**Version**: v0.2.0  
**Date**: November 19, 2025  
**Severity**: High (breaks OpenAI-compatible clients)

## Summary

The Groq Go SDK's streaming implementation does not provide `finish_reason` in any chunks before returning `io.EOF`, causing OpenAI-compatible clients to fail with "Model stream ended without a finish reason" errors.

## Environment

- **SDK Version**: groq-go v0.2.0
- **Go Version**: 1.24.6
- **Models Tested**: 
  - `moonshotai/kimi-k2-instruct-0905`
  - `qwen/qwen3-32b` (Qwen-Code)
- **Client**: Zaguan CoreX (OpenAI-compatible proxy)

## Expected Behavior

According to OpenAI's streaming specification, the final chunk should contain:
```json
{
  "choices": [{
    "index": 0,
    "delta": {},
    "finish_reason": "stop"  // ← Should be present
  }]
}
```

## Actual Behavior

The SDK's `stream.Next()` returns `io.EOF` without ever providing a chunk with `finish_reason` set.

### Observed Stream Sequence

```
1. Multiple chunks with content (finish_reason: "")
2. Empty chunk (model: "", choices: 0, usage: false)
3. io.EOF returned
4. No chunk with finish_reason ever received
```

## Reproduction

### Code

```go
stream, err := client.Chat.CreateStream(ctx, &types.CreateChatCompletionRequest{
    Model: "moonshotai/kimi-k2-instruct-0905",
    Messages: []types.ChatCompletionMessageParam{
        {Role: types.RoleUser, Content: "Hello"},
    },
})

sawFinishReason := false
for {
    chunk, err := stream.Next(ctx)
    if err != nil {
        if errors.Is(err, io.EOF) {
            fmt.Printf("EOF reached. Saw finish reason: %v\n", sawFinishReason)
            break
        }
        return err
    }
    
    for _, choice := range chunk.Choices {
        if choice.FinishReason != "" {
            sawFinishReason = true
            fmt.Printf("Finish reason: %s\n", choice.FinishReason)
        }
    }
}
```

### Output

```
EOF reached. Saw finish reason: false
```

## Debug Logs from Production

```json
{"msg":"Groq SDK stream chunk received","model":"moonshotai/kimi-k2-instruct-0905","choices_count":1,"has_usage":false}
{"msg":"Groq SDK stream chunk received","model":"moonshotai/kimi-k2-instruct-0905","choices_count":1,"has_usage":false}
{"msg":"Groq SDK stream chunk received","model":"moonshotai/kimi-k2-instruct-0905","choices_count":1,"has_usage":false}
...
{"msg":"Groq SDK stream chunk received","model":"","choices_count":0,"has_usage":false}
{"msg":"Skipping empty chunk (no choices, no usage)"}
{"msg":"Skipping nil chunk from Groq SDK stream"}
{"msg":"Chat request succeeded","duration":1284179523}
```

**Note**: No log message showing "Groq stream chunk has finish reason" ever appears.

## Root Cause Analysis

### Hypothesis 1: SDK Issue

Looking at `groq/chat/stream.go:42-46`:

```go
if strings.HasPrefix(event.Data, "[DONE]") {
    // Should we close here? Or wait for channel close?
    // Usually [DONE] is the last message.
    return nil, io.EOF
}
```

The SDK returns `io.EOF` immediately when it sees `[DONE]`, potentially before processing a final chunk with `finish_reason`.

### Hypothesis 2: API Issue

The Groq API might not be sending `finish_reason` in streaming mode at all, which would be a server-side bug.

### Hypothesis 3: Parsing Issue

The SDK might be receiving chunks with `finish_reason` but not properly parsing or forwarding them.

## Impact

### Affected Clients

- **OpenWebUI**: Shows "API Error: Model stream ended without a finish reason"
- **Any OpenAI-compatible client**: Expects finish_reason in final chunk
- **Langchain/LlamaIndex**: May fail validation checks

### Workaround Required

We've implemented a workaround in Zaguan CoreX:

```go
sawFinishReason := false
for {
    chunk, err := stream.Next(ctx)
    if err != nil {
        if errors.Is(err, io.EOF) {
            if !sawFinishReason {
                // Synthesize final chunk with finish_reason
                outputChan <- core.ChatResponseChunk{
                    Model: modelName,
                    Choices: []core.Choice{{
                        Index:        0,
                        FinishReason: "stop",
                        Delta:        core.ChatMessage{},
                    }},
                }
            }
            return
        }
    }
    
    for _, choice := range chunk.Choices {
        if choice.FinishReason != "" {
            sawFinishReason = true
        }
    }
}
```

## Comparison with Other SDKs

### OpenAI Go SDK (go-openai)

The OpenAI SDK properly provides finish_reason in the final chunk before EOF:

```go
// Final chunk from OpenAI
{
  "choices": [{
    "index": 0,
    "delta": {},
    "finish_reason": "stop"
  }],
  "usage": {...}
}
```

### Expected Groq Behavior

Groq should match OpenAI's behavior for compatibility.

## Debugging Steps Taken

1. ✅ Added detailed logging for every chunk received
2. ✅ Logged chunk properties (model, choices_count, has_usage)
3. ✅ Logged when finish_reason is detected (never triggered)
4. ✅ Tested with multiple models (kimi-k2, qwen-code)
5. ✅ Compared with old OpenAI-compatible implementation
6. ✅ Verified types.ChatCompletionChunkChoice has FinishReason field

## Proposed Solutions

### Option 1: SDK Fix (Recommended)

Ensure the SDK processes all chunks before returning EOF:

```go
// In stream.go
if strings.HasPrefix(event.Data, "[DONE]") {
    // Process any pending chunks first
    // Then return EOF
    return nil, io.EOF
}
```

### Option 2: API Fix

If the issue is server-side, Groq API should send a final chunk with `finish_reason` before `[DONE]`.

### Option 3: Documentation

If this is intentional, document that clients must synthesize finish_reason on EOF.

## Test Case

```go
func TestStreamFinishReason(t *testing.T) {
    client, _ := groq.NewClient(groq.WithAPIKey(os.Getenv("GROQ_API_KEY")))
    
    stream, err := client.Chat.CreateStream(context.Background(), 
        &types.CreateChatCompletionRequest{
            Model: "llama-3.3-70b-versatile",
            Messages: []types.ChatCompletionMessageParam{
                {Role: types.RoleUser, Content: "Say 'test'"},
            },
        })
    require.NoError(t, err)
    defer stream.Close()
    
    sawFinishReason := false
    for {
        chunk, err := stream.Next(context.Background())
        if err != nil {
            if errors.Is(err, io.EOF) {
                break
            }
            require.NoError(t, err)
        }
        
        for _, choice := range chunk.Choices {
            if choice.FinishReason != "" {
                sawFinishReason = true
            }
        }
    }
    
    assert.True(t, sawFinishReason, "Stream should provide finish_reason before EOF")
}
```

## Additional Information

### Related Issues

- Similar issue in other streaming SDKs?
- OpenAI compatibility requirements

### Logs Available

We have production logs showing:
- Every chunk received from the SDK
- Chunk properties (model, choices, usage)
- Timing information
- No finish_reason ever detected

### Contact

For more information or to provide test data, contact:
- **Project**: Zaguan CoreX
- **Issue Tracker**: https://git.kekepower.com/ZaguanAI/zaguancorex

## Checklist for SDK Maintainers

- [ ] Verify Groq API actually sends finish_reason in streaming
- [ ] Check if SDK is properly parsing finish_reason field
- [ ] Ensure finish_reason is forwarded before EOF
- [ ] Add test case for finish_reason in streaming
- [ ] Update documentation if behavior is intentional
- [ ] Consider OpenAI compatibility requirements

## Workaround Status

✅ **Workaround implemented** in Zaguan CoreX v0.37.0-beta7  
⚠️ **Not ideal**: Clients shouldn't need to synthesize finish reasons  
🔧 **Proper fix needed**: SDK or API should provide finish_reason

---

**Priority**: High  
**Affects**: All streaming requests  
**Workaround**: Available but not ideal  
**Fix Required**: SDK or API update

**Thank you for maintaining this SDK!** We're happy to provide more debug information or test cases as needed.

---

## Investigation Results (November 19, 2025)

### Conclusion: This is an API Issue, Not an SDK Bug ✅

After thorough investigation, we've determined:

1. **The SDK correctly handles `finish_reason`** when the API sends it
   - Test added: `groq/chat/finish_reason_test.go`
   - Diagnostic tool added: `groq/examples/debug_finish_reason/main.go`
   - The SDK properly processes chunks with `finish_reason` before `[DONE]`

2. **The Groq API is not sending `finish_reason`** for certain models
   - Models affected: `moonshotai/kimi-k2-instruct-0905`, `qwen/qwen3-32b`
   - The API sends an empty chunk (0 choices) before `[DONE]`
   - This violates OpenAI streaming specification

3. **The Python SDK has the same behavior**
   - Both SDKs return EOF when they see `[DONE]`
   - Neither SDK synthesizes finish_reason (correct behavior)

### Why the SDK Should NOT Synthesize finish_reason

- **Transparency**: The SDK should return exactly what the API sends
- **Debugging**: Users need to know the actual API behavior
- **Correctness**: The SDK can't know what finish_reason should be

### Recommendation

**For Groq API Team**: Update the API to send a final chunk with `finish_reason` before `[DONE]` for OpenAI compatibility.

**For SDK Users**: Implement the workaround in your application layer (see `FINISH_REASON_ANALYSIS.md`).

### Files Added

- `FINISH_REASON_ANALYSIS.md` - Complete analysis and recommendations
- `groq/chat/finish_reason_test.go` - Test proving SDK works correctly
- `groq/examples/debug_finish_reason/main.go` - Diagnostic tool for real API testing
