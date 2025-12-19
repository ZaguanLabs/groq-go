package types

import (
	"github.com/ZaguanLabs/groq-go/groq/option"
)

// ChatCompletion represents a chat completion response
type ChatCompletion struct {
	ID                string                 `json:"id"`
	Choices           []ChatCompletionChoice `json:"choices"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	SystemFingerprint string                 `json:"system_fingerprint,omitempty"`
	Object            string                 `json:"object"`
	Usage             *CompletionUsage       `json:"usage,omitempty"`
	McpListTools      []McpListTool          `json:"mcp_list_tools,omitempty"`  // MCP tool discovery
	ServiceTier       string                 `json:"service_tier,omitempty"`    // Service tier used (auto, on_demand, flex, performance)
	UsageBreakdown    *UsageBreakdown        `json:"usage_breakdown,omitempty"` // Per-model usage for compound AI
	XGroq             *XGroq                 `json:"x_groq,omitempty"`          // Groq-specific metadata
}

// McpListTool represents MCP tool discovery information
type McpListTool struct {
	ID          string           `json:"id,omitempty"`
	ServerLabel string           `json:"server_label,omitempty"`
	Tools       []McpListToolDef `json:"tools,omitempty"`
	Type        string           `json:"type,omitempty"`
}

// McpListToolDef represents a single MCP tool definition
type McpListToolDef struct {
	Annotations interface{}            `json:"annotations,omitempty"` // Additional metadata
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"` // JSON Schema
	Name        string                 `json:"name,omitempty"`
}

// ChatCompletionChoice represents a choice in a chat completion
type ChatCompletionChoice struct {
	FinishReason FinishReason            `json:"finish_reason"`
	Index        int                     `json:"index"`
	Logprobs     *ChatCompletionLogprobs `json:"logprobs,omitempty"`
	Message      ChatCompletionMessage   `json:"message"`
}

// ChatCompletionLogprobs represents log probability information
type ChatCompletionLogprobs struct {
	Content []ChatCompletionTokenLogprob `json:"content"`
}

// ChatCompletionTokenLogprob represents a token logprob
type ChatCompletionTokenLogprob struct {
	Token       string       `json:"token"`
	Logprob     float64      `json:"logprob"`
	Bytes       []int        `json:"bytes,omitempty"`
	TopLogprobs []TopLogprob `json:"top_logprobs,omitempty"`
}

// TopLogprob represents a top logprob
type TopLogprob struct {
	Token   string  `json:"token"`
	Logprob float64 `json:"logprob"`
	Bytes   []int   `json:"bytes,omitempty"`
}

// ChatCompletionMessage represents a message in a chat completion response
type ChatCompletionMessage struct {
	Role          Role           `json:"role"`
	Content       string         `json:"content"` // Can be null? Usually string or null.
	Refusal       string         `json:"refusal,omitempty"`
	ToolCalls     []ToolCall     `json:"tool_calls,omitempty"`
	FunctionCall  *FunctionCall  `json:"function_call,omitempty"`
	Annotations   []Annotation   `json:"annotations,omitempty"`
	ExecutedTools []ExecutedTool `json:"executed_tools,omitempty"`
	Reasoning     *string        `json:"reasoning,omitempty"`
}

// ToolCall represents a tool call
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// ChatCompletionChunk represents a streamed chat completion chunk
type ChatCompletionChunk struct {
	ID                string                      `json:"id"`
	Choices           []ChatCompletionChunkChoice `json:"choices"`
	Created           int64                       `json:"created"`
	Model             string                      `json:"model"`
	SystemFingerprint string                      `json:"system_fingerprint,omitempty"`
	Object            string                      `json:"object"`
	Usage             *CompletionUsage            `json:"usage,omitempty"`  // Only in final chunk when stream_options.include_usage is true
	XGroq             *XGroqStream                `json:"x_groq,omitempty"` // Groq-specific metadata in streaming responses
}

// XGroq represents Groq-specific metadata in non-streaming responses
type XGroq struct {
	ID    string      `json:"id"`
	Debug *XGroqDebug `json:"debug,omitempty"`
	Seed  *int        `json:"seed,omitempty"`
	Usage *XGroqUsage `json:"usage,omitempty"` // Hardware cache statistics
}

