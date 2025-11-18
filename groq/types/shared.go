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
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	PromptTime       float64 `json:"prompt_time,omitempty"`     // Groq specific
	CompletionTime   float64 `json:"completion_time,omitempty"` // Groq specific
	TotalTime        float64 `json:"total_time,omitempty"`      // Groq specific
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
