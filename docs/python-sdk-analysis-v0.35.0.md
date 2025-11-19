# Groq Python SDK v0.35.0 Deep Dive Analysis

## Executive Summary

This document provides a comprehensive analysis of the Groq Python SDK v0.35.0 (released 2025-11-18) to inform improvements to the Go SDK. The analysis covers new features, API changes, type definitions, and best practices that should be incorporated into our Go implementation.

---

## 1. Recent Changes (v0.35.0)

### 1.1 Key Updates

**Features:**
- **API Updates** (commits d6ec93d, 48c8f11): Two API updates were merged, likely adding new fields or endpoints.

**Bug Fixes:**
- **Chat Completion Streaming Types** (commit 833cf83): Fixed streaming response types, particularly around `ChatCompletionChunk`.

**Previous Version Context (v0.34.1):**
- Added annotations to streaming response
- Fixed stream closing without requiring full consumption
- Python 3.14 compatibility

---

## 2. Critical New Features Missing from Our Go SDK

### 2.1 Compound AI Systems

**What it is:** Groq now supports "Compound" models that combine multiple models and tools for complex tasks.

**Python API:**
```python
chat_completion = client.chat.completions.create(
    model="compound-beta",  # or "compound-beta-mini"
    messages=[...],
    compound_custom={
        "models": {
            "answering_model": "llama-3.3-70b-versatile",
            "reasoning_model": "qwen/qwen3-32b"
        },
        "tools": {
            "enabled_tools": ["web_search", "code_interpreter"],
            "wolfram_settings": {
                "authorization": "API_KEY"
            }
        }
    }
)
```

**Go Implementation Needed:**
```go
type CompoundCustom struct {
    Models *CompoundCustomModels `json:"models,omitempty"`
    Tools  *CompoundCustomTools  `json:"tools,omitempty"`
}

type CompoundCustomModels struct {
    AnsweringModel *option.Optional[string] `json:"answering_model,omitempty"`
    ReasoningModel *option.Optional[string] `json:"reasoning_model,omitempty"`
}

type CompoundCustomTools struct {
    EnabledTools     []string                          `json:"enabled_tools,omitempty"`
    WolframSettings  *CompoundCustomToolsWolframSettings `json:"wolfram_settings,omitempty"`
}
```

**Impact:** HIGH - This is a major new feature for advanced AI workflows.

---

### 2.2 Enhanced Streaming Response (`x_groq` Field)

**What it is:** Streaming responses now include a `x_groq` field with additional metadata.

**Python Type:**
```python
class XGroq(BaseModel):
    id: Optional[str] = None  # Groq request ID (first & final chunk)
    debug: Optional[XGroqDebug] = None  # Token IDs/strings (if debug=true)
    seed: Optional[int] = None  # Seed used (final chunk)
    usage: Optional[CompletionUsage] = None  # Usage stats (final chunk)
    usage_breakdown: Optional[UsageBreakdown] = None  # Per-model usage
    error: Optional[str] = None  # Error string if stream stopped early
```

**Go Implementation Needed:**
```go
type XGroq struct {
    ID             *string                  `json:"id,omitempty"`
    Debug          *XGroqDebug              `json:"debug,omitempty"`
    Seed           *int                     `json:"seed,omitempty"`
    Usage          *CompletionUsage         `json:"usage,omitempty"`
    UsageBreakdown *UsageBreakdown          `json:"usage_breakdown,omitempty"`
    Error          *string                  `json:"error,omitempty"`
}

type XGroqDebug struct {
    InputTokenIDs  []int    `json:"input_token_ids,omitempty"`
    InputTokens    []string `json:"input_tokens,omitempty"`
    OutputTokenIDs []int    `json:"output_token_ids,omitempty"`
    OutputTokens   []string `json:"output_tokens,omitempty"`
}

type UsageBreakdown struct {
    Models []UsageBreakdownModel `json:"models"`
}

type UsageBreakdownModel struct {
    Model string           `json:"model"`
    Usage CompletionUsage  `json:"usage"`
}
```

**Update `ChatCompletionChunk`:**
```go
type ChatCompletionChunk struct {
    // ... existing fields ...
    XGroq *XGroq `json:"x_groq,omitempty"`
}
```

**Impact:** MEDIUM - Important for debugging and multi-model tracking.

---

### 2.3 Annotations and Citations

**What it is:** Responses can include annotations that cite documents or function calls.

