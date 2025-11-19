# Finish Reason Diagnostic Tool

This tool helps diagnose whether the Groq API is sending `finish_reason` in streaming responses.

## Purpose

OpenAI-compatible clients expect to receive a `finish_reason` field in the final chunk before the stream ends. This tool verifies whether the Groq API is sending this field correctly.

## Usage

```bash
export GROQ_API_KEY=your_api_key_here
go run ./groq/examples/debug_finish_reason/main.go
```

## Expected Output

### If API is Working Correctly ✅

```
Testing finish_reason in streaming responses...
================================================

Streaming chunks:
Chunk 1:
  Model: "llama-3.3-70b-versatile"
  Choices: 1
  Choice 0:
    Content: "Hello"
    FinishReason: ""

Chunk 2:
  Model: "llama-3.3-70b-versatile"
  Choices: 1
  Choice 0:
    Content: " World"
    FinishReason: ""

Chunk 3:
  Model: "llama-3.3-70b-versatile"
  Choices: 1
  Choice 0:
    Content: ""
    FinishReason: "stop"
    ✓ FINISH REASON DETECTED: "stop"

✓ Received EOF after 3 chunks

================================================
Test Results:
Total chunks received: 3
✓ SUCCESS: finish_reason was received before EOF
```

### If API Has Issues ⚠️

```
Testing finish_reason in streaming responses...
================================================

Streaming chunks:
Chunk 1:
  Model: "moonshotai/kimi-k2-instruct-0905"
  Choices: 1
  Choice 0:
    Content: "Hello"
    FinishReason: ""

Chunk 2:
  Model: ""
  Choices: 0

✓ Received EOF after 2 chunks

================================================
Test Results:
Total chunks received: 2
✗ FAIL: finish_reason was NOT received before EOF

Last chunk details:
  Model: ""
  Choices: 0

This indicates the Groq API is not sending finish_reason in streaming mode.
This is a server-side issue, not an SDK bug.
```

## What This Tells You

- **SUCCESS**: The API is sending `finish_reason` correctly, and the SDK is processing it
- **FAIL**: The API is not sending `finish_reason`, which breaks OpenAI compatibility

## Related Files

- `FINISH_REASON_ANALYSIS.md` - Complete analysis of the issue
- `groq/chat/finish_reason_test.go` - Unit test proving SDK handles finish_reason correctly
- `docs/GROQ_SDK_BUG_REPORT.md` - Original bug report and investigation results
