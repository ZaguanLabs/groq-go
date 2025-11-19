package groq

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ZaguanLabs/groq-go/groq/option"
)

func TestNewClient_Defaults(t *testing.T) {
	os.Setenv("GROQ_API_KEY", "test-key")
	defer os.Unsetenv("GROQ_API_KEY")

	c, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if c.config.BaseURL != DefaultBaseURL {
		t.Errorf("BaseURL = %s, want %s", c.config.BaseURL, DefaultBaseURL)
	}
	if c.config.APIKey != "test-key" {
		t.Errorf("APIKey = %s, want test-key", c.config.APIKey)
	}
}

func TestNewClient_MissingAPIKey(t *testing.T) {
	os.Unsetenv("GROQ_API_KEY")

	_, err := NewClient()
	if err == nil {
		t.Error("NewClient() expected error for missing API key, got nil")
	}
}

func TestNewClient_Options(t *testing.T) {
	c, err := NewClient(
		WithAPIKey("custom-key"),
		WithBaseURL("https://custom.api.com"),
		WithTimeout(10*time.Second),
		WithHeader("X-Custom", "value"),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if c.config.APIKey != "custom-key" {
		t.Errorf("APIKey = %s, want custom-key", c.config.APIKey)
	}
	if c.config.BaseURL != "https://custom.api.com" {
		t.Errorf("BaseURL = %s, want https://custom.api.com", c.config.BaseURL)
	}
	if c.config.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", c.config.Timeout)
	}
	if val := c.config.Headers["X-Custom"]; val != "value" {
		t.Errorf("Header X-Custom = %s, want value", val)
	}
}

func TestNewClient_Logger(t *testing.T) {
	os.Setenv("GROQ_API_KEY", "test-key")
	defer os.Unsetenv("GROQ_API_KEY")

	logger := &LeveledLogger{Level: LevelDebug}
	c, err := NewClient(WithLogger(logger))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if c.config.Logger != logger {
		t.Error("Logger not set correctly")
	}
}

func TestClient_PostStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data: hello\n\n"))
	}))
	defer server.Close()

	c, err := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	resp, err := c.PostStream(context.Background(), "/stream", nil)
	if err != nil {
		t.Fatalf("PostStream error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"test"}`))
	}))
	defer server.Close()

	c, err := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	var result map[string]string
	err = c.Post(context.Background(), "/test", map[string]string{"key": "value"}, &result)
	if err != nil {
		t.Fatalf("Post error: %v", err)
	}

	if result["id"] != "test" {
		t.Errorf("result id = %s, want test", result["id"])
	}
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %s, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":"test"}`))
	}))
	defer server.Close()

	c, err := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	var result map[string]string
	err = c.Get(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}

	if result["data"] != "test" {
		t.Errorf("result data = %s, want test", result["data"])
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %s, want DELETE", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"deleted":true}`))
	}))
	defer server.Close()

	c, err := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	var result map[string]bool
	err = c.Delete(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	if !result["deleted"] {
		t.Error("expected deleted = true")
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		wantErrType interface{}
	}{
		{"BadRequest", 400, &BadRequestError{}},
		{"Unauthorized", 401, &AuthenticationError{}},
		{"Forbidden", 403, &PermissionDeniedError{}},
		{"NotFound", 404, &NotFoundError{}},
		{"Conflict", 409, &ConflictError{}},
		{"UnprocessableEntity", 422, &UnprocessableEntityError{}},
		{"RateLimit", 429, &RateLimitError{}},
		{"InternalServerError", 500, &InternalServerError{}},
		{"BadGateway", 502, &InternalServerError{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(`{"error":"test error"}`))
			}))
			defer server.Close()

			c, _ := NewClient(
				WithAPIKey("test-key"),
				WithBaseURL(server.URL),
			)

			var result map[string]string
			err := c.Post(context.Background(), "/test", nil, &result)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			// Check error type
			switch tt.wantErrType.(type) {
			case *BadRequestError:
				if _, ok := err.(*BadRequestError); !ok {
					t.Errorf("error type = %T, want *BadRequestError", err)
				}
			case *AuthenticationError:
				if _, ok := err.(*AuthenticationError); !ok {
					t.Errorf("error type = %T, want *AuthenticationError", err)
				}
			case *PermissionDeniedError:
				if _, ok := err.(*PermissionDeniedError); !ok {
					t.Errorf("error type = %T, want *PermissionDeniedError", err)
				}
			case *NotFoundError:
				if _, ok := err.(*NotFoundError); !ok {
					t.Errorf("error type = %T, want *NotFoundError", err)
				}
			case *ConflictError:
				if _, ok := err.(*ConflictError); !ok {
					t.Errorf("error type = %T, want *ConflictError", err)
				}
			case *UnprocessableEntityError:
				if _, ok := err.(*UnprocessableEntityError); !ok {
					t.Errorf("error type = %T, want *UnprocessableEntityError", err)
				}
			case *RateLimitError:
				if _, ok := err.(*RateLimitError); !ok {
					t.Errorf("error type = %T, want *RateLimitError", err)
				}
			case *InternalServerError:
				if _, ok := err.(*InternalServerError); !ok {
					t.Errorf("error type = %T, want *InternalServerError", err)
				}
			}
		})
	}
}

