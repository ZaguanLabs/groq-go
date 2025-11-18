package groq

import (
	"fmt"
	"net/http"
)

// GroqError is the base error type for the SDK
type GroqError struct {
	Message string
	Request *http.Request
}

func (e *GroqError) Error() string { return e.Message }

// APIError represents an error returned by the API
type APIError struct {
	GroqError
	Response   *http.Response
	StatusCode int
	Body       interface{} // Parsed JSON or raw string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Error code: %d - %v", e.StatusCode, e.Body)
}

// Specific status errors

// BadRequestError corresponds to 400 Bad Request
type BadRequestError struct{ APIError }

// AuthenticationError corresponds to 401 Unauthorized
type AuthenticationError struct{ APIError }

// PermissionDeniedError corresponds to 403 Forbidden
type PermissionDeniedError struct{ APIError }

// NotFoundError corresponds to 404 Not Found
type NotFoundError struct{ APIError }

// ConflictError corresponds to 409 Conflict
type ConflictError struct{ APIError }

// UnprocessableEntityError corresponds to 422 Unprocessable Entity
type UnprocessableEntityError struct{ APIError }

// RateLimitError corresponds to 429 Too Many Requests
type RateLimitError struct{ APIError }

// InternalServerError corresponds to 500+ Server Errors
type InternalServerError struct{ APIError }

// ConnectionError represents a connection error
type ConnectionError struct{ GroqError }

// TimeoutError represents a timeout error
type TimeoutError struct{ ConnectionError }

// ValidationError represents a validation error
type ValidationError struct{ GroqError }
