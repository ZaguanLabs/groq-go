# Groq Go SDK v0.3.0 Release Notes

**Release Date:** December 3, 2025  
**Status:** Stable  
**Python SDK Parity:** v0.37.0

## Overview

This release brings full API parity with the official Groq Python SDK v0.37.0, adding new fields and types for enhanced usage statistics, MCP tool discovery, JSON Schema structured outputs, and audio enhancements.

## New Features

### 1. Enhanced Usage Statistics

**CompletionUsage** now includes additional fields for detailed token tracking:

```go
type CompletionUsage struct {
    PromptTokens            int                      `json:"prompt_tokens"`
    CompletionTokens        int                      `json:"completion_tokens"`
    TotalTokens             int                      `json:"total_tokens"`
    PromptTime              float64                  `json:"prompt_time,omitempty"`
    CompletionTime          float64                  `json:"completion_time,omitempty"`
    TotalTime               float64                  `json:"total_time,omitempty"`
    QueueTime               float64                  `json:"queue_time,omitempty"`                // NEW
    CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"` // NEW
    PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`     // NEW
}

type CompletionTokensDetails struct {
    ReasoningTokens int `json:"reasoning_tokens"` // Tokens used for reasoning
}

type PromptTokensDetails struct {
    CachedTokens int `json:"cached_tokens"` // Tokens served from cache
}
```

### 2. MCP Tool Discovery

New types for Model Context Protocol (MCP) tool discovery:

```go
type McpListTool struct {
    ID          string           `json:"id,omitempty"`
    ServerLabel string           `json:"server_label,omitempty"`
    Tools       []McpListToolDef `json:"tools,omitempty"`
    Type        string           `json:"type,omitempty"`
}

type McpListToolDef struct {
    Annotations interface{}            `json:"annotations,omitempty"`
    Description string                 `json:"description,omitempty"`
    InputSchema map[string]interface{} `json:"input_schema,omitempty"`
    Name        string                 `json:"name,omitempty"`
}
```

### 3. JSON Schema Structured Output

Enhanced **ResponseFormat** now supports JSON Schema for structured outputs:

```go
type ResponseFormat struct {
    Type       string                    `json:"type"`                  // "text", "json_object", or "json_schema"
    JSONSchema *ResponseFormatJSONSchema `json:"json_schema,omitempty"` // Required when type is "json_schema"
}

type ResponseFormatJSONSchema struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description,omitempty"`
    Schema      map[string]interface{} `json:"schema,omitempty"`
    Strict      *bool                  `json:"strict,omitempty"`
}
```

**Example Usage:**
```go
resp, err := client.Chat.Create(ctx, &types.CreateChatCompletionRequest{
    Model: "llama-3.3-70b-versatile",
    Messages: []types.ChatCompletionMessageParam{
        {Role: types.RoleUser, Content: "Extract the name and age from: John is 30 years old"},
    },
    ResponseFormat: &types.ResponseFormat{
        Type: "json_schema",
        JSONSchema: &types.ResponseFormatJSONSchema{
            Name: "person_info",
            Schema: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "name": map[string]interface{}{"type": "string"},
                    "age":  map[string]interface{}{"type": "integer"},
                },
                "required": []string{"name", "age"},
            },
            Strict: ptr(true),
        },
    },
})
```

### 4. Hardware Cache Statistics

New **XGroqCacheStats** type for hardware cache metrics:

```go
type XGroqCacheStats struct {
    DramCachedTokens int `json:"dram_cached_tokens,omitempty"` // Tokens served from DRAM cache
    SramCachedTokens int `json:"sram_cached_tokens,omitempty"` // Tokens served from SRAM cache
}
```

### 5. Audio Enhancements

**CreateSpeechRequest** now supports sample rate selection:

```go
type CreateSpeechRequest struct {
    Model          string                    `json:"model"`
    Input          string                    `json:"input"`
    Voice          string                    `json:"voice"`
    ResponseFormat *option.Optional[string]  `json:"response_format,omitempty"` // flac, mp3, mulaw, ogg, wav
    SampleRate     *option.Optional[int]     `json:"sample_rate,omitempty"`     // NEW: 8000-48000 Hz
    Speed          *option.Optional[float64] `json:"speed,omitempty"`
}
```

**CreateTranscriptionRequest** now supports URL-based transcription:

```go
type CreateTranscriptionRequest struct {
    File                   interface{}               `json:"file"`
    Model                  string                    `json:"model"`
    Language               *option.Optional[string]  `json:"language,omitempty"`
    Prompt                 *option.Optional[string]  `json:"prompt,omitempty"`
    ResponseFormat         *option.Optional[string]  `json:"response_format,omitempty"`
    Temperature            *option.Optional[float64] `json:"temperature,omitempty"`
    TimestampGranularities []string                  `json:"timestamp_granularities[],omitempty"`
    URL                    *option.Optional[string]  `json:"url,omitempty"` // NEW: Required for Batch API
}
```

### 6. ChatCompletion Enhancements

**ChatCompletion** now includes additional metadata fields:

```go
type ChatCompletion struct {
    ID                string                 `json:"id"`
    Choices           []ChatCompletionChoice `json:"choices"`
    Created           int64                  `json:"created"`
    Model             string                 `json:"model"`
    SystemFingerprint string                 `json:"system_fingerprint,omitempty"`
    Object            string                 `json:"object"`
    Usage             *CompletionUsage       `json:"usage,omitempty"`
    McpListTools      []McpListTool          `json:"mcp_list_tools,omitempty"`  // NEW
    ServiceTier       string                 `json:"service_tier,omitempty"`    // NEW
    UsageBreakdown    *UsageBreakdown        `json:"usage_breakdown,omitempty"` // NEW
    XGroq             *XGroq                 `json:"x_groq,omitempty"`          // NEW
}
```

## Quality Metrics

| Metric | Value |
|--------|-------|
| Test Coverage | 73.5%+ |
| Total Tests | 238 |
| Pass Rate | 100% |
| Race Conditions | 0 |
| Backward Compatibility | ✅ Full |

## Breaking Changes

None. All changes are backward compatible.

## Migration Guide

No migration required. Simply update your dependency:

```bash
go get github.com/ZaguanLabs/groq-go@v0.3.0
```

## Files Changed

- `groq/constants.go` - Version updated to 0.3.0
- `groq/types/shared.go` - Added CompletionTokensDetails, PromptTokensDetails
- `groq/types/chat.go` - Added McpListTool, ResponseFormatJSONSchema, XGroqCacheStats
- `groq/types/audio.go` - Added SampleRate, URL fields
- `VERSION` - Updated to 0.3.0
- `CHANGELOG.md` - Added v0.3.0 release notes
- `README.md` - Updated for v0.3.0

## Acknowledgments

This release synchronizes with the official Groq Python SDK v0.37.0, ensuring Go developers have access to all the latest API features.
