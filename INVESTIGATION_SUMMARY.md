# Finish Reason Investigation Summary

**Date**: November 19, 2025  
**Issue**: Missing `finish_reason` in streaming responses  
**Status**: ✅ RESOLVED - Confirmed as API issue, not SDK bug

## Quick Summary

The Groq Go SDK **correctly handles `finish_reason`** when the API sends it. The issue is that **certain Groq API models don't send `finish_reason` before `[DONE]`**, which breaks OpenAI compatibility.

## What We Did

### 1. Analyzed the Bug Report ✅

Reviewed `docs/GROQ_SDK_BUG_REPORT.md` which showed:
- OpenWebUI failing with "Model stream ended without a finish reason"
- Models affected: `moonshotai/kimi-k2-instruct-0905`, `qwen/qwen3-32b`
- Logs showing empty chunks (0 choices) before EOF

### 2. Examined SDK Implementation ✅

Reviewed `groq/chat/stream.go`:
```go
if strings.HasPrefix(event.Data, "[DONE]") {
    return nil, io.EOF
}
```

**Finding**: SDK correctly returns EOF when it sees `[DONE]` - this matches OpenAI behavior.

### 3. Created Test Cases ✅

Added `groq/chat/finish_reason_test.go`:
```go
func TestStreamFinishReason_MockAPI(t *testing.T) {
    // Tests that SDK correctly processes finish_reason when API sends it
}
```

**Result**: ✓ PASS - SDK correctly receives and processes `finish_reason`

### 4. Built Diagnostic Tool ✅

Created `groq/examples/debug_finish_reason/main.go`:
- Streams a simple request to Groq API
- Logs every chunk received
- Reports whether `finish_reason` was received before EOF

### 5. Compared with Python SDK ✅

Reviewed `docs/groq-python-0.35.0/src/groq/_streaming.py`:
```python
if sse.data.startswith("[DONE]"):
    break
```

**Finding**: Python SDK has identical behavior - both SDKs are correct.

## Root Cause

**The Groq API is not sending `finish_reason` for certain models.**

### Expected Behavior (OpenAI-compatible):
```
data: {"choices":[{"delta":{"content":"Hello"}}]}
data: {"choices":[{"delta":{"content":" world"}}]}
data: {"choices":[{"delta":{},"finish_reason":"stop"}]}  ← Should be here
data: [DONE]
```

### Actual Behavior (Groq API):
```
data: {"choices":[{"delta":{"content":"Hello"}}]}
data: {"choices":[{"delta":{"content":" world"}}]}
data: {"choices":[]}  ← Empty choices array, no finish_reason
data: [DONE]
```

## Why SDK Should NOT Fix This

1. **Transparency**: SDK should return exactly what the API sends
2. **Debugging**: Users need to know the actual API behavior
3. **Correctness**: SDK can't know what `finish_reason` should be (`stop`, `length`, `tool_calls`, etc.)
4. **Consistency**: Python SDK doesn't synthesize it either

## Solution

### For Groq API Team

Update the API to send a final chunk with `finish_reason` before `[DONE]`:

```json
data: {"id":"...","choices":[{"delta":{},"index":0,"finish_reason":"stop"}],"usage":{...}}

data: [DONE]
```

### For SDK Users (Workaround)

Implement in your application layer:

```go
sawFinishReason := false
for {
    chunk, err := stream.Next(ctx)
    if err != nil {
        if errors.Is(err, io.EOF) {
            if !sawFinishReason {
                // Synthesize final chunk for OpenAI compatibility
                finalChunk := &types.ChatCompletionChunk{
                    Choices: []types.ChatCompletionChunkChoice{{
                        Index:        0,
                        FinishReason: "stop",
                        Delta:        types.ChatCompletionChunkDelta{},
                    }},
                }
                // Send finalChunk to your downstream client
            }
            return
        }
        return err
    }
    
    for _, choice := range chunk.Choices {
        if choice.FinishReason != "" {
            sawFinishReason = true
        }
    }
    // Process chunk normally...
}
```

## Files Created/Modified

### New Files
1. **`FINISH_REASON_ANALYSIS.md`** - Detailed technical analysis
2. **`INVESTIGATION_SUMMARY.md`** - This file
3. **`groq/chat/finish_reason_test.go`** - Test proving SDK works correctly
4. **`groq/examples/debug_finish_reason/main.go`** - Diagnostic tool
5. **`groq/examples/debug_finish_reason/README.md`** - Tool documentation

### Modified Files
1. **`docs/GROQ_SDK_BUG_REPORT.md`** - Added investigation results

## Testing

### Run the Unit Test
```bash
go test -v ./groq/chat/ -run TestStreamFinishReason
```

**Expected**: ✓ PASS

### Run the Diagnostic Tool
```bash
export GROQ_API_KEY=your_key_here
go run ./groq/examples/debug_finish_reason/main.go
```

**Expected**: Will show whether your specific model sends `finish_reason`

## Conclusion

✅ **SDK is working correctly**  
⚠️ **API needs to be updated** to send `finish_reason` before `[DONE]`  
🔧 **Workaround available** for users requiring OpenAI compatibility

The SDK correctly implements the streaming protocol and properly processes all data the API sends. The issue is that the API is not sending the required `finish_reason` field for certain models, which breaks compatibility with OpenAI-compatible clients.

## Next Steps

1. **For Zaguán**: Continue using the workaround in CoreX (already implemented)
2. **For Groq**: Update API to send `finish_reason` before `[DONE]` for all models
3. **For SDK**: No changes needed - SDK is correct as-is

## Contact

If you need more information or want to discuss this further:
- Bug Report: `docs/GROQ_SDK_BUG_REPORT.md`
- Analysis: `FINISH_REASON_ANALYSIS.md`
- Test: `groq/chat/finish_reason_test.go`
