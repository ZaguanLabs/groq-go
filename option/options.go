package option

import (
	"time"
)

// RequestOptions holds per-request configuration
type RequestOptions struct {
	Headers     map[string]string
	QueryParams map[string]string
	Timeout     *time.Duration
	MaxRetries  *int

	// Internal
	IdempotencyKey string
}

// RequestOption allows configuring a request
type RequestOption func(*RequestOptions)

// WithRequestTimeout sets the timeout for a specific request
func WithRequestTimeout(d time.Duration) RequestOption {
	return func(o *RequestOptions) { o.Timeout = &d }
}

// WithRequestMaxRetries sets the max retries for a specific request
func WithRequestMaxRetries(n int) RequestOption {
	return func(o *RequestOptions) { o.MaxRetries = &n }
}

// WithRequestHeader adds a header to a specific request
func WithRequestHeader(key, value string) RequestOption {
	return func(o *RequestOptions) {
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		o.Headers[key] = value
	}
}

// WithRequestQuery adds a query parameter to a specific request
func WithRequestQuery(key, value string) RequestOption {
	return func(o *RequestOptions) {
		if o.QueryParams == nil {
			o.QueryParams = make(map[string]string)
		}
		o.QueryParams[key] = value
	}
}
