package xrplsale

import (
	"context"
	"fmt"
	"time"
)

// ProjectsService handles project-related operations
type ProjectsService struct {
	client *Client
}

// ListProjectsOptions represents options for listing projects
type ListProjectsOptions struct {
	Status   string `url:"status,omitempty"`
	Page     int    `url:"page,omitempty"`
	Limit    int    `url:"limit,omitempty"`
	SortBy   string `url:"sort_by,omitempty"`
	SortOrder string `url:"sort_order,omitempty"`
}

// List retrieves a list of projects
func (ps *ProjectsService) List(ctx context.Context, opts *ListProjectsOptions) (*PaginatedResponse[Project], error) {
	params := make(map[string]string)
	if opts != nil {
		if opts.Status != "" {
			params["status"] = opts.Status
		}
		if opts.Page > 0 {
			params["page"] = fmt.Sprintf("%d", opts.Page)
		}
		if opts.Limit > 0 {
			params["limit"] = fmt.Sprintf("%d", opts.Limit)
		}
		if opts.SortBy != "" {
			params["sort_by"] = opts.SortBy
		}
		if opts.SortOrder != "" {
			params["sort_order"] = opts.SortOrder
		}
	}
	
	var result PaginatedResponse[Project]
	err := ps.client.Get(ctx, "/projects", params, &result)
	return &result, err
}

// GetActive retrieves active projects
func (ps *ProjectsService) GetActive(ctx context.Context, page, limit int) (*PaginatedResponse[Project], error) {
	return ps.List(ctx, &ListProjectsOptions{
		Status: "active",
		Page:   page,
		Limit:  limit,
	})
}

// Get retrieves a specific project
func (ps *ProjectsService) Get(ctx context.Context, projectID string) (*Project, error) {
	var project Project
	err := ps.client.Get(ctx, fmt.Sprintf("/projects/%s", projectID), nil, &project)
	return &project, err
}

// Create creates a new project
func (ps *ProjectsService) Create(ctx context.Context, project *CreateProjectRequest) (*Project, error) {
	var result Project
	err := ps.client.Post(ctx, "/projects", project, &result)
	return &result, err
}

// Update updates a project
func (ps *ProjectsService) Update(ctx context.Context, projectID string, updates map[string]interface{}) (*Project, error) {
	var result Project
	err := ps.client.Patch(ctx, fmt.Sprintf("/projects/%s", projectID), updates, &result)
	return &result, err
}

// Launch launches a project
func (ps *ProjectsService) Launch(ctx context.Context, projectID string) (*Project, error) {
	var result Project
	err := ps.client.Post(ctx, fmt.Sprintf("/projects/%s/launch", projectID), nil, &result)
	return &result, err
}

// GetStats retrieves project statistics
func (ps *ProjectsService) GetStats(ctx context.Context, projectID string) (*ProjectStats, error) {
	var stats ProjectStats
	err := ps.client.Get(ctx, fmt.Sprintf("/projects/%s/stats", projectID), nil, &stats)
	return &stats, err
}

// InvestmentsService handles investment-related operations
type InvestmentsService struct {
	client *Client
}

// Create creates a new investment
func (is *InvestmentsService) Create(ctx context.Context, investment *CreateInvestmentRequest) (*Investment, error) {
	var result Investment
	err := is.client.Post(ctx, "/investments", investment, &result)
	return &result, err
}

// Get retrieves a specific investment
func (is *InvestmentsService) Get(ctx context.Context, investmentID string) (*Investment, error) {
	var investment Investment
	err := is.client.Get(ctx, fmt.Sprintf("/investments/%s", investmentID), nil, &investment)
	return &investment, err
}

// GetByProject retrieves investments for a project
func (is *InvestmentsService) GetByProject(ctx context.Context, projectID string, page, limit int) (*PaginatedResponse[Investment], error) {
	params := map[string]string{
		"page":  fmt.Sprintf("%d", page),
		"limit": fmt.Sprintf("%d", limit),
	}
	
	var result PaginatedResponse[Investment]
	err := is.client.Get(ctx, fmt.Sprintf("/projects/%s/investments", projectID), params, &result)
	return &result, err
}

// GetInvestorSummary retrieves an investor's summary
func (is *InvestmentsService) GetInvestorSummary(ctx context.Context, investorAccount string) (*InvestorSummary, error) {
	var summary InvestorSummary
	err := is.client.Get(ctx, fmt.Sprintf("/investors/%s/summary", investorAccount), nil, &summary)
	return &summary, err
}

// Simulate simulates an investment
func (is *InvestmentsService) Simulate(ctx context.Context, simulation *SimulateInvestmentRequest) (*SimulationResult, error) {
	var result SimulationResult
	err := is.client.Post(ctx, "/investments/simulate", simulation, &result)
	return &result, err
}

// AnalyticsService handles analytics operations
type AnalyticsService struct {
	client *Client
}

// GetPlatformAnalytics retrieves platform-wide analytics
func (as *AnalyticsService) GetPlatformAnalytics(ctx context.Context) (*PlatformAnalytics, error) {
	var analytics PlatformAnalytics
	err := as.client.Get(ctx, "/analytics/platform", nil, &analytics)
	return &analytics, err
}