**Python Types:**
```python
class Annotation(BaseModel):
    type: Literal["document_citation", "function_citation"]
    document_citation: Optional[AnnotationDocumentCitation] = None
    function_citation: Optional[AnnotationFunctionCitation] = None

class AnnotationDocumentCitation(BaseModel):
    document_id: str
    start_index: int
    end_index: int

class AnnotationFunctionCitation(BaseModel):
    tool_call_id: str
    start_index: int
    end_index: int
```

**Go Implementation Needed:**
```go
type Annotation struct {
    Type              string                     `json:"type"` // "document_citation" or "function_citation"
    DocumentCitation  *AnnotationDocumentCitation `json:"document_citation,omitempty"`
    FunctionCitation  *AnnotationFunctionCitation `json:"function_citation,omitempty"`
}

type AnnotationDocumentCitation struct {
    DocumentID string `json:"document_id"`
    StartIndex int    `json:"start_index"`
    EndIndex   int    `json:"end_index"`
}

type AnnotationFunctionCitation struct {
    ToolCallID string `json:"tool_call_id"`
    StartIndex int    `json:"start_index"`
    EndIndex   int    `json:"end_index"`
}
```

**Update `ChatCompletionMessage` and `ChoiceDelta`:**
```go
type ChatCompletionMessage struct {
    // ... existing fields ...
    Annotations []Annotation `json:"annotations,omitempty"`
}

type ChoiceDelta struct {
    // ... existing fields ...
    Annotations []Annotation `json:"annotations,omitempty"`
}
```

**Impact:** MEDIUM - Enables citation-based responses.

---

### 2.4 Executed Tools (Compound AI)

**What it is:** For compound AI systems, responses include details about tools that were executed.

**Python Types:**
```python
class ExecutedTool(BaseModel):
    arguments: str
    index: int
    type: str
    browser_results: Optional[List[ExecutedToolBrowserResult]] = None
    code_results: Optional[List[ExecutedToolCodeResult]] = None
    output: Optional[str] = None
    search_results: Optional[ExecutedToolSearchResults] = None
```

**Go Implementation Needed:**
```go
type ExecutedTool struct {
    Arguments      string                      `json:"arguments"`
    Index          int                         `json:"index"`
    Type           string                      `json:"type"`
    BrowserResults []ExecutedToolBrowserResult `json:"browser_results,omitempty"`
    CodeResults    []ExecutedToolCodeResult    `json:"code_results,omitempty"`
    Output         *string                     `json:"output,omitempty"`
    SearchResults  *ExecutedToolSearchResults  `json:"search_results,omitempty"`
}

type ExecutedToolBrowserResult struct {
    Title       string  `json:"title"`
    URL         string  `json:"url"`
    Content     *string `json:"content,omitempty"`
    LiveViewURL *string `json:"live_view_url,omitempty"`
}

type ExecutedToolCodeResult struct {
    Chart  *ExecutedToolCodeResultChart   `json:"chart,omitempty"`
    Charts []ExecutedToolCodeResultChart  `json:"charts,omitempty"`
    PNG    *string                        `json:"png,omitempty"` // Base64 encoded
    Text   *string                        `json:"text,omitempty"`
}

type ExecutedToolSearchResults struct {
    Images  []string                         `json:"images,omitempty"`
    Results []ExecutedToolSearchResultsResult `json:"results,omitempty"`
}
```

**Impact:** HIGH - Critical for compound AI system support.

---

### 2.5 Documents in Chat Requests

**What it is:** You can now provide documents as context for the conversation.

**Python API:**
```python
chat_completion = client.chat.completions.create(
    model="llama-3.3-70b-versatile",
    messages=[...],
    documents=[
        {
            "id": "doc1",
            "source": {
                "type": "text",
                "text": "Document content here..."
            }
        },
        {
            "id": "doc2",
            "source": {
                "type": "json",
                "data": {"key": "value"}
            }
        }
    ],
    citation_options="enabled"
)
```

**Go Implementation Needed:**
```go
type Document struct {
    ID     *string         `json:"id,omitempty"`
    Source DocumentSource  `json:"source"`
}

type DocumentSource struct {
    Type string                 `json:"type"` // "text" or "json"
    Text *string                `json:"text,omitempty"`
    Data map[string]interface{} `json:"data,omitempty"`
}
```

**Update `CreateChatCompletionRequest`:**
```go
type CreateChatCompletionRequest struct {
    // ... existing fields ...
    Documents        []Document                   `json:"documents,omitempty"`
    CitationOptions  *option.Optional[string]     `json:"citation_options,omitempty"` // "enabled" or "disabled"
}
```

**Impact:** HIGH - Enables RAG-like workflows.

---