// XGroqStream represents Groq-specific metadata in streaming responses
type XGroqStream struct {
	ID             *string          `json:"id,omitempty"`              // Sent in first and final chunk
	Debug          *XGroqDebug      `json:"debug,omitempty"`           // Only when debug=true
	Seed           *int             `json:"seed,omitempty"`            // Sent in final chunk
	Usage          *CompletionUsage `json:"usage,omitempty"`           // Usage stats in final chunk
	UsageBreakdown *UsageBreakdown  `json:"usage_breakdown,omitempty"` // Per-model usage for compound AI
	Error          *string          `json:"error,omitempty"`           // Error if stream stopped early
}

// XGroqUsage represents Groq-specific hardware cache statistics
type XGroqUsage struct {
	DramCachedTokens *int `json:"dram_cached_tokens,omitempty"` // Number of tokens served from DRAM cache
	SramCachedTokens *int `json:"sram_cached_tokens,omitempty"` // Number of tokens served from SRAM cache
}

// XGroqDebug represents debug information
type XGroqDebug struct {
	InputTokenIDs  []int    `json:"input_token_ids,omitempty"`
	InputTokens    []string `json:"input_tokens,omitempty"`
	OutputTokenIDs []int    `json:"output_token_ids,omitempty"`
	OutputTokens   []string `json:"output_tokens,omitempty"`
}

// UsageBreakdown represents per-model usage statistics
type UsageBreakdown struct {
	Models []UsageBreakdownModel `json:"models"`
}

// UsageBreakdownModel represents usage for a single model
type UsageBreakdownModel struct {
	Model string          `json:"model"`
	Usage CompletionUsage `json:"usage"`
}

// ChatCompletionChunkChoice represents a choice in a chunk
type ChatCompletionChunkChoice struct {
	Delta        ChatCompletionChunkDelta `json:"delta"`
	FinishReason FinishReason             `json:"finish_reason"` // Can be null/empty
	Index        int                      `json:"index"`
	Logprobs     *ChatCompletionLogprobs  `json:"logprobs,omitempty"`
}

// ChatCompletionChunkDelta represents a delta in a chunk
type ChatCompletionChunkDelta struct {
	Role          Role           `json:"role,omitempty"`
	Content       string         `json:"content,omitempty"`
	Refusal       string         `json:"refusal,omitempty"`
	ToolCalls     []ToolCall     `json:"tool_calls,omitempty"`
	FunctionCall  *FunctionCall  `json:"function_call,omitempty"`
	Annotations   []Annotation   `json:"annotations,omitempty"`
	ExecutedTools []ExecutedTool `json:"executed_tools,omitempty"`
	Reasoning     *string        `json:"reasoning,omitempty"`
}

// ChatCompletionMessageParam represents an input message
// Content can be:
// - string for simple text messages
// - []ContentPart for multimodal messages (text, images, documents)
type ChatCompletionMessageParam struct {
	Role         Role          `json:"role"`
	Content      interface{}   `json:"content"`                 // string or []ContentPart
	Name         string        `json:"name,omitempty"`          // Optional name for participant
	ToolCalls    []ToolCall    `json:"tool_calls,omitempty"`    // For assistant messages with tool calls
	ToolCallID   string        `json:"tool_call_id,omitempty"`  // For tool messages
	FunctionCall *FunctionCall `json:"function_call,omitempty"` // Deprecated: use tool_calls instead
	Reasoning    *string       `json:"reasoning,omitempty"`     // For assistant messages with reasoning_format=parsed
}

