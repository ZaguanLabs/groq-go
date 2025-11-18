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
	Role         Role          `json:"role"`
	Content      string        `json:"content"` // Can be null? Usually string or null.
	Refusal      string        `json:"refusal,omitempty"`
	ToolCalls    []ToolCall    `json:"tool_calls,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
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
	Usage             *CompletionUsage            `json:"usage,omitempty"`  // Optional in chunks? Yes, usually final chunk.
	XGroq             *XGroqUsage                 `json:"x_groq,omitempty"` // Groq specific usage in stream?
}

// XGroqUsage represents Groq specific usage info
type XGroqUsage struct {
	Usage *CompletionUsage `json:"usage,omitempty"`
	ID    string           `json:"id,omitempty"` // Request ID?
	Error *ErrorObject     `json:"error,omitempty"`
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
	Role         Role          `json:"role,omitempty"`
	Content      string        `json:"content,omitempty"`
	Refusal      string        `json:"refusal,omitempty"`
	ToolCalls    []ToolCall    `json:"tool_calls,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

// ChatCompletionMessageParam represents an input message
// We use a struct with pointers for optional fields to support omitempty correctly via standard json
// or just strict fields.
// But wait, Content can be []ContentPart for vision.
// For now, let's assume text content.
type ChatCompletionMessageParam struct {
	Role       Role        `json:"role"`
	Content    interface{} `json:"content"` // string or []ContentPart
	Name       string      `json:"name,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

// CreateChatCompletionRequest represents the request body
type CreateChatCompletionRequest struct {
	Messages         []ChatCompletionMessageParam `json:"messages"`
	Model            string                       `json:"model"`
	FrequencyPenalty *option.Optional[float64]    `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int               `json:"logit_bias,omitempty"`
	Logprobs         *option.Optional[bool]       `json:"logprobs,omitempty"`
	TopLogprobs      *option.Optional[int]        `json:"top_logprobs,omitempty"`
	MaxTokens        *option.Optional[int]        `json:"max_tokens,omitempty"`
	N                *option.Optional[int]        `json:"n,omitempty"`
	PresencePenalty  *option.Optional[float64]    `json:"presence_penalty,omitempty"`
	ResponseFormat   *ResponseFormat              `json:"response_format,omitempty"`
	Seed             *option.Optional[int]        `json:"seed,omitempty"`
	Stop             interface{}                  `json:"stop,omitempty"` // string or []string
	Stream           *option.Optional[bool]       `json:"stream,omitempty"`
	Temperature      *option.Optional[float64]    `json:"temperature,omitempty"`
	TopP             *option.Optional[float64]    `json:"top_p,omitempty"`
	Tools            []ChatCompletionTool         `json:"tools,omitempty"`
	ToolChoice       interface{}                  `json:"tool_choice,omitempty"` // string or ToolChoice object
	User             string                       `json:"user,omitempty"`

	// Groq specific?
	ParallelToolCalls *option.Optional[bool] `json:"parallel_tool_calls,omitempty"`
}

// ResponseFormat represents the response format
type ResponseFormat struct {
	Type string `json:"type"` // "text" or "json_object"
}

// ChatCompletionTool represents a tool definition
type ChatCompletionTool struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}