// GetProjectAnalytics retrieves project-specific analytics
func (as *AnalyticsService) GetProjectAnalytics(ctx context.Context, projectID string, startDate, endDate time.Time) (*ProjectAnalytics, error) {
	params := map[string]string{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	}
	
	var analytics ProjectAnalytics
	err := as.client.Get(ctx, fmt.Sprintf("/analytics/projects/%s", projectID), params, &analytics)
	return &analytics, err
}

// GetMarketTrends retrieves market trends
func (as *AnalyticsService) GetMarketTrends(ctx context.Context, period string) (*MarketTrends, error) {
	params := map[string]string{"period": period}
	var trends MarketTrends
	err := as.client.Get(ctx, "/analytics/trends", params, &trends)
	return &trends, err
}

// ExportData exports analytics data
func (as *AnalyticsService) ExportData(ctx context.Context, exportReq *ExportDataRequest) (*ExportResult, error) {
	var result ExportResult
	err := as.client.Post(ctx, "/analytics/export", exportReq, &result)
	return &result, err
}

// AuthService handles authentication operations
type AuthService struct {
	client *Client
}

// GenerateChallenge generates an authentication challenge
func (as *AuthService) GenerateChallenge(ctx context.Context, walletAddress string) (*AuthChallenge, error) {
	req := map[string]string{"wallet_address": walletAddress}
	var challenge AuthChallenge
	err := as.client.Post(ctx, "/auth/challenge", req, &challenge)
	return &challenge, err
}

// Authenticate authenticates with wallet signature
func (as *AuthService) Authenticate(ctx context.Context, authReq *AuthRequest) (*AuthResponse, error) {
	var response AuthResponse
	err := as.client.Post(ctx, "/auth/wallet", authReq, &response)
	if err == nil && response.Token != "" {
		as.client.SetAuthToken(response.Token)
	}
	return &response, err
}

// Refresh refreshes the authentication token
func (as *AuthService) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	req := map[string]string{"refresh_token": refreshToken}
	var response AuthResponse
	err := as.client.Post(ctx, "/auth/refresh", req, &response)
	if err == nil && response.Token != "" {
		as.client.SetAuthToken(response.Token)
	}
	return &response, err
}

// Logout logs out the current session
func (as *AuthService) Logout(ctx context.Context) error {
	return as.client.Post(ctx, "/auth/logout", nil, nil)
}

// GetProfile retrieves the current user profile
func (as *AuthService) GetProfile(ctx context.Context) (*UserProfile, error) {
	var profile UserProfile
	err := as.client.Get(ctx, "/auth/profile", nil, &profile)
	return &profile, err
}

// WebhooksService handles webhook operations
type WebhooksService struct {
	client *Client
}

// Register registers a new webhook
func (ws *WebhooksService) Register(ctx context.Context, webhook *RegisterWebhookRequest) (*Webhook, error) {
	var result Webhook
	err := ws.client.Post(ctx, "/webhooks", webhook, &result)
	return &result, err
}

// List retrieves all webhooks
func (ws *WebhooksService) List(ctx context.Context) ([]*Webhook, error) {
	var webhooks []*Webhook
	err := ws.client.Get(ctx, "/webhooks", nil, &webhooks)
	return webhooks, err
}

// Get retrieves a specific webhook
func (ws *WebhooksService) Get(ctx context.Context, webhookID string) (*Webhook, error) {
	var webhook Webhook
	err := ws.client.Get(ctx, fmt.Sprintf("/webhooks/%s", webhookID), nil, &webhook)
	return &webhook, err
}

// Update updates a webhook
func (ws *WebhooksService) Update(ctx context.Context, webhookID string, updates map[string]interface{}) (*Webhook, error) {
	var webhook Webhook
	err := ws.client.Patch(ctx, fmt.Sprintf("/webhooks/%s", webhookID), updates, &webhook)
	return &webhook, err
}

// Delete deletes a webhook
func (ws *WebhooksService) Delete(ctx context.Context, webhookID string) error {
	return ws.client.Delete(ctx, fmt.Sprintf("/webhooks/%s", webhookID), nil)
}

// Test tests a webhook delivery
func (ws *WebhooksService) Test(ctx context.Context, webhookID string) error {
	return ws.client.Post(ctx, fmt.Sprintf("/webhooks/%s/test", webhookID), nil, nil)
}

// GetDeliveries retrieves webhook delivery logs
func (ws *WebhooksService) GetDeliveries(ctx context.Context, webhookID string, page, limit int) (*PaginatedResponse[WebhookDelivery], error) {
	params := map[string]string{
		"page":  fmt.Sprintf("%d", page),
		"limit": fmt.Sprintf("%d", limit),
	}
	
	var result PaginatedResponse[WebhookDelivery]
	err := ws.client.Get(ctx, fmt.Sprintf("/webhooks/%s/deliveries", webhookID), params, &result)
	return &result, err
}