// CreateChatCompletionRequest represents the request body
type CreateChatCompletionRequest struct {
	Messages              []ChatCompletionMessageParam `json:"messages"`
	Model                 string                       `json:"model"`
	FrequencyPenalty      *option.Optional[float64]    `json:"frequency_penalty,omitempty"`
	LogitBias             map[string]int               `json:"logit_bias,omitempty"`
	Logprobs              *option.Optional[bool]       `json:"logprobs,omitempty"`
	TopLogprobs           *option.Optional[int]        `json:"top_logprobs,omitempty"`
	MaxTokens             *option.Optional[int]        `json:"max_tokens,omitempty"`
	MaxCompletionTokens   *option.Optional[int]        `json:"max_completion_tokens,omitempty"`
	N                     *option.Optional[int]        `json:"n,omitempty"`
	PresencePenalty       *option.Optional[float64]    `json:"presence_penalty,omitempty"`
	ResponseFormat        *ResponseFormat              `json:"response_format,omitempty"`
	Seed                  *option.Optional[int]        `json:"seed,omitempty"`
	Stop                  interface{}                  `json:"stop,omitempty"` // string or []string
	Stream                *option.Optional[bool]       `json:"stream,omitempty"`
	Temperature           *option.Optional[float64]    `json:"temperature,omitempty"`
	TopP                  *option.Optional[float64]    `json:"top_p,omitempty"`
	Tools                 []ChatCompletionTool         `json:"tools,omitempty"`
	ToolChoice            interface{}                  `json:"tool_choice,omitempty"` // string or ToolChoice object
	User                  string                       `json:"user,omitempty"`
	ParallelToolCalls     *option.Optional[bool]       `json:"parallel_tool_calls,omitempty"`
	DisableToolValidation bool                         `json:"disable_tool_validation,omitempty"`

	// Deprecated: use Tools instead
	Functions    []FunctionDefinition `json:"functions,omitempty"`
	FunctionCall interface{}          `json:"function_call,omitempty"` // "none", "auto", "required", or {"name": "func_name"}

	// Compound AI
	CompoundCustom *CompoundCustom `json:"compound_custom,omitempty"`

	// Documents and Citations
	Documents       []Document               `json:"documents,omitempty"`
	CitationOptions *option.Optional[string] `json:"citation_options,omitempty"`

	// Reasoning
	ReasoningEffort  *option.Optional[string] `json:"reasoning_effort,omitempty"`
	ReasoningFormat  *option.Optional[string] `json:"reasoning_format,omitempty"`
	IncludeReasoning *option.Optional[bool]   `json:"include_reasoning,omitempty"`

	// Search
	SearchSettings *SearchSettings `json:"search_settings,omitempty"`
	ExcludeDomains []string        `json:"exclude_domains,omitempty"` // Deprecated
	IncludeDomains []string        `json:"include_domains,omitempty"` // Deprecated

	// Service tier
	ServiceTier *option.Optional[string] `json:"service_tier,omitempty"`

	// Metadata (not currently supported)
	Metadata map[string]string      `json:"metadata,omitempty"`
	Store    *option.Optional[bool] `json:"store,omitempty"`
}

// ResponseFormat represents the response format
// Type can be "text", "json_object", or "json_schema"
type ResponseFormat struct {
	Type       string                    `json:"type"`                  // "text", "json_object", or "json_schema"
	JSONSchema *ResponseFormatJSONSchema `json:"json_schema,omitempty"` // Required when type is "json_schema"
}

// ResponseFormatJSONSchema represents structured output configuration
type ResponseFormatJSONSchema struct {
	Name        string                 `json:"name"`                  // Schema name (a-z, A-Z, 0-9, underscores, dashes, max 64 chars)
	Description string                 `json:"description,omitempty"` // Description of the response format
	Schema      map[string]interface{} `json:"schema,omitempty"`      // JSON Schema object
	Strict      *bool                  `json:"strict,omitempty"`      // Enable strict schema adherence
}

// ChatCompletionTool represents a tool definition
type ChatCompletionTool struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}

// CompoundCustom represents custom configuration for Compound AI
type CompoundCustom struct {
	Models *CompoundCustomModels `json:"models,omitempty"`
	Tools  *CompoundCustomTools  `json:"tools,omitempty"`
}

// CompoundCustomModels represents custom model selection
type CompoundCustomModels struct {
	AnsweringModel *option.Optional[string] `json:"answering_model,omitempty"`
	ReasoningModel *option.Optional[string] `json:"reasoning_model,omitempty"`
}

// CompoundCustomTools represents tool configuration
type CompoundCustomTools struct {
	EnabledTools    []string                            `json:"enabled_tools,omitempty"`
	WolframSettings *CompoundCustomToolsWolframSettings `json:"wolfram_settings,omitempty"`
}

// CompoundCustomToolsWolframSettings represents Wolfram tool settings
type CompoundCustomToolsWolframSettings struct {
	Authorization *option.Optional[string] `json:"authorization,omitempty"`
}

// Document represents a document for context
type Document struct {
	ID     *string         `json:"id,omitempty"`
	Source *DocumentSource `json:"source"`
}

