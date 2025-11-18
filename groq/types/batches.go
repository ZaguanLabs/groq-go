package types

import "github.com/ZaguanLabs/groq-go/groq/option"

// Batch represents a batch job
type Batch struct {
	ID               string              `json:"id"`
	Object           string              `json:"object"`
	Endpoint         string              `json:"endpoint"`
	Errors           *BatchErrors        `json:"errors,omitempty"`
	InputFileID      string              `json:"input_file_id"`
	CompletionWindow string              `json:"completion_window"`
	Status           string              `json:"status"` // validating, failed, in_progress, finalizing, completed, expired, cancelling, cancelled
	OutputFileID     string              `json:"output_file_id,omitempty"`
	ErrorFileID      string              `json:"error_file_id,omitempty"`
	CreatedAt        int64               `json:"created_at"`
	InProgressAt     int64               `json:"in_progress_at,omitempty"`
	ExpiresAt        int64               `json:"expires_at,omitempty"`
	FinalizingAt     int64               `json:"finalizing_at,omitempty"`
	CompletedAt      int64               `json:"completed_at,omitempty"`
	FailedAt         int64               `json:"failed_at,omitempty"`
	ExpiredAt        int64               `json:"expired_at,omitempty"`
	CancellingAt     int64               `json:"cancelling_at,omitempty"`
	CancelledAt      int64               `json:"cancelled_at,omitempty"`
	RequestCounts    *BatchRequestCounts `json:"request_counts,omitempty"`
	Metadata         map[string]string   `json:"metadata,omitempty"`
}

// BatchErrors represents errors in a batch
type BatchErrors struct {
	Object string        `json:"object"`
	Data   []ErrorObject `json:"data"`
}

// BatchRequestCounts represents request counts in a batch
type BatchRequestCounts struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}

// CreateBatchRequest represents parameters to create a batch
type CreateBatchRequest struct {
	InputFileID      string            `json:"input_file_id"`
	Endpoint         string            `json:"endpoint"`
	CompletionWindow string            `json:"completion_window"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// BatchListResponse represents a list of batches
type BatchListResponse struct {
	Object  string  `json:"object"`
	Data    []Batch `json:"data"`
	FirstID string  `json:"first_id,omitempty"`
	LastID  string  `json:"last_id,omitempty"`
	HasMore bool    `json:"has_more"`
}

// ListBatchesRequest represents parameters to list batches
type ListBatchesRequest struct {
	After *option.Optional[string] `json:"after,omitempty"`
	Limit *option.Optional[int]    `json:"limit,omitempty"`
}
