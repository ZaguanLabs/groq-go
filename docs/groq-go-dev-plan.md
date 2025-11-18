# Go Groq SDK Port – Comprehensive Development Plan

**Version:** 0.34.1 Port  
**Target Go Version:** 1.21+  
**Python SDK Reference:** groq-python 0.34.1 (Stainless-generated)

---

## 1. Objectives & Scope

### 1.1 Primary Goals
- **Feature Parity**: Deliver a first-class Go SDK mirroring the Python package's ergonomics (typed requests/responses, retries, streaming, raw responses, pagination) so Go teams can adopt Groq APIs without hand-writing HTTP calls.@docs/groq-0.34.1/README.md#6-485
- **Dual-Mode Support**: Maintain feature parity across synchronous and asynchronous flows (Go blocking client vs. goroutine-friendly streaming helpers) with consistent configuration knobs for API key, base URL, timeout, retries, and custom headers/queries.@docs/groq-0.34.1/src/groq/_client.py#45-628
- **Production-Grade Tooling**: Provide batteries-included tooling: lint/format/test automation, semantic versioning, and reproducible releases similar to the Python project's hatch/rye/ruff pipeline.@docs/groq-0.34.1/pyproject.toml#61-267

### 1.2 Non-Goals
- **Code Generation Automation**: Initial release will manually port types; future automation via OpenAPI spec is deferred
- **Pydantic Equivalence**: Go will use struct tags + custom JSON marshalers instead of runtime validation framework
- **Python 3.9+ Compatibility**: Go-specific idioms take precedence over Python syntax mimicry

## 2. Python SDK Deep Architecture Analysis

### 2.1 Module Structure & Exports
**Entry Point** (`__init__.py`):@docs/groq-0.34.1/src/groq/__init__.py#1-93
- Exports 74 symbols including `Groq`, `AsyncGroq`, `BaseModel`, error classes, constants, and type utilities
- Uses `__all__` for explicit public API surface
- Modifies `__module__` attribute of exported symbols to point to `groq` for cleaner error messages
- Lazy-loads resources via `_resources_proxy` when not type-checking
- Initializes logging via `_setup_logging()` on import

**Client Classes**:@docs/groq-0.34.1/src/groq/_client.py#45-628
- `Groq(SyncAPIClient)`: Synchronous client with `httpx.Client` backend
- `AsyncGroq(AsyncAPIClient)`: Async client with `httpx.AsyncClient` backend
- Both inherit from `BaseClient[HttpxClientT, DefaultStreamT]` generic
- Expose 6 resource properties via `@cached_property`: `chat`, `embeddings`, `audio`, `models`, `batches`, `files`
- Provide `with_raw_response` and `with_streaming_response` wrappers for each resource
- Implement `copy()`/`with_options()` for immutable client cloning with overrides

### 2.2 Base Client Infrastructure
**BaseClient** (`_base_client.py`):@docs/groq-0.34.1/src/groq/_base_client.py#359-799
- **Configuration**:
  - `_version`: SDK version string
  - `_base_url`: URL with enforced trailing slash
  - `max_retries`: Default 2 retries@docs/groq-0.34.1/src/groq/_constants.py#10
  - `timeout`: Default 60s total, 5s connect@docs/groq-0.34.1/src/groq/_constants.py#9
  - `_strict_response_validation`: Enable Pydantic validation errors
  - `_idempotency_header`: Optional header name for idempotency keys
  - `_custom_headers`, `_custom_query`: User-provided defaults

- **HTTP Client Defaults**:@docs/groq-0.34.1/src/groq/_base_client.py#787-793
  - Connection limits: 100 max connections, 20 keepalive@docs/groq-0.34.1/src/groq/_constants.py#11
  - Follow redirects: `True`
  - Timeout: `httpx.Timeout(timeout=60, connect=5.0)`

- **Request Building**:@docs/groq-0.34.1/src/groq/_base_client.py#473-556
  - Merges base URL with relative paths (strips leading `/` from path)
  - Builds headers: default + auth + custom + retry count + read timeout
  - Handles multipart/form-data by removing Content-Type to let httpx set boundary
  - Serializes query params via `Querystring.stringify()` with comma array format
  - Supports `extra_json` merging into request body
  - Adds `x-stainless-retry-count` and `x-stainless-read-timeout` headers

- **Retry Logic**:@docs/groq-0.34.1/src/groq/_base_client.py#724-781
  - Retries on: 408 (timeout), 409 (conflict), 429 (rate limit), >=500 (server errors)
  - Respects `x-should-retry` header if present
  - Exponential backoff: `min(0.5 * 2^retries, 8.0)` seconds@docs/groq-0.34.1/src/groq/_constants.py#13-14
  - Jitter: ±25% randomization
  - Honors `Retry-After` header (both seconds and HTTP-date formats)
  - Also supports non-standard `retry-after-ms` header for millisecond precision

