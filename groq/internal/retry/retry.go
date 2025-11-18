package retry

import (
	"context"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	InitialDelay = 500 * time.Millisecond
	MaxDelay     = 8 * time.Second
)

type Config struct {
	MaxRetries  int
	ShouldRetry func(*http.Response) bool
}

func DefaultShouldRetry(resp *http.Response) bool {
	// Check x-should-retry header
	if retry := resp.Header.Get("x-should-retry"); retry != "" {
		return retry == "true"
	}

	// Retry on specific status codes
	switch resp.StatusCode {
	case 408, 409, 429: // Timeout, Conflict, Rate Limit
		return true
	}

	return resp.StatusCode >= 500
}

func CalculateBackoff(attempt int, resp *http.Response) time.Duration {
	// Check Retry-After header
	if retryAfter := resp.Header.Get("retry-after"); retryAfter != "" {
		if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds > 0 && seconds <= 60 {
			return time.Duration(seconds) * time.Second
		}
		// Could also parse HTTP-date format here
	}

	// Check retry-after-ms header (non-standard)
	if retryMs := resp.Header.Get("retry-after-ms"); retryMs != "" {
		if ms, err := strconv.Atoi(retryMs); err == nil {
			return time.Duration(ms) * time.Millisecond
		}
	}

	// Exponential backoff with jitter
	delay := math.Min(float64(InitialDelay)*math.Pow(2, float64(attempt)), float64(MaxDelay))
	jitter := 1 - 0.25*rand.Float64()
	return time.Duration(delay * jitter)
}

func Do(ctx context.Context, cfg Config, fn func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		resp, err = fn()

		if err != nil {
			// Network errors (err != nil) are typically retried if safe/idempotent?
			// For now we don't retry network errors automatically unless we inspect them
			// But usually SDKs retry on connection errors.
			// The plan doesn't specify handling network error retries explicitly in the code snippet,
			// but logic implies we return if error.
			// Let's stick to the plan's snippet which returns err if not nil.
			// Wait, the plan snippet says:
			// if err != nil || !cfg.ShouldRetry(resp) { return resp, err }
			// This means it does NOT retry on error? That seems wrong for a robust SDK.
			// However, checking the snippet in 3.4:
			// if err != nil || !cfg.ShouldRetry(resp) { return resp, err }
			// This implies only status code retries.
			// But `DefaultShouldRetry` takes `resp`.
			// If `err` is not nil, `resp` is nil.
			// So we return `nil, err`.
			// If we want to retry on timeout/connection error, we need to handle `err`.
			// I will stick to the plan's snippet for now.
			return resp, err
		}

		if !cfg.ShouldRetry(resp) {
			return resp, err
		}

		if attempt < cfg.MaxRetries {
			backoff := CalculateBackoff(attempt, resp)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return resp, err
}
