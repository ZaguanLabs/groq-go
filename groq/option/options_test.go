package option

import (
	"testing"
	"time"
)

func TestWithRequestTimeout(t *testing.T) {
	opts := &RequestOptions{}
	timeout := 10 * time.Second

	WithRequestTimeout(timeout)(opts)

	if opts.Timeout == nil {
		t.Fatal("Timeout not set")
	}
	if *opts.Timeout != timeout {
		t.Errorf("Timeout = %v, want %v", *opts.Timeout, timeout)
	}
}

func TestWithRequestMaxRetries(t *testing.T) {
	opts := &RequestOptions{}
	maxRetries := 5

	WithRequestMaxRetries(maxRetries)(opts)

	if opts.MaxRetries == nil {
		t.Fatal("MaxRetries not set")
	}
	if *opts.MaxRetries != maxRetries {
		t.Errorf("MaxRetries = %d, want %d", *opts.MaxRetries, maxRetries)
	}
}

func TestWithRequestHeader(t *testing.T) {
	opts := &RequestOptions{}

	WithRequestHeader("X-Custom", "value1")(opts)
	WithRequestHeader("X-Another", "value2")(opts)

	if opts.Headers == nil {
		t.Fatal("Headers not initialized")
	}
	if opts.Headers["X-Custom"] != "value1" {
		t.Errorf("Header X-Custom = %s, want value1", opts.Headers["X-Custom"])
	}
	if opts.Headers["X-Another"] != "value2" {
		t.Errorf("Header X-Another = %s, want value2", opts.Headers["X-Another"])
	}
}

func TestWithRequestQuery(t *testing.T) {
	opts := &RequestOptions{}

	WithRequestQuery("limit", "10")(opts)
	WithRequestQuery("offset", "20")(opts)

	if opts.QueryParams == nil {
		t.Fatal("QueryParams not initialized")
	}
	if opts.QueryParams["limit"] != "10" {
		t.Errorf("Query limit = %s, want 10", opts.QueryParams["limit"])
	}
	if opts.QueryParams["offset"] != "20" {
		t.Errorf("Query offset = %s, want 20", opts.QueryParams["offset"])
	}
}

func TestWithIdempotencyKey(t *testing.T) {
	opts := &RequestOptions{}
	key := "test-idempotency-key"

	WithIdempotencyKey(key)(opts)

	if opts.IdempotencyKey != key {
		t.Errorf("IdempotencyKey = %s, want %s", opts.IdempotencyKey, key)
	}
}

func TestRequestOptions_Multiple(t *testing.T) {
	opts := &RequestOptions{}

	// Apply multiple options
	timeout := 5 * time.Second
	maxRetries := 3

	WithRequestTimeout(timeout)(opts)
	WithRequestMaxRetries(maxRetries)(opts)
	WithRequestHeader("X-Test", "value")(opts)
	WithRequestQuery("param", "value")(opts)
	WithIdempotencyKey("key-123")(opts)

	// Verify all options were applied
	if opts.Timeout == nil || *opts.Timeout != timeout {
		t.Error("Timeout not set correctly")
	}
	if opts.MaxRetries == nil || *opts.MaxRetries != maxRetries {
		t.Error("MaxRetries not set correctly")
	}
	if opts.Headers["X-Test"] != "value" {
		t.Error("Header not set correctly")
	}
	if opts.QueryParams["param"] != "value" {
		t.Error("Query param not set correctly")
	}
	if opts.IdempotencyKey != "key-123" {
		t.Error("IdempotencyKey not set correctly")
	}
}

func TestRequestOptions_HeaderOverwrite(t *testing.T) {
	opts := &RequestOptions{}

	WithRequestHeader("X-Test", "value1")(opts)
	WithRequestHeader("X-Test", "value2")(opts)

	if opts.Headers["X-Test"] != "value2" {
		t.Errorf("Header X-Test = %s, want value2 (should be overwritten)", opts.Headers["X-Test"])
	}
}

func TestRequestOptions_QueryOverwrite(t *testing.T) {
	opts := &RequestOptions{}

	WithRequestQuery("param", "value1")(opts)
	WithRequestQuery("param", "value2")(opts)

	if opts.QueryParams["param"] != "value2" {
		t.Errorf("Query param = %s, want value2 (should be overwritten)", opts.QueryParams["param"])
	}
}