### 2.3 Authentication & Headers
**Auth Headers**:@docs/groq-0.34.1/src/groq/_client.py#149-153
- `Authorization: Bearer {api_key}`
- API key loaded from `GROQ_API_KEY` env var if not provided
- Raises `GroqError` if API key missing

**Default Headers**:@docs/groq-0.34.1/src/groq/_base_client.py#644-653
- `Accept: application/json`
- `Content-Type: application/json`
- `User-Agent: Groq/Python {version}`
- `X-Stainless-Async: false` (sync) or `async:{library}` (async)
- Platform headers: OS, Python version, architecture

**Base URL**:@docs/groq-0.34.1/src/groq/_client.py#84-87
- Default: `https://api.groq.com`
- Overridable via `GROQ_BASE_URL` env var or constructor arg

### 2.4 Resources & API Methods
**Resource Pattern**:@docs/groq-0.34.1/src/groq/resources/chat/completions.py#32-512
- Each resource inherits `SyncAPIResource` or `AsyncAPIResource`
- Exposes `with_raw_response` and `with_streaming_response` wrappers
- Methods use `@overload` for stream vs. non-stream return types
- All methods accept `extra_headers`, `extra_query`, `extra_body`, `timeout` overrides

**Chat Completions**:@docs/groq-0.34.1/src/groq/resources/chat/completions.py#244-512
- Endpoint: `POST /openai/v1/chat/completions`
- 30+ parameters including: messages, model, temperature, max_tokens, stream, tools, etc.
- Streaming returns `Stream[ChatCompletionChunk]`, non-streaming returns `ChatCompletion`
- Supports function calling, tool use, structured outputs (JSON schema)
- Special parameters: `citation_options`, `compound_custom`, `reasoning_effort`, `search_settings`

**Embeddings**:@docs/groq-0.34.1/src/groq/resources/embeddings.py#47-100
- Endpoint: `POST /openai/v1/embeddings`
- Parameters: input (str or list), model, encoding_format (float/base64), user
- Returns `CreateEmbeddingResponse` with embedding vectors

**Models**:@docs/groq-0.34.1/src/groq/resources/models.py#24-127
- `GET /openai/v1/models` - list all models
- `GET /openai/v1/models/{model}` - retrieve specific model
- `DELETE /openai/v1/models/{model}` - delete model

### 2.5 Type System & Validation
**BaseModel**:@docs/groq-0.34.1/src/groq/_models.py#84-384
- Extends `pydantic.BaseModel` with custom `to_dict()` and `to_json()` methods
- `to_dict()` options: mode (json/python), use_api_names, exclude_unset/defaults/none
- `to_json()` returns indented JSON string (default indent=2)
- Custom `construct()` method for validation-free instantiation
- Supports both Pydantic v1 and v2 via compatibility layer
- `model_fields_set` property tracks which fields were explicitly set

**Type Construction**:@docs/groq-0.34.1/src/groq/_models.py#469-586
- `construct_type()`: Loose coercion with nested value construction
- Handles unions via discriminated union metadata or fallback to first valid variant
- Supports `dict`, `list`, `BaseModel`, `datetime`, `date`, `float` coercion
- Discriminated unions use `PropertyInfo.discriminator` annotation
- Caches discriminator details in `WeakKeyDictionary`

**Validation**:@docs/groq-0.34.1/src/groq/_models.py#719-724
- `validate_type()`: Strict validation using Pydantic's `TypeAdapter` (v2) or `RootModel` (v1)
- Raises `pydantic.ValidationError` on mismatch
- Used when `_strict_response_validation=True`

### 2.6 File Handling
**File Upload Types**:@docs/groq-0.34.1/src/groq/_files.py#23-124
- Accepts: `bytes`, `io.IOBase`, `os.PathLike`, or tuples `(filename, content, content_type, headers)`
- Sync: Reads `PathLike` via `pathlib.Path.read_bytes()`
- Async: Reads `PathLike` via `anyio.Path.read_bytes()` for non-blocking I/O
- Transforms to httpx-compatible file types before request

**Multipart Encoding**:@docs/groq-0.34.1/src/groq/_base_client.py#558-586
- Serializes nested dicts/arrays using `Querystring.stringify_items()` with `array_format="brackets"`
- Converts duplicate keys to lists for httpx multipart handling
- Removes `Content-Type: multipart/form-data` header to let httpx set boundary

