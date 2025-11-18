package groq

import (
	"os"
	"testing"
	"time"
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
