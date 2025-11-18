package groq

import (
	"net/http"
	"time"
)

const (
	// Version is the current version of the SDK
	Version = "0.1.0-alpha"

	// DefaultBaseURL is the default API endpoint
	DefaultBaseURL = "https://api.groq.com"

	// DefaultMaxRetries is the default number of retries for failed requests
	DefaultMaxRetries = 2

	// Headers
	HeaderRawResponse  = "X-Stainless-Raw-Response"
	HeaderOverrideCast = "____stainless_override_cast_to"

	// InitialRetryDelay is the initial backoff duration
	InitialRetryDelay = 500 * time.Millisecond

	// MaxRetryDelay is the maximum backoff duration
	MaxRetryDelay = 8 * time.Second
)

var (
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 60 * time.Second

	// DefaultConnectTimeout is the default connection timeout
	DefaultConnectTimeout = 5 * time.Second

	// DefaultIdleConnTimeout is the default idle connection timeout
	DefaultIdleConnTimeout = 90 * time.Second
)

// DefaultTransport returns a http.Transport with default settings
func DefaultTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     DefaultIdleConnTimeout,
	}
}
