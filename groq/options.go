package groq

import (
	"net/http"
	"time"
)

// ClientConfig holds configuration for the Client
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
	HTTPClient       *http.Client // Optional custom client
}

// ClientOption allows configuring the Client
type ClientOption func(*ClientConfig)

// WithAPIKey sets the API key
func WithAPIKey(key string) ClientOption {
	return func(c *ClientConfig) { c.APIKey = key }
}

// WithBaseURL sets the base URL
func WithBaseURL(url string) ClientOption {
	return func(c *ClientConfig) { c.BaseURL = url }
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(n int) ClientOption {
	return func(c *ClientConfig) { c.MaxRetries = n }
}

// WithTimeout sets the request timeout
func WithTimeout(d time.Duration) ClientOption {
	return func(c *ClientConfig) { c.Timeout = d }
}

// WithHTTPClient sets the custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *ClientConfig) { c.HTTPClient = client }
}

// WithHeader adds a default header to all requests
func WithHeader(key, value string) ClientOption {
	return func(c *ClientConfig) { c.Headers[key] = value }
}

// WithLogger sets the logger
func WithLogger(l Logger) ClientOption {
	return func(c *ClientConfig) { c.Logger = l }
}