### 2.7 Streaming Implementation
**SSE Decoder**:@docs/groq-0.34.1/src/groq/_streaming.py#264-367
- Implements Server-Sent Events (SSE) protocol per WHATWG spec
- Parses fields: `event`, `data`, `id`, `retry`
- Handles chunk boundaries: `\r\r`, `\n\n`, `\r\n\r\n`
- Accumulates multi-line `data` fields with `\n` separator
- Ignores comment lines (starting with `:`)
- Filters null bytes in `id` field

**Stream Classes**:@docs/groq-0.34.1/src/groq/_streaming.py#22-221
- `Stream[T]`: Sync iterator over SSE events
- `AsyncStream[T]`: Async iterator over SSE events
- Both implement context managers for automatic cleanup
- Parse JSON from `sse.data` and yield typed chunks
- Detect `[DONE]` sentinel to terminate stream
- Handle error events: extract message from `{"error": {"message": "..."}}`
- Close HTTP response after iteration completes

**Usage Pattern**:@docs/groq-0.34.1/examples/chat_completion_streaming.py#5-53
```python
stream = client.chat.completions.create(stream=True, ...)
for chunk in stream:
    print(chunk.choices[0].delta.content, end="")
    if chunk.choices[0].finish_reason:
        print(chunk.x_groq.usage)  # Usage on final chunk
```

### 2.8 Response Parsing
**APIResponse**:@docs/groq-0.34.1/src/groq/_response.py#48-267
- Wraps `httpx.Response` with typed parsing
- Exposes: `headers`, `status_code`, `url`, `method`, `http_version`, `elapsed`, `retries_taken`
- `parse()` method with type override support
- Handles: `str`, `bytes`, `int`, `float`, `bool`, `BaseModel`, `dict`, `list`, `httpx.Response`
- Validates `Content-Type: application/json` if `_strict_response_validation=True`
- Falls back to text response if JSON parsing fails

**Raw Response Wrappers**:@docs/groq-0.34.1/src/groq/resources/embeddings.py#179-213
- `{Resource}WithRawResponse`: Returns `APIResponse[T]` instead of `T`
- `{Resource}WithStreamingResponse`: Returns streaming context manager
- Implemented via decorator pattern wrapping resource methods

### 2.9 Querystring Encoding
**Querystring Class**:@docs/groq-0.34.1/src/groq/_qs.py#23-151
- **Array Formats**:
  - `comma`: `key=val1,val2,val3`
  - `repeat`: `key=val1&key=val2&key=val3`
  - `brackets`: `key[]=val1&key[]=val2`
  - `indices`: Not implemented
- **Nested Formats**:
  - `dots`: `parent.child=value`
  - `brackets`: `parent[child]=value`
- **Primitive Serialization**:
  - `True` → `"true"`, `False` → `"false"`, `None` → `""`
  - Numbers/strings → `str(value)`
- Default: `array_format="repeat"`, `nested_format="brackets"`
- Groq client overrides: `array_format="comma"`@docs/groq-0.34.1/src/groq/_client.py#145-148

### 2.10 Error Hierarchy
**Exception Classes**:@docs/groq-0.34.1/src/groq/_exceptions.py#21-109
```
GroqError (base)
├── APIError
│   ├── APIConnectionError
│   │   └── APITimeoutError
│   ├── APIResponseValidationError
│   └── APIStatusError
│       ├── BadRequestError (400)
│       ├── AuthenticationError (401)
│       ├── PermissionDeniedError (403)
│       ├── NotFoundError (404)
│       ├── ConflictError (409)
│       ├── UnprocessableEntityError (422)
│       ├── RateLimitError (429)
│       └── InternalServerError (>=500)
```

**Error Attributes**:
- `message`: Human-readable error description
- `request`: `httpx.Request` object
- `response`: `httpx.Response` object (for status errors)
- `status_code`: HTTP status code
- `body`: Parsed JSON or raw text from response

**Error Construction**:@docs/groq-0.34.1/src/groq/_base_client.py#400-429
- Attempts JSON parse of response body
- Falls back to raw text if JSON invalid
- Formats message: `"Error code: {status_code} - {body}"`
- Maps status codes to specific error classes via `_make_status_error()`

### 2.11 Constants & Defaults
**Key Constants**:@docs/groq-0.34.1/src/groq/_constants.py#1-15
- `RAW_RESPONSE_HEADER = "X-Stainless-Raw-Response"`
- `OVERRIDE_CAST_TO_HEADER = "____stainless_override_cast_to"`
- `DEFAULT_TIMEOUT = httpx.Timeout(timeout=60, connect=5.0)`
- `DEFAULT_MAX_RETRIES = 2`
- `DEFAULT_CONNECTION_LIMITS = httpx.Limits(max_connections=100, max_keepalive_connections=20)`
- `INITIAL_RETRY_DELAY = 0.5`
- `MAX_RETRY_DELAY = 8.0`

