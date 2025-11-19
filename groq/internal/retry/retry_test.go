package retry

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestCalculateBackoff(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		attempt  int
		minDelay time.Duration
		maxDelay time.Duration
	}{
		{
			name:     "exponential backoff attempt 0",
			attempt:  0,
			minDelay: 375 * time.Millisecond, // 500ms * 0.75
			maxDelay: 500 * time.Millisecond,
		},
		{
			name:     "exponential backoff attempt 1",
			attempt:  1,
			minDelay: 750 * time.Millisecond, // 1000ms * 0.75
			maxDelay: 1000 * time.Millisecond,
		},
		{
			name:     "retry-after seconds",
			headers:  map[string]string{"Retry-After": "10"},
			minDelay: 10 * time.Second,
			maxDelay: 10 * time.Second,
		},
		{
			name:     "retry-after http date",
			headers:  map[string]string{"Retry-After": time.Now().UTC().Add(5 * time.Second).Format(http.TimeFormat)},
			minDelay: 4 * time.Second, // Allow some buffer for execution time
			maxDelay: 6 * time.Second,
		},
		{
			name:     "retry-after-ms",
			headers:  map[string]string{"retry-after-ms": "1500"},
			minDelay: 1500 * time.Millisecond,
			maxDelay: 1500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{Header: make(http.Header)}
			for k, v := range tt.headers {
				resp.Header.Set(k, v)
			}

			delay := CalculateBackoff(tt.attempt, resp)

			if delay < tt.minDelay || delay > tt.maxDelay {
				t.Errorf("CalculateBackoff() = %v, want between %v and %v", delay, tt.minDelay, tt.maxDelay)
			}
		})
	}
}

func TestDefaultShouldRetry(t *testing.T) {
	tests := []struct {
		status int
		retry  string
		want   bool
	}{
		{200, "", false},
		{400, "", false},
		{401, "", false},
		{403, "", false},
		{404, "", false},
		{408, "", true},
		{409, "", true},
		{429, "", true},
		{500, "", true},
		{502, "", true},
		{503, "", true},
		{400, "true", true}, // x-should-retry override
	}

	for _, tt := range tests {
		resp := &http.Response{
			StatusCode: tt.status,
			Header:     make(http.Header),
		}
		if tt.retry != "" {
			resp.Header.Set("x-should-retry", tt.retry)
		}

		if got := DefaultShouldRetry(resp); got != tt.want {
			t.Errorf("DefaultShouldRetry(%d, %s) = %v, want %v", tt.status, tt.retry, got, tt.want)
		}
	}
}

func TestDo_Success(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	cfg := Config{
		MaxRetries:  3,
		ShouldRetry: DefaultShouldRetry,
	}

	resp, err := Do(ctx, cfg, func() (*http.Response, error) {
		attempts++
		return &http.Response{StatusCode: 200}, nil
	})

	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1", attempts)
	}
}

func TestDo_RetryOnce(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	cfg := Config{
		MaxRetries:  3,
		ShouldRetry: DefaultShouldRetry,
	}

	resp, err := Do(ctx, cfg, func() (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return &http.Response{StatusCode: 500, Header: make(http.Header)}, nil
		}
		return &http.Response{StatusCode: 200}, nil
	})

	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
	if attempts != 2 {
		t.Errorf("attempts = %d, want 2", attempts)
	}
}

func TestDo_MaxRetriesExceeded(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	cfg := Config{
		MaxRetries:  2,
		ShouldRetry: DefaultShouldRetry,
	}

	resp, err := Do(ctx, cfg, func() (*http.Response, error) {
		attempts++
		return &http.Response{StatusCode: 500, Header: make(http.Header)}, nil
	})

	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if resp.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", resp.StatusCode)
	}
	// MaxRetries=2 means 3 total attempts (initial + 2 retries)
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
}

func TestDo_NetworkError(t *testing.T) {
	ctx := context.Background()

	cfg := Config{
		MaxRetries:  3,
		ShouldRetry: DefaultShouldRetry,
	}

	_, err := Do(ctx, cfg, func() (*http.Response, error) {
		return nil, &url.Error{Op: "Get", URL: "http://test", Err: errors.New("connection refused")}
	})

	if err == nil {
		t.Fatal("Do() expected error, got nil")
	}
}

func TestDo_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	cfg := Config{
		MaxRetries:  3,
		ShouldRetry: DefaultShouldRetry,
	}

	_, err := Do(ctx, cfg, func() (*http.Response, error) {
		return &http.Response{StatusCode: 500, Header: make(http.Header)}, nil
	})

	if err == nil {
		t.Fatal("Do() expected context error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestCalculateBackoff_MaxDelay(t *testing.T) {
	// Test that backoff doesn't exceed max delay
	resp := &http.Response{Header: make(http.Header)}

	// Attempt 10 should be capped at MaxDelay
	delay := CalculateBackoff(10, resp)
	if delay > MaxDelay {
		t.Errorf("CalculateBackoff(10) = %v, want <= %v", delay, MaxDelay)
	}
}

func TestCalculateBackoff_InvalidRetryAfter(t *testing.T) {
	resp := &http.Response{Header: make(http.Header)}
	resp.Header.Set("Retry-After", "invalid")

	// Should fall back to exponential backoff
	delay := CalculateBackoff(0, resp)
	if delay < 375*time.Millisecond || delay > 500*time.Millisecond {
		t.Errorf("CalculateBackoff with invalid Retry-After = %v, want exponential backoff", delay)
	}
}