func TestClient_Headers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("Authorization = %s, want Bearer test-key", auth)
		}
		if ua := r.Header.Get("User-Agent"); ua == "" {
			t.Error("User-Agent header missing")
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %s, want application/json", ct)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	var result map[string]string
	c.Post(context.Background(), "/test", map[string]string{}, &result)
}

func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result map[string]string
	err := c.Post(ctx, "/test", nil, &result)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestClient_StrictValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)
	c.config.StrictValidation = true

	var result map[string]string
	err := c.Post(context.Background(), "/test", nil, &result)
	if err == nil {
		t.Fatal("expected content-type validation error")
	}
}

func TestClient_ResourcesInitialized(t *testing.T) {
	c, _ := NewClient(WithAPIKey("test-key"))

	if c.Chat == nil {
		t.Error("Chat not initialized")
	}
	if c.Embeddings == nil {
		t.Error("Embeddings not initialized")
	}
	if c.Audio == nil {
		t.Error("Audio not initialized")
	}
	if c.Models == nil {
		t.Error("Models not initialized")
	}
	if c.Batches == nil {
		t.Error("Batches not initialized")
	}
	if c.Files == nil {
		t.Error("Files not initialized")
	}
}

func TestClient_CustomHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithHTTPClient(customClient),
	)

	if c.httpClient != customClient {
		t.Error("custom HTTP client not set")
	}
}

func TestClient_PostForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "multipart/form-data") {
			t.Errorf("Content-Type = %s, want multipart/form-data", ct)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"file-123"}`))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	type FormData struct {
		File    io.Reader `json:"file"`
		Purpose string    `json:"purpose"`
	}

	var result map[string]string
	err := c.PostForm(context.Background(), "/test", &FormData{
		File:    strings.NewReader("test content"),
		Purpose: "batch",
	}, &result)

	if err != nil {
		t.Fatalf("PostForm error: %v", err)
	}

	if result["id"] != "file-123" {
		t.Errorf("result id = %s, want file-123", result["id"])
	}
}

func TestClient_GetStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %s, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("binary content"))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := c.GetStream(context.Background(), "/test")
	if err != nil {
		t.Fatalf("GetStream error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	content, _ := io.ReadAll(resp.Body)
	if string(content) != "binary content" {
		t.Errorf("content = %s, want binary content", string(content))
	}
}

func TestClient_WithMaxRetries(t *testing.T) {
	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithMaxRetries(5),
	)

	if c.config.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", c.config.MaxRetries)
	}
}

func TestClient_DefaultHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if custom := r.Header.Get("X-Custom-Header"); custom != "custom-value" {
			t.Errorf("X-Custom-Header = %s, want custom-value", custom)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithHeader("X-Custom-Header", "custom-value"),
	)

	var result map[string]string
	c.Get(context.Background(), "/test", &result)
}

func TestClient_RetryWithBackoff(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limited"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithMaxRetries(3),
	)

	var result map[string]interface{}
	err := c.Post(context.Background(), "/test", nil, &result)
	if err != nil {
		t.Fatalf("Post error: %v", err)
	}

	if attempts != 2 {
		t.Errorf("attempts = %d, want 2", attempts)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		t.Error("expected success = true")
	}
}

func TestClient_PostStreamStrictValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not sse"))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)
	c.config.StrictValidation = true

	_, err := c.PostStream(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("expected content-type validation error")
	}
}

func TestClient_GetStreamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	_, err := c.GetStream(context.Background(), "/test")
	if err == nil {
		t.Fatal("expected error for 404 status")
	}
}

func TestClient_PostFormError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid form data"}`))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	type FormData struct {
		File io.Reader `json:"file"`
	}

	var result map[string]string
	err := c.PostForm(context.Background(), "/test", &FormData{File: strings.NewReader("test")}, &result)
	if err == nil {
		t.Fatal("expected error for 400 status")
	}

	if _, ok := err.(*BadRequestError); !ok {
		t.Errorf("expected BadRequestError, got %T", err)
	}
}

func TestClient_BuildURLWithQueryParams(t *testing.T) {
	c, _ := NewClient(WithAPIKey("test-key"))

	reqOpts := &option.RequestOptions{
		QueryParams: map[string]string{
			"limit":  "10",
			"offset": "20",
		},
	}

	url := c.buildURL("/test", reqOpts)
	if !strings.Contains(url, "limit=10") {
		t.Error("URL missing limit parameter")
	}
	if !strings.Contains(url, "offset=20") {
		t.Error("URL missing offset parameter")
	}
}

func TestClient_ExecuteJSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	var result map[string]string
	err := c.Post(context.Background(), "/test", nil, &result)
	if err == nil {
		t.Fatal("expected JSON unmarshal error")
	}
}

func TestClient_BaseURLFromEnv(t *testing.T) {
	os.Setenv("GROQ_BASE_URL", "https://custom.api.com")
	os.Setenv("GROQ_API_KEY", "test-key")
	defer os.Unsetenv("GROQ_BASE_URL")
	defer os.Unsetenv("GROQ_API_KEY")

	c, _ := NewClient()
	if c.config.BaseURL != "https://custom.api.com" {
		t.Errorf("BaseURL = %s, want https://custom.api.com", c.config.BaseURL)
	}
}