### 2.6 Reasoning Features

**What it is:** Support for reasoning models with configurable reasoning output.

**Python API:**
```python
chat_completion = client.chat.completions.create(
    model="qwen/qwen3-32b",
    messages=[...],
    reasoning_effort="medium",  # "none", "default", "low", "medium", "high"
    reasoning_format="parsed",  # "hidden", "raw", "parsed"
    include_reasoning=True
)
```

**Go Implementation Needed:**
```go
type CreateChatCompletionRequest struct {
    // ... existing fields ...
    ReasoningEffort  *option.Optional[string] `json:"reasoning_effort,omitempty"`
    ReasoningFormat  *option.Optional[string] `json:"reasoning_format,omitempty"`
    IncludeReasoning *option.Optional[bool]   `json:"include_reasoning,omitempty"`
}

type ChatCompletionMessage struct {
    // ... existing fields ...
    Reasoning *string `json:"reasoning,omitempty"`
}

type ChoiceDelta struct {
    // ... existing fields ...
    Reasoning *string `json:"reasoning,omitempty"`
}
```

**Impact:** MEDIUM - Important for reasoning models.

---

### 2.7 Search Settings

**What it is:** Fine-grained control over web search when using search tools.

**Python API:**
```python
chat_completion = client.chat.completions.create(
    model="compound-beta",
    messages=[...],
    search_settings={
        "country": "united states",
        "include_domains": ["wikipedia.org", "arxiv.org"],
        "exclude_domains": ["example.com"],
        "include_images": True
    }
)
```

**Go Implementation Needed:**
```go
type SearchSettings struct {
    Country        *option.Optional[string] `json:"country,omitempty"`
    IncludeDomains []string                 `json:"include_domains,omitempty"`
    ExcludeDomains []string                 `json:"exclude_domains,omitempty"`
    IncludeImages  *option.Optional[bool]   `json:"include_images,omitempty"`
}
```

**Update `CreateChatCompletionRequest`:**
```go
type CreateChatCompletionRequest struct {
    // ... existing fields ...
    SearchSettings *SearchSettings `json:"search_settings,omitempty"`
}
```

**Impact:** MEDIUM - Useful for search-enabled models.

---

### 2.8 New Request Parameters

**Additional fields in `CreateChatCompletionRequest`:**

```go
type CreateChatCompletionRequest struct {
    // ... existing fields ...
    
    // Service tier selection
    ServiceTier *option.Optional[string] `json:"service_tier,omitempty"` // "auto", "on_demand", "flex", "performance"
    
    // Token limits
    MaxCompletionTokens *option.Optional[int] `json:"max_completion_tokens,omitempty"`
    
    // Tool validation
    DisableToolValidation bool `json:"disable_tool_validation,omitempty"`
    
    // Metadata (not currently supported but in API)
    Metadata map[string]string `json:"metadata,omitempty"`
    
    // Store flag (not currently supported but in API)
    Store *option.Optional[bool] `json:"store,omitempty"`
}
```

**Impact:** LOW-MEDIUM - Completeness and future-proofing.

---

### 2.9 New Model Identifiers

**Models added in recent versions:**
- `compound-beta`
- `compound-beta-mini`
- `meta-llama/llama-4-maverick-17b-128e-instruct`
- `meta-llama/llama-4-scout-17b-16e-instruct`
- `meta-llama/llama-guard-4-12b`
- `moonshotai/kimi-k2-instruct`
- `openai/gpt-oss-120b`
- `openai/gpt-oss-20b`
- `qwen/qwen3-32b`

**Go Implementation:**
Update `groq/types/enums.go` or create model constants:
```go
const (
    ModelCompoundBeta                = "compound-beta"
    ModelCompoundBetaMini            = "compound-beta-mini"
    ModelLlama4Maverick17B           = "meta-llama/llama-4-maverick-17b-128e-instruct"
    ModelLlama4Scout17B              = "meta-llama/llama-4-scout-17b-16e-instruct"
    ModelLlamaGuard412B              = "meta-llama/llama-guard-4-12b"
    ModelKimiK2Instruct              = "moonshotai/kimi-k2-instruct"
    ModelGPTOSS120B                  = "openai/gpt-oss-120b"
    ModelGPTOSS20B                   = "openai/gpt-oss-20b"
    ModelQwen332B                    = "qwen/qwen3-32b"
    // ... existing models ...
)
```

**Impact:** LOW - Easy to add, important for completeness.

---

## 3. Architecture Insights from Python SDK

### 3.1 Streaming Implementation