## 3. Proposed Go SDK Architecture

### 3.1 Package Structure
```
groq-go/
├── client.go              # Client, ClientConfig, NewClient()
├── options.go             # RequestOptions, ClientOptions
├── errors.go              # Error types matching Python hierarchy
├── constants.go           # Timeouts, limits, headers
├── chat/
│   ├── completions.go     # Completions service
│   └── types.go           # Chat-specific types
├── embeddings/
│   └── embeddings.go
├── audio/
│   ├── speech.go
│   ├── transcriptions.go
│   └── translations.go
├── models/
│   └── models.go
├── batches/
│   └── batches.go
├── files/
│   └── files.go
├── types/
│   ├── chat.go            # ChatCompletion, ChatCompletionChunk, etc.
│   ├── embeddings.go
│   ├── shared.go          # ErrorObject, FunctionDefinition, etc.
│   └── enums.go           # Role, FinishReason, etc.
├── internal/
│   ├── retry/
│   │   └── retry.go       # Exponential backoff logic
│   ├── sse/
│   │   └── decoder.go     # SSE parsing
│   ├── querystring/
│   │   └── encode.go      # Query param encoding
│   └── files/
│       └── upload.go      # File upload helpers
└── examples/
    ├── chat_completion.go
    └── streaming.go
```

### 3.2 Core Client Design
**Client Struct**:
```go
type Client struct {
    httpClient *http.Client
    config     *ClientConfig
    
    // Resources
    Chat       *chat.Completions
    Embeddings *embeddings.Embeddings
    Audio      *audio.Audio
    Models     *models.Models
    Batches    *batches.Batches
    Files      *files.Files
}

type ClientConfig struct {
    APIKey         string
    BaseURL        string
    MaxRetries     int
    Timeout        time.Duration
    ConnectTimeout time.Duration
    Headers        map[string]string
    QueryParams    map[string]string
    
    // Advanced
    StrictValidation bool
    Logger           Logger
    HTTPClient       *http.Client  // Optional custom client
}
```

**Constructor**:
```go
func NewClient(opts ...ClientOption) (*Client, error) {
    cfg := &ClientConfig{
        APIKey:         os.Getenv("GROQ_API_KEY"),
        BaseURL:        getEnvOrDefault("GROQ_BASE_URL", "https://api.groq.com"),
        MaxRetries:     2,
        Timeout:        60 * time.Second,
        ConnectTimeout: 5 * time.Second,
        Headers:        make(map[string]string),
        QueryParams:    make(map[string]string),
    }
    
    for _, opt := range opts {
        opt(cfg)
    }
    
    if cfg.APIKey == "" {
        return nil, errors.New("API key required: set GROQ_API_KEY or use WithAPIKey()")
    }
    
    if cfg.HTTPClient == nil {
        cfg.HTTPClient = &http.Client{
            Timeout: cfg.Timeout,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 20,
                IdleConnTimeout:     90 * time.Second,
                DialContext: (&net.Dialer{
                    Timeout: cfg.ConnectTimeout,
                }).DialContext,
            },
        }
    }
    
    c := &Client{
        httpClient: cfg.HTTPClient,
        config:     cfg,
    }
    
    // Initialize resources
    c.Chat = chat.NewCompletions(c)
    c.Embeddings = embeddings.New(c)
    // ... etc
    
    return c, nil
}
```

**Functional Options**:
```go
type ClientOption func(*ClientConfig)

func WithAPIKey(key string) ClientOption {
    return func(c *ClientConfig) { c.APIKey = key }
}

func WithBaseURL(url string) ClientOption {
    return func(c *ClientConfig) { c.BaseURL = url }
}

func WithMaxRetries(n int) ClientOption {
    return func(c *ClientConfig) { c.MaxRetries = n }
}

func WithTimeout(d time.Duration) ClientOption {
    return func(c *ClientConfig) { c.Timeout = d }
}

func WithHTTPClient(client *http.Client) ClientOption {
    return func(c *ClientConfig) { c.HTTPClient = client }
}

func WithHeader(key, value string) ClientOption {
    return func(c *ClientConfig) { c.Headers[key] = value }
}
```

