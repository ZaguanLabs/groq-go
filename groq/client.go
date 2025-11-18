package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/ZaguanLabs/groq-go/groq/audio"
	"github.com/ZaguanLabs/groq-go/groq/batches"
	"github.com/ZaguanLabs/groq-go/groq/chat"
	"github.com/ZaguanLabs/groq-go/groq/embeddings"
	"github.com/ZaguanLabs/groq-go/groq/files"
	"github.com/ZaguanLabs/groq-go/groq/internal/form"
	"github.com/ZaguanLabs/groq-go/groq/internal/querystring"
	"github.com/ZaguanLabs/groq-go/groq/internal/retry"
	"github.com/ZaguanLabs/groq-go/groq/models"
	"github.com/ZaguanLabs/groq-go/groq/option"
)

// Client is the Groq API client
type Client struct {
	httpClient *http.Client
	config     *ClientConfig

	// Resources
	Chat       *chat.Completions
	Embeddings *embeddings.Embeddings
	Audio      *audio.Audio
	Models     *models.Models
	Batches    *batches.Batches
	Files      *files.Files
}

// NewClient creates a new Groq API client
func NewClient(opts ...ClientOption) (*Client, error) {
	cfg := &ClientConfig{
		APIKey:         os.Getenv("GROQ_API_KEY"),
		BaseURL:        getEnvOrDefault("GROQ_BASE_URL", DefaultBaseURL),
		MaxRetries:     DefaultMaxRetries,
		Timeout:        DefaultTimeout,
		ConnectTimeout: DefaultConnectTimeout,
		Headers:        make(map[string]string),
		QueryParams:    make(map[string]string),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.APIKey == "" {
		return nil, &GroqError{Message: "API key required: set GROQ_API_KEY or use WithAPIKey()"}
	}

	if cfg.Logger == nil {
		cfg.Logger = defaultLogger
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     DefaultIdleConnTimeout,
				DialContext: (&net.Dialer{
					Timeout: cfg.ConnectTimeout,
				}).DialContext,
			},
		}
	}

	c := &Client{
		httpClient: cfg.HTTPClient,
		config:     cfg,
	}

	// Initialize resources
	c.Chat = chat.NewCompletions(c)
	c.Embeddings = embeddings.New(c)
	c.Audio = audio.New(c)
	c.Models = models.New(c)
	c.Batches = batches.New(c)
	c.Files = files.New(c)

	return c, nil
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Post sends a POST request
func (c *Client) Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
	reqOpts := &option.RequestOptions{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
	for _, opt := range opts {
		opt(reqOpts)
	}

	// Build request
	req, err := c.buildRequest(ctx, http.MethodPost, path, body, reqOpts)
	if err != nil {
		return err
	}

	// Execute with retry
	resp, err := c.doWithRetry(ctx, req, reqOpts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle errors
	if resp.StatusCode >= 400 {
		return c.handleError(resp)
	}

	// Validate Content-Type if strict validation is enabled
	if c.config.StrictValidation {
		ct := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			return &APIError{
				GroqError: GroqError{
					Message: fmt.Sprintf("Expected Content-Type application/json, got %s", ct),
					Request: req,
				},
				Response:   resp,
				StatusCode: resp.StatusCode,
			}
		}
	}

	// Parse response
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return &GroqError{Message: fmt.Sprintf("error decoding response: %v", err)}
		}
	}

	return nil
}

// PostStream sends a POST request and returns the raw response for streaming
func (c *Client) PostStream(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error) {
	reqOpts := &option.RequestOptions{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
	for _, opt := range opts {
		opt(reqOpts)
	}

	// Build request
	req, err := c.buildRequest(ctx, http.MethodPost, path, body, reqOpts)
	if err != nil {
		return nil, err
	}

	// Execute with retry
	resp, err := c.doWithRetry(ctx, req, reqOpts)
	if err != nil {
		return nil, err
	}

	// Handle errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, c.handleError(resp)
	}

	// Validate Content-Type if strict validation is enabled
	if c.config.StrictValidation {
		ct := resp.Header.Get("Content-Type")
		// Streaming usually returns text/event-stream
		if !strings.HasPrefix(ct, "text/event-stream") && !strings.HasPrefix(ct, "application/json") {
			defer resp.Body.Close()
			return nil, &APIError{
				GroqError: GroqError{
					Message: fmt.Sprintf("Expected Content-Type text/event-stream or application/json, got %s", ct),
					Request: req,
				},
				Response:   resp,
				StatusCode: resp.StatusCode,
			}
		}
	}

	return resp, nil
}

// Get sends a GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
	reqOpts := &option.RequestOptions{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
	for _, opt := range opts {
		opt(reqOpts)
	}

	req, err := c.buildRequest(ctx, http.MethodGet, path, nil, reqOpts)
	if err != nil {
		return err
	}

	return c.execute(ctx, req, result, reqOpts)
}

// GetStream sends a GET request and returns the raw response
func (c *Client) GetStream(ctx context.Context, path string, opts ...option.RequestOption) (*http.Response, error) {
	reqOpts := &option.RequestOptions{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
	for _, opt := range opts {
		opt(reqOpts)
	}

	req, err := c.buildRequest(ctx, http.MethodGet, path, nil, reqOpts)
	if err != nil {
		return nil, err
	}

	resp, err := c.doWithRetry(ctx, req, reqOpts)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, c.handleError(resp)
	}

	return resp, nil
}

// Delete sends a DELETE request
func (c *Client) Delete(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
	reqOpts := &option.RequestOptions{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
	for _, opt := range opts {
		opt(reqOpts)
	}

	req, err := c.buildRequest(ctx, http.MethodDelete, path, nil, reqOpts)
	if err != nil {
		return err
	}

	return c.execute(ctx, req, result, reqOpts)
}

// execute handles the common request execution logic
func (c *Client) execute(ctx context.Context, req *http.Request, result interface{}, opts *option.RequestOptions) error {
	resp, err := c.doWithRetry(ctx, req, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return c.handleError(resp)
	}

	if c.config.StrictValidation {
		ct := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			// Some endpoints might return empty body with 204 No Content
			if resp.StatusCode != http.StatusNoContent {
				return &APIError{
					GroqError: GroqError{
						Message: fmt.Sprintf("Expected Content-Type application/json, got %s", ct),
						Request: req,
					},
					Response:   resp,
					StatusCode: resp.StatusCode,
				}
			}
		}
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return &GroqError{Message: fmt.Sprintf("error decoding response: %v", err)}
		}
	}

	return nil
}

// PostForm sends a POST request with multipart/form-data
func (c *Client) PostForm(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error {
	reqOpts := &option.RequestOptions{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
	for _, opt := range opts {
		opt(reqOpts)
	}

	// Encode form
	enc := form.NewEncoder()
	contentType, bodyReader, err := enc.Encode(formStruct)
	if err != nil {
		return fmt.Errorf("form encode: %w", err)
	}

	// Set Content-Type header
	reqOpts.Headers["Content-Type"] = contentType

	// Build request
	url := c.buildURL(path, reqOpts)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}

	// Set other headers
	c.setHeaders(req, reqOpts)

	return c.execute(ctx, req, result, reqOpts)
}

func (c *Client) buildRequest(ctx context.Context, method, path string, body interface{}, opts *option.RequestOptions) (*http.Request, error) {
	url := c.buildURL(path, opts)

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req, opts)
	return req, nil
}

func (c *Client) buildURL(path string, opts *option.RequestOptions) string {
	baseURL := strings.TrimSuffix(c.config.BaseURL, "/")
	path = strings.TrimPrefix(path, "/")
	u := fmt.Sprintf("%s/%s", baseURL, path)

	// Merge query params
	query := make(map[string]interface{})
	for k, v := range c.config.QueryParams {
		query[k] = v
	}
	if opts != nil {
		for k, v := range opts.QueryParams {
			query[k] = v
		}
	}

	if len(query) > 0 {
		if qs, err := querystring.Stringify(query); err == nil && qs != "" {
			u += "?" + qs
		} else if err != nil && c.config.Logger != nil {
			c.config.Logger.Warn("Failed to serialize query parameters: %v", err)
		}
	}

	return u
}

func (c *Client) setHeaders(req *http.Request, opts *option.RequestOptions) {
	// Default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("Groq/Go %s", Version))
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	// Platform headers
	req.Header.Set("X-Stainless-Lang", "go")
	req.Header.Set("X-Stainless-Package-Version", Version)
	req.Header.Set("X-Stainless-OS", runtime.GOOS)
	req.Header.Set("X-Stainless-Arch", runtime.GOARCH)
	req.Header.Set("X-Stainless-Runtime", runtime.Version())

	// Custom headers from client config
	for k, v := range c.config.Headers {
		req.Header.Set(k, v)
	}

	// Request specific headers
	if opts != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
		if opts.IdempotencyKey != "" {
			req.Header.Set("Idempotency-Key", opts.IdempotencyKey)
		}
	}
}

func (c *Client) doWithRetry(ctx context.Context, req *http.Request, opts *option.RequestOptions) (*http.Response, error) {
	maxRetries := c.config.MaxRetries
	if opts.MaxRetries != nil {
		maxRetries = *opts.MaxRetries
	}

	retryCfg := retry.Config{
		MaxRetries:  maxRetries,
		ShouldRetry: retry.DefaultShouldRetry,
	}

	return retry.Do(ctx, retryCfg, func() (*http.Response, error) {
		// Clone request for each attempt (rewind body if needed)
		// For now assuming body is bytes.Reader which is seekable or we need to handle it
		// http.NewRequestWithContext wraps bytes.Reader properly.
		// But if we retry, we need to ensure the body is read from start.
		// `http.Client.Do` handles rewinding `GetBody` if set.
		// `http.NewRequest` sets `GetBody` for `bytes.Reader`.

		return c.httpClient.Do(req)
	})
}

func (c *Client) handleError(resp *http.Response) error {
	// Parse error body
	var body interface{}
	respBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &body)

	apiErr := APIError{
		GroqError: GroqError{
			Message: fmt.Sprintf("Error code: %d - %s", resp.StatusCode, string(respBytes)),
			Request: resp.Request,
		},
		Response:   resp,
		StatusCode: resp.StatusCode,
		Body:       body,
	}

	switch resp.StatusCode {
	case 400:
		return &BadRequestError{APIError: apiErr}
	case 401:
		return &AuthenticationError{APIError: apiErr}
	case 403:
		return &PermissionDeniedError{APIError: apiErr}
	case 404:
		return &NotFoundError{APIError: apiErr}
	case 409:
		return &ConflictError{APIError: apiErr}
	case 422:
		return &UnprocessableEntityError{APIError: apiErr}
	case 429:
		return &RateLimitError{APIError: apiErr}
	}

	if resp.StatusCode >= 500 {
		return &InternalServerError{APIError: apiErr}
	}

	return &apiErr
}
