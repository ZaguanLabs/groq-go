package types

// ErrorObject represents an error returned by the API
type ErrorObject struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param,omitempty"`
	Code    interface{} `json:"code,omitempty"`
}

// FunctionDefinition represents a function definition
type FunctionDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"` // JSON Schema object
}

// CompletionUsage represents usage statistics
type CompletionUsage struct {
	PromptTokens            int                      `json:"prompt_tokens"`
	CompletionTokens        int                      `json:"completion_tokens"`
	TotalTokens             int                      `json:"total_tokens"`
	PromptTime              float64                  `json:"prompt_time,omitempty"`               // Groq specific
	CompletionTime          float64                  `json:"completion_time,omitempty"`           // Groq specific
	TotalTime               float64                  `json:"total_time,omitempty"`                // Groq specific
	QueueTime               float64                  `json:"queue_time,omitempty"`                // Time spent in queue
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"` // Breakdown of completion tokens
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`     // Breakdown of prompt tokens
}

// CompletionTokensDetails represents breakdown of completion tokens
type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"` // Number of tokens used for reasoning
}

// PromptTokensDetails represents breakdown of prompt tokens
type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"` // Number of tokens that were cached and reused
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