### 3.3 Request Building
**Request Options**:
```go
type RequestOptions struct {
    Headers     map[string]string
    QueryParams map[string]string
    Timeout     *time.Duration
    MaxRetries  *int
    
    // Internal
    idempotencyKey string
}

func (c *Client) buildRequest(ctx context.Context, method, path string, body interface{}, opts *RequestOptions) (*http.Request, error) {
    url := c.buildURL(path, opts)
    
    var bodyReader io.Reader
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("marshal request body: %w", err)
        }
        bodyReader = bytes.NewReader(jsonData)
    }
    
    req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
    if err != nil {
        return nil, err
    }
    
    c.setHeaders(req, opts)
    return req, nil
}

func (c *Client) setHeaders(req *http.Request, opts *RequestOptions) {
    // Default headers
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", fmt.Sprintf("Groq/Go %s", Version))
    req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
    
    // Platform headers
    req.Header.Set("X-Stainless-Lang", "go")
    req.Header.Set("X-Stainless-Package-Version", Version)
    req.Header.Set("X-Stainless-OS", runtime.GOOS)
    req.Header.Set("X-Stainless-Arch", runtime.GOARCH)
    req.Header.Set("X-Stainless-Runtime", runtime.Version())
    
    // Custom headers
    for k, v := range c.config.Headers {
        req.Header.Set(k, v)
    }
    
    if opts != nil {
        for k, v := range opts.Headers {
            req.Header.Set(k, v)
        }
    }
}
```

### 3.4 Retry Logic
**Retry Implementation**:
```go
package retry

import (
    "context"
    "math"
    "math/rand"
    "net/http"
    "strconv"
    "time"
)

const (
    InitialDelay = 500 * time.Millisecond
    MaxDelay     = 8 * time.Second
)

type Config struct {
    MaxRetries int
    ShouldRetry func(*http.Response) bool
}

func DefaultShouldRetry(resp *http.Response) bool {
    // Check x-should-retry header
    if retry := resp.Header.Get("x-should-retry"); retry != "" {
        return retry == "true"
    }
    
    // Retry on specific status codes
    switch resp.StatusCode {
    case 408, 409, 429:  // Timeout, Conflict, Rate Limit
        return true
    }
    
    return resp.StatusCode >= 500
}

func CalculateBackoff(attempt int, resp *http.Response) time.Duration {
    // Check Retry-After header
    if retryAfter := resp.Header.Get("retry-after"); retryAfter != "" {
        if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds > 0 && seconds <= 60 {
            return time.Duration(seconds) * time.Second
        }
        // Could also parse HTTP-date format here
    }
    
    // Check retry-after-ms header (non-standard)
    if retryMs := resp.Header.Get("retry-after-ms"); retryMs != "" {
        if ms, err := strconv.Atoi(retryMs); err == nil {
            return time.Duration(ms) * time.Millisecond
        }
    }
    
    // Exponential backoff with jitter
    delay := math.Min(float64(InitialDelay)*math.Pow(2, float64(attempt)), float64(MaxDelay))
    jitter := 1 - 0.25*rand.Float64()
    return time.Duration(delay * jitter)
}

func Do(ctx context.Context, cfg Config, fn func() (*http.Response, error)) (*http.Response, error) {
    var resp *http.Response
    var err error
    
    for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
        resp, err = fn()
        
        if err != nil || !cfg.ShouldRetry(resp) {
            return resp, err
        }
        
        if attempt < cfg.MaxRetries {
            backoff := CalculateBackoff(attempt, resp)
            select {
            case <-time.After(backoff):
            case <-ctx.Done():
                return nil, ctx.Err()
            }
        }
    }
    
    return resp, err
}
```

## 4. API Surface Parity

### 4.1 Complete API Inventory
**Chat Completions**:@docs/groq-0.34.1/api.md#15-47
- `POST /openai/v1/chat/completions`
- Types: `ChatCompletion`, `ChatCompletionChunk`, `ChatCompletionMessage`, `ChatCompletionMessageParam`, `ChatCompletionTool`, `ChatCompletionToolChoiceOption`, etc.
- 30+ parameters: messages, model, temperature, max_tokens, stream, tools, tool_choice, response_format, etc.
- Special features: function calling, structured outputs (JSON schema), reasoning modes, citations, web search

**Embeddings**:@docs/groq-0.34.1/api.md#49-59
- `POST /openai/v1/embeddings`
- Types: `CreateEmbeddingResponse`, `Embedding`
- Parameters: input, model, encoding_format, user

**Audio**:@docs/groq-0.34.1/api.md#61-91
- Speech: `POST /openai/v1/audio/speech` → `BinaryAPIResponse`
- Transcriptions: `POST /openai/v1/audio/transcriptions` → `Transcription`
- Translations: `POST /openai/v1/audio/translations` → `Translation`

**Models**:@docs/groq-0.34.1/api.md#93-105
- `GET /openai/v1/models/{model}` → `Model`
- `GET /openai/v1/models` → `ModelListResponse`
- `DELETE /openai/v1/models/{model}` → `ModelDeleted`

