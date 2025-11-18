package types

import (
	"github.com/ZaguanLabs/groq-go/groq/option"
)

// CreateEmbeddingRequest represents the request body for embeddings
type CreateEmbeddingRequest struct {
	Input          interface{}              `json:"input"` // string or []string
	Model          string                   `json:"model"`
	EncodingFormat *option.Optional[string] `json:"encoding_format,omitempty"` // "float" or "base64"
	User           string                   `json:"user,omitempty"`
}

// CreateEmbeddingResponse represents the response for embeddings
type CreateEmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []Embedding     `json:"data"`
	Model  string          `json:"model"`
	Usage  CompletionUsage `json:"usage"`
}

// Embedding represents a single embedding
type Embedding struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
	Object    string    `json:"object"`
}
