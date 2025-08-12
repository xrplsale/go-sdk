package xrplsale

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// Version is the SDK version
	Version = "1.0.0"
	
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second
	
	// DefaultMaxRetries is the default maximum retry attempts
	DefaultMaxRetries = 3
	
	// ProductionBaseURL is the production API base URL
	ProductionBaseURL = "https://api.xrpl.sale/v1"
	
	// TestnetBaseURL is the testnet API base URL
	TestnetBaseURL = "https://api-testnet.xrpl.sale/v1"
)

// Environment represents the API environment
type Environment string

const (
	Production Environment = "production"
	Testnet    Environment = "testnet"
)

// Config holds the client configuration
type Config struct {
	APIKey        string
	Environment   Environment
	BaseURL       string
	Timeout       time.Duration
	MaxRetries    int
	RetryWaitTime time.Duration
	WebhookSecret string
	Debug         bool
}

// Client is the main XRPL.Sale SDK client
type Client struct {
	config     *Config
	httpClient *resty.Client
	authToken  string
	
	// Services
	Auth        *AuthService
	Projects    *ProjectsService
	Investments *InvestmentsService
	Analytics   *AnalyticsService
	Webhooks    *WebhooksService
}

// NewClient creates a new XRPL.Sale client
func NewClient(apiKey string) *Client {
	return NewClientWithConfig(&Config{
		APIKey:      apiKey,
		Environment: Production,
	})
}

// NewClientWithConfig creates a new client with custom configuration
func NewClientWithConfig(config *Config) *Client {
	// Set defaults
	if config.Environment == "" {
		config.Environment = Production
	}
	
	if config.BaseURL == "" {
		if config.Environment == Testnet {
			config.BaseURL = TestnetBaseURL
		} else {
			config.BaseURL = ProductionBaseURL
		}
	}
	
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}
	
	if config.MaxRetries == 0 {
		config.MaxRetries = DefaultMaxRetries
	}
	
	if config.RetryWaitTime == 0 {
		config.RetryWaitTime = 1 * time.Second
	}
	
	// Create HTTP client
	httpClient := resty.New().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetRetryCount(config.MaxRetries).
		SetRetryWaitTime(config.RetryWaitTime).
		SetRetryMaxWaitTime(10 * time.Second).
		SetHeader("User-Agent", "XRPL.Sale-Go-SDK/"+Version).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")
	
	if config.Debug {
		httpClient.SetDebug(true)
	}
	
	// Add retry conditions
	httpClient.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return err != nil || r.StatusCode() >= 500
		},
	)
	
	client := &Client{
		config:     config,
		httpClient: httpClient,
	}
	
	// Initialize services
	client.Auth = &AuthService{client: client}
	client.Projects = &ProjectsService{client: client}
	client.Investments = &InvestmentsService{client: client}
	client.Analytics = &AnalyticsService{client: client}
	client.Webhooks = &WebhooksService{client: client}
	
	// Set API key header if provided
	if config.APIKey != "" {
		httpClient.SetHeader("X-API-Key", config.APIKey)
	}
	
	return client
}

// SetAuthToken sets the authentication token for requests
func (c *Client) SetAuthToken(token string) {
	c.authToken = token
	c.httpClient.SetAuthToken(token)
}

// Request makes an authenticated API request
func (c *Client) Request(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	req := c.httpClient.R().
		SetContext(ctx)
	
	if body != nil {
		req.SetBody(body)
	}
	
	if result != nil {
		req.SetResult(result)
	}
	
	// Set error structure
	apiError := &APIError{}
	req.SetError(apiError)
	
	var resp *resty.Response
	var err error
	
	switch method {
	case http.MethodGet:
		resp, err = req.Get(endpoint)
	case http.MethodPost:
		resp, err = req.Post(endpoint)
	case http.MethodPut:
		resp, err = req.Put(endpoint)
	case http.MethodPatch:
		resp, err = req.Patch(endpoint)
	case http.MethodDelete:
		resp, err = req.Delete(endpoint)
	default:
		return fmt.Errorf("unsupported method: %s", method)
	}
	
	if err != nil {
		return err
	}
	
	// Check for error response
	if resp.IsError() {
		if apiError.Message != "" {
			return apiError
		}
		return fmt.Errorf("API error: %d %s", resp.StatusCode(), resp.Status())
	}
	
	return nil
}

// Get makes a GET request
func (c *Client) Get(ctx context.Context, endpoint string, params map[string]string, result interface{}) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(result).
		SetError(&APIError{})
	
	resp, err := req.Get(endpoint)
	if err != nil {
		return err
	}
	
	if resp.IsError() {
		apiErr := resp.Error().(*APIError)
		if apiErr.Message != "" {
			return apiErr
		}
		return fmt.Errorf("API error: %d", resp.StatusCode())
	}
	
	return nil
}

// Post makes a POST request
func (c *Client) Post(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, endpoint, body, result)
}

// Put makes a PUT request
func (c *Client) Put(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPut, endpoint, body, result)
}

// Patch makes a PATCH request
func (c *Client) Patch(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPatch, endpoint, body, result)
}

// Delete makes a DELETE request
func (c *Client) Delete(ctx context.Context, endpoint string, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, endpoint, nil, result)
}

// VerifyWebhookSignature verifies a webhook signature
func (c *Client) VerifyWebhookSignature(payload []byte, signature string) bool {
	if c.config.WebhookSecret == "" {
		return false
	}
	
	mac := hmac.New(sha256.New, []byte(c.config.WebhookSecret))
	mac.Write(payload)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// ParseWebhookEvent parses a webhook event from JSON
func (c *Client) ParseWebhookEvent(payload []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
	}
	return &event, nil
}