**Batches**:@docs/groq-0.34.1/api.md#107-125
- `POST /openai/v1/batches` → `BatchCreateResponse`
- `GET /openai/v1/batches/{batch_id}` → `BatchRetrieveResponse`
- `GET /openai/v1/batches` → `BatchListResponse`
- `POST /openai/v1/batches/{batch_id}/cancel` → `BatchCancelResponse`

**Files**:@docs/groq-0.34.1/api.md#127-142
- `POST /openai/v1/files` → `FileCreateResponse`
- `GET /openai/v1/files` → `FileListResponse`
- `DELETE /openai/v1/files/{file_id}` → `FileDeleteResponse`
- `GET /openai/v1/files/{file_id}/content` → `BinaryAPIResponse`
- `GET /openai/v1/files/{file_id}` → `FileInfoResponse`

### 4.2 Type Mapping Strategy
**Python → Go Type Conversions**:
- `str` → `string`
- `int` → `int` or `int64` (context-dependent)
- `float` → `float64`
- `bool` → `bool`
- `List[T]` → `[]T`
- `Dict[str, T]` → `map[string]T`
- `Union[A, B]` → Interface or type switch
- `Optional[T]` → `*T` (pointer for optionality)
- `Literal["a", "b"]` → Custom string type with validation
- `TypedDict` → `struct` with JSON tags

**Omit Pattern**:
Python uses `Omit` sentinel to distinguish unset from `None`. Go equivalent:
```go
type Optional[T any] struct {
    Value T
    Set   bool
}

func Some[T any](v T) Optional[T] {
    return Optional[T]{Value: v, Set: true}
}

func None[T any]() Optional[T] {
    return Optional[T]{Set: false}
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
    if !o.Set {
        return []byte("null"), nil
    }
    return json.Marshal(o.Value)
}
```

### 4.3 Method Signature Pattern
**Example: Chat Completions Create**:
```go
type CreateChatCompletionRequest struct {
    Messages    []ChatCompletionMessageParam `json:"messages"`
    Model       string                        `json:"model"`
    
    // Optional parameters
    Temperature      Optional[float64]                `json:"temperature,omitempty"`
    MaxTokens        Optional[int]                    `json:"max_tokens,omitempty"`
    Stream           Optional[bool]                   `json:"stream,omitempty"`
    Tools            Optional[[]ChatCompletionTool]   `json:"tools,omitempty"`
    ToolChoice       Optional[ToolChoiceOption]       `json:"tool_choice,omitempty"`
    ResponseFormat   Optional[ResponseFormat]         `json:"response_format,omitempty"`
    Stop             Optional[[]string]               `json:"stop,omitempty"`
    Seed             Optional[int]                    `json:"seed,omitempty"`
    // ... 20+ more optional fields
}

func (c *Completions) Create(ctx context.Context, req *CreateChatCompletionRequest, opts ...RequestOption) (*ChatCompletion, error) {
    if req.Stream.Set && req.Stream.Value {
        return nil, errors.New("use CreateStream for streaming requests")
    }
    
    var result ChatCompletion
    err := c.client.post(ctx, "/openai/v1/chat/completions", req, &result, opts...)
    return &result, err
}

func (c *Completions) CreateStream(ctx context.Context, req *CreateChatCompletionRequest, opts ...RequestOption) (*Stream[ChatCompletionChunk], error) {
    req.Stream = Some(true)
    return c.client.postStream(ctx, "/openai/v1/chat/completions", req, opts...)
}
```

### 4.4 Per-Request Options
**RequestOption Pattern**:
```go
type RequestOption func(*RequestOptions)

func WithTimeout(d time.Duration) RequestOption {
    return func(o *RequestOptions) { o.Timeout = &d }
}

func WithMaxRetries(n int) RequestOption {
    return func(o *RequestOptions) { o.MaxRetries = &n }
}

func WithHeader(key, value string) RequestOption {
    return func(o *RequestOptions) {
        if o.Headers == nil {
            o.Headers = make(map[string]string)
        }
        o.Headers[key] = value
    }
}

func WithQuery(key, value string) RequestOption {
    return func(o *RequestOptions) {
        if o.QueryParams == nil {
            o.QueryParams = make(map[string]string)
        }
        o.QueryParams[key] = value
    }
}

// Usage
resp, err := client.Chat.Completions.Create(ctx, req,
    WithTimeout(30*time.Second),
    WithMaxRetries(5),
    WithHeader("X-Custom", "value"),
)
```

## 5. Streaming & SSE Implementation

