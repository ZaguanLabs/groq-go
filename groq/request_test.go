package groq

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/option"
)

func TestClient_buildURL(t *testing.T) {
	c, _ := NewClient(WithAPIKey("test"), WithBaseURL("https://api.example.com/v1"))

	// Test merging config and request query params
	c.config.QueryParams = map[string]string{
		"global": "val",
	}

	opts := &option.RequestOptions{
		QueryParams: map[string]string{
			"req": "val2",
		},
	}

	got := c.buildURL("chat/completions", opts)

	// Order of params is not guaranteed by map iteration unless sorted in encode,
	// but our stringify sorts.
	if !strings.Contains(got, "https://api.example.com/v1/chat/completions") {
		t.Errorf("URL base mismatch: %s", got)
	}
	if !strings.Contains(got, "global=val") {
		t.Error("Missing global param")
	}
	if !strings.Contains(got, "req=val2") {
		t.Error("Missing request param")
	}
}

func TestClient_setHeaders(t *testing.T) {
	c, _ := NewClient(
		WithAPIKey("test-key"),
		WithHeader("X-Global", "global-val"),
	)

	req, _ := http.NewRequest("GET", "https://example.com", nil)

	opts := &option.RequestOptions{
		Headers: map[string]string{
			"X-Request": "request-val",
		},
	}

	c.setHeaders(req, opts)

	if req.Header.Get("Authorization") != "Bearer test-key" {
		t.Error("Missing Authorization header")
	}
	if req.Header.Get("X-Global") != "global-val" {
		t.Error("Missing global header")
	}
	if req.Header.Get("X-Request") != "request-val" {
		t.Error("Missing request header")
	}
	if req.Header.Get("X-Stainless-Lang") != "go" {
		t.Error("Missing platform header")
	}
}