// DocumentSource represents the source of a document
type DocumentSource struct {
	Type string                 `json:"type"` // "text" or "json"
	Text *string                `json:"text,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// Annotation represents a citation or reference
type Annotation struct {
	Type             string                      `json:"type"` // "document_citation" or "function_citation"
	DocumentCitation *AnnotationDocumentCitation `json:"document_citation,omitempty"`
	FunctionCitation *AnnotationFunctionCitation `json:"function_citation,omitempty"`
}

// AnnotationDocumentCitation represents a document citation
type AnnotationDocumentCitation struct {
	DocumentID string `json:"document_id"`
	StartIndex int    `json:"start_index"`
	EndIndex   int    `json:"end_index"`
}

// AnnotationFunctionCitation represents a function call citation
type AnnotationFunctionCitation struct {
	ToolCallID string `json:"tool_call_id"`
	StartIndex int    `json:"start_index"`
	EndIndex   int    `json:"end_index"`
}

// ExecutedTool represents a tool that was executed
type ExecutedTool struct {
	Arguments      string                      `json:"arguments"`
	Index          int                         `json:"index"`
	Type           string                      `json:"type"`
	BrowserResults []ExecutedToolBrowserResult `json:"browser_results,omitempty"`
	CodeResults    []ExecutedToolCodeResult    `json:"code_results,omitempty"`
	Output         *string                     `json:"output,omitempty"`
	SearchResults  *ExecutedToolSearchResults  `json:"search_results,omitempty"`
}

// ExecutedToolBrowserResult represents browser execution result
type ExecutedToolBrowserResult struct {
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Content     *string `json:"content,omitempty"`
	LiveViewURL *string `json:"live_view_url,omitempty"`
}

// ExecutedToolCodeResult represents code execution result
type ExecutedToolCodeResult struct {
	Chart  *ExecutedToolCodeResultChart  `json:"chart,omitempty"`
	Charts []ExecutedToolCodeResultChart `json:"charts,omitempty"`
	PNG    *string                       `json:"png,omitempty"` // Base64 encoded
	Text   *string                       `json:"text,omitempty"`
}

// ExecutedToolCodeResultChart represents a chart from code execution
type ExecutedToolCodeResultChart struct {
	Elements    []ExecutedToolCodeResultChartElement `json:"elements"`
	Type        string                               `json:"type"` // "bar", "box_and_whisker", "line", "pie", "scatter", "superchart", "unknown"
	Title       *string                              `json:"title,omitempty"`
	XLabel      *string                              `json:"x_label,omitempty"`
	XScale      *string                              `json:"x_scale,omitempty"`
	XTickLabels []string                             `json:"x_tick_labels,omitempty"`
	XTicks      []float64                            `json:"x_ticks,omitempty"`
	XUnit       *string                              `json:"x_unit,omitempty"`
	YLabel      *string                              `json:"y_label,omitempty"`
	YScale      *string                              `json:"y_scale,omitempty"`
	YTickLabels []string                             `json:"y_tick_labels,omitempty"`
	YTicks      []float64                            `json:"y_ticks,omitempty"`
	YUnit       *string                              `json:"y_unit,omitempty"`
}

// ExecutedToolCodeResultChartElement represents a chart element
type ExecutedToolCodeResultChartElement struct {
	Label         string      `json:"label"`
	Angle         *float64    `json:"angle,omitempty"`
	FirstQuartile *float64    `json:"first_quartile,omitempty"`
	Group         *string     `json:"group,omitempty"`
	Max           *float64    `json:"max,omitempty"`
	Median        *float64    `json:"median,omitempty"`
	Min           *float64    `json:"min,omitempty"`
	Outliers      []float64   `json:"outliers,omitempty"`
	Points        [][]float64 `json:"points,omitempty"`
	Radius        *float64    `json:"radius,omitempty"`
	ThirdQuartile *float64    `json:"third_quartile,omitempty"`
	Value         *float64    `json:"value,omitempty"`
}

// ExecutedToolSearchResults represents search results
type ExecutedToolSearchResults struct {
	Images  []string                          `json:"images,omitempty"`
	Results []ExecutedToolSearchResultsResult `json:"results,omitempty"`
}

// ExecutedToolSearchResultsResult represents a single search result
type ExecutedToolSearchResultsResult struct {
	Content *string  `json:"content,omitempty"`
	Score   *float64 `json:"score,omitempty"`
	Title   *string  `json:"title,omitempty"`
	URL     *string  `json:"url,omitempty"`
}

// SearchSettings represents web search configuration
type SearchSettings struct {
	Country        *option.Optional[string] `json:"country,omitempty"`
	IncludeDomains []string                 `json:"include_domains,omitempty"`
	ExcludeDomains []string                 `json:"exclude_domains,omitempty"`
	IncludeImages  *option.Optional[bool]   `json:"include_images,omitempty"`
}
