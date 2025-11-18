package retry

import (
	"net/http"
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