### 5.1 SSE Decoder (internal/sse)
```go
type Event struct {
    Event string
    Data  string
    ID    string
    Retry int
}

type Decoder struct {
    event      string
    data       []string
    id         string
    retry      int
}

func (d *Decoder) Decode(r io.Reader) (<-chan Event, <-chan error) {
    events := make(chan Event)
    errors := make(chan error, 1)
    
    go func() {
        defer close(events)
        defer close(errors)
        
        scanner := bufio.NewScanner(r)
        var buffer []byte
        
        for scanner.Scan() {
            line := scanner.Bytes()
            buffer = append(buffer, line...)
            buffer = append(buffer, '\n')
            
            // Check for chunk boundary
            if bytes.HasSuffix(buffer, []byte("\n\n")) ||
               bytes.HasSuffix(buffer, []byte("\r\r")) ||
               bytes.HasSuffix(buffer, []byte("\r\n\r\n")) {
                
                event := d.parseChunk(buffer)
                if event != nil {
                    events <- *event
                }
                buffer = buffer[:0]
            }
        }
        
        if err := scanner.Err(); err != nil {
            errors <- err
        }
    }()
    
    return events, errors
}
```

### 5.2 Stream Type
```go
type Stream[T any] struct {
    decoder *sse.Decoder
    resp    *http.Response
    events  <-chan sse.Event
    errors  <-chan error
}

func (s *Stream[T]) Next(ctx context.Context) (*T, error) {
    select {
    case event, ok := <-s.events:
        if !ok {
            return nil, io.EOF
        }
        
        if strings.HasPrefix(event.Data, "[DONE]") {
            return nil, io.EOF
        }
        
        var result T
        if err := json.Unmarshal([]byte(event.Data), &result); err != nil {
            return nil, fmt.Errorf("unmarshal SSE data: %w", err)
        }
        return &result, nil
        
    case err := <-s.errors:
        return nil, err
        
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func (s *Stream[T]) Close() error {
    return s.resp.Body.Close()
}
```

## 6. Error Handling

### 6.1 Error Type Hierarchy
```go
// Base error
type Error struct {
    Message string
    Request *http.Request
}

func (e *Error) Error() string { return e.Message }

// API errors with response
type APIError struct {
    Error
    Response   *http.Response
    StatusCode int
    Body       interface{}  // Parsed JSON or raw string
}

// Specific status errors
type BadRequestError struct{ APIError }          // 400
type AuthenticationError struct{ APIError }      // 401
type PermissionDeniedError struct{ APIError }    // 403
type NotFoundError struct{ APIError }            // 404
type ConflictError struct{ APIError }            // 409
type UnprocessableEntityError struct{ APIError } // 422
type RateLimitError struct{ APIError }           // 429
type InternalServerError struct{ APIError }      // >=500

// Connection errors
type ConnectionError struct{ Error }
type TimeoutError struct{ ConnectionError }
type ValidationError struct{ Error }
```

## 7. Complete Implementation Roadmap

### Phase 1: Foundation (Week 1-2)
**Deliverables**:
- [x] Module setup: `go.mod`, directory structure, CI/CD (GitHub Actions)
- [x] Core client: `Client`, `ClientConfig`, `NewClient()` with env loading
- [x] HTTP client wrapper with default timeouts/limits
- [x] Error hierarchy: all error types with proper wrapping
- [x] Constants: timeouts, limits, headers
- [x] Logging: structured logger interface + env toggle (`GROQ_LOG`)
- [x] Unit tests: client initialization, config validation

### Phase 2: Request Infrastructure (Week 2-3)
**Deliverables**:
- [x] Request building: headers, auth, URL construction
- [x] Retry logic: exponential backoff, jitter, `Retry-After` parsing
- [x] Query string encoding: comma arrays, nested objects
- [x] Request options: functional options pattern
- [x] Response parsing: JSON unmarshaling, content-type validation
- [x] Unit tests: retry scenarios, query encoding, header merging

### Phase 3: Type System (Week 3-4)
**Deliverables**:
- [x] `Optional[T]` generic with JSON marshaling
- [x] Chat types: `ChatCompletion`, `ChatCompletionChunk`, message types
- [x] Embedding types: `CreateEmbeddingResponse`, `Embedding`
- [x] Model types: `Model`, `ModelListResponse`, `ModelDeleted`
- [x] Shared types: `ErrorObject`, `FunctionDefinition`, `CompletionUsage`
- [x] Enums: `Role`, `FinishReason`, model literals
- [x] Unit tests: JSON round-trip, optional field handling

### Phase 4: Chat Completions (Week 4-5)
**Deliverables**:
- [x] `chat.Completions` service struct
- [x] `Create()` method with full parameter support
- [x] `CreateStream()` for streaming responses
- [x] SSE decoder implementation
- [x] Stream type with `Next()` iterator
- [x] Integration tests with real API
- [x] Examples: basic completion, streaming, function calling