**Python approach:**
- Returns an iterator that yields `ChatCompletionChunk` objects
- Final chunk contains usage information in `x_groq.usage`
- Stream can be closed early without consuming all chunks
- Supports `with` context manager for automatic cleanup

**Our Go implementation is solid but could add:**
```go
// Add helper to check if chunk is final
func (c *ChatCompletionChunk) IsFinal() bool {
    return c.XGroq != nil && c.XGroq.Usage != nil
}

// Add helper to get usage from stream
func (s *Stream[T]) GetUsage() (*CompletionUsage, error) {
    // Consume stream until final chunk and return usage
}
```

### 3.2 Error Handling

**Python SDK has specific error types:**
- `APIConnectionError` - Network issues
- `APIStatusError` - HTTP status errors (4xx, 5xx)
- `APITimeoutError` - Timeout errors
- `RateLimitError` - 429 errors
- `BadRequestError` - 400 errors
- `AuthenticationError` - 401 errors
- etc.

**Our Go SDK already has this covered** in `groq/errors.go`.

### 3.3 Type Safety

**Python uses:**
- `TypedDict` for request parameters (similar to our structs)
- `Pydantic` models for responses (similar to our structs with JSON tags)
- `Literal` types for enums (we use string constants)
- `Union` types for polymorphic fields (we use interfaces or `interface{}`)

**Recommendation:** Our approach is idiomatic for Go. No changes needed.

---

## 4. Testing Insights

### 4.1 Test Coverage Areas

From Python SDK tests directory:
- Streaming response handling
- Error scenarios (network, timeout, API errors)
- Pagination (if applicable)
- File uploads
- Async/sync parity
- Type validation

**Recommendation:** Add integration tests for:
- Compound AI workflows
- Document-based chat
- Reasoning models
- Citation parsing

---

## 5. Priority Implementation Roadmap

### Phase 1: Critical (Immediate)
1. **Compound AI Support**
   - Add `CompoundCustom` types
   - Add `compound_custom` field to request
   - Add `ExecutedTool` types
   - Update response types

2. **XGroq Metadata**
   - Add `XGroq` type
   - Update `ChatCompletionChunk`
   - Add helper methods

3. **Documents & Citations**
   - Add `Document` types
   - Add `Annotation` types
   - Update request/response types

### Phase 2: Important (Short-term)
4. **Reasoning Support**
   - Add reasoning fields to request
   - Add reasoning fields to response

5. **Search Settings**
   - Add `SearchSettings` type
   - Update request type

6. **New Request Parameters**
   - Add `service_tier`, `max_completion_tokens`, etc.

### Phase 3: Nice-to-have (Medium-term)
7. **Model Constants**
   - Add new model identifiers

8. **Enhanced Examples**
   - Compound AI example
   - Document-based RAG example
   - Reasoning model example

---

## 6. Breaking Changes to Consider

**None identified.** All new features are additive (optional fields).

---

## 7. Documentation Updates Needed

1. Update README with compound AI example
2. Add migration guide section for new features
3. Update CHANGELOG with new types
4. Add GoDoc comments for all new types

---

## 8. Compatibility Notes

- Python SDK supports Python 3.9+ (dropped 3.8 in v0.34.0)
- Our Go SDK requires Go 1.18+ (for generics)
- No breaking changes in Python SDK v0.35.0
- All new features are backward compatible

---

## 9. Recommendations Summary

### Must-Have (Before v1.0)
✅ Compound AI system support
✅ XGroq metadata in streaming
✅ Documents and citations
✅ Reasoning model support

### Should-Have (v1.0 or v1.1)
✅ Search settings
✅ New request parameters
✅ Model constant updates

### Nice-to-Have (Future)
✅ Enhanced helper methods
✅ More comprehensive examples
✅ Integration tests for new features

---

## 10. Code Quality Observations

**Python SDK strengths to emulate:**
- Comprehensive type hints
- Clear separation of concerns (resources, types, utils)
- Extensive examples
- Good error messages
- Automatic retry logic

**Our Go SDK already does well:**
- Strong typing with generics
- Idiomatic Go patterns
- Context support throughout
- Comprehensive error types
- Retry logic with backoff

---

## Conclusion

The Groq Python SDK v0.35.0 introduces significant new features, particularly around **Compound AI systems**, **document-based context**, and **enhanced streaming metadata**. Our Go SDK is architecturally sound but needs to add these new types and fields to maintain feature parity.

The implementation is straightforward as all changes are additive. No breaking changes are required. Priority should be given to Compound AI support and XGroq metadata as these represent the most significant new capabilities.