### Phase 5: Remaining Resources (Week 5-7)
**Deliverables**:
- [x] Embeddings service
- [x] Audio: speech, transcriptions, translations
- [x] Models: list, retrieve, delete
- [x] Batches: CRUD operations
- [x] Files: upload, download, list, delete
- [x] Multipart form-data handling
- [x] Binary response handling
- [x] Integration tests for each resource

### Phase 6: Advanced Features (Week 7-8)
**Deliverables**:
- [ ] Raw response wrappers
- [ ] Pagination helpers (if needed)
- [ ] Context cancellation throughout
- [ ] Idempotency key generation
- [ ] Platform headers (OS, arch, Go version)
- [ ] Comprehensive examples
- [ ] Performance benchmarks

### Phase 7: Documentation & Release (Week 8-9)
**Deliverables**:
- [ ] README with quickstart, examples, API reference
- [ ] GoDoc comments on all public APIs
- [ ] Migration guide from Python SDK
- [ ] CHANGELOG following semver
- [ ] GitHub release with binaries
- [ ] Announce v0.1.0-alpha

## 8. Testing Strategy

### 8.1 Unit Tests
- Mock HTTP client for deterministic tests
- Cover all error paths, retry scenarios
- Validate JSON marshaling/unmarshaling
- Test optional field handling
- Query string encoding edge cases

### 8.2 Integration Tests
- Gated by `GROQ_API_KEY` env var
- Test each resource endpoint
- Verify streaming works end-to-end
- File upload/download round-trip
- Rate limit handling

### 8.3 Contract Tests
- JSON fixtures from Python SDK
- Ensure Go types match Python responses
- Validate request serialization matches Python

### 8.4 Benchmarks
- Chat completion throughput
- Streaming latency
- Memory allocation profiling

## 9. Open Questions & Risks

1. **OpenAPI Spec Access**: Need spec for automated type generation
2. **Discriminated Unions**: Go lacks native support; use type switches or interfaces
3. **Streaming Error Handling**: Ensure errors mid-stream are surfaced properly
4. **Context Propagation**: Verify all blocking calls respect context cancellation
5. **Binary Responses**: Audio speech endpoint returns raw bytes; need special handling
6. **Pagination**: Check if any endpoints use cursor-based pagination

---

**Total Estimated Effort**: 9 weeks (1 developer)  
**Critical Path**: Type system → Chat completions → Streaming  
**Success Metrics**: 100% API coverage, <100ms overhead vs Python, >90% test coverage

---

## Appendix: Python SDK Technical Reference

### A.1 Dependencies
**Core**:@docs/groq-0.34.1/pyproject.toml#10-17
- `httpx>=0.23.0, <1` - HTTP client
- `pydantic>=1.9.0, <3` - Data validation
- `typing-extensions>=4.10, <5` - Backports
- `anyio>=3.5.0, <5` - Async I/O
- `distro>=1.7.0, <2` - OS detection
- `sniffio` - Async library detection

**Dev**:@docs/groq-0.34.1/pyproject.toml#46-59
- `pyright==1.1.399`, `mypy` - Type checking
- `ruff` - Linting/formatting
- `pytest`, `pytest-asyncio`, `pytest-xdist` - Testing
- `respx` - HTTP mocking
- `nox` - Test automation

### A.2 Example Usage Patterns
**Basic Chat**:@docs/groq-0.34.1/examples/chat_completion.py#1-46
```python
client = Groq()
completion = client.chat.completions.create(
    messages=[{"role": "user", "content": "Hello"}],
    model="mixtral-8x7b-32768",
    temperature=0.5,
    max_tokens=1024,
)
print(completion.choices[0].message.content)
```

**Streaming**:@docs/groq-0.34.1/examples/chat_completion_streaming.py#5-53
```python
stream = client.chat.completions.create(stream=True, ...)
for chunk in stream:
    print(chunk.choices[0].delta.content, end="")
    if chunk.choices[0].finish_reason:
        print(chunk.x_groq.usage)
```

### A.3 Key Implementation Details
- **Stainless Generated**: SDK auto-generated from OpenAPI spec@docs/groq-0.34.1/README.md#10
- **Pydantic v1/v2 Compat**: Supports both via compatibility layer@docs/groq-0.34.1/src/groq/_models.py#84-384
- **Cached Properties**: Resources lazy-loaded via `@cached_property`@docs/groq-0.34.1/src/groq/_client.py#100-134
- **Response Validation**: Optional strict mode raises errors on schema mismatch@docs/groq-0.34.1/src/groq/_client.py#70
- **Platform Headers**: Sends OS, Python version, architecture@docs/groq-0.34.1/src/groq/_base_client.py#684-688

---

**Document Version**: 1.0  
**Last Updated**: 2024-11-18  
**Python SDK Version**: 0.34.1  
**Total Lines Analyzed**: ~10,000+ from Python SDK
