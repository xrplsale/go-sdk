# XRPL.Sale Go SDK

Official Go SDK for integrating with the XRPL.Sale platform - the native XRPL launchpad for token sales and project funding.

[![Go Reference](https://pkg.go.dev/badge/github.com/xrplsale/go-sdk.svg)](https://pkg.go.dev/github.com/xrplsale/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/xrplsale/go-sdk)](https://goreportcard.com/report/github.com/xrplsale/go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- 🚀 **Idiomatic Go** - Built following Go best practices and conventions
- 🔄 **Context Support** - Full context.Context support for cancellation and timeouts
- 🔐 **XRPL Wallet Authentication** - Seamless wallet integration
- 📊 **Project Management** - Create, launch, and manage token sales
- 💰 **Investment Tracking** - Monitor investments and analytics
- 🔔 **Webhook Support** - Real-time event notifications with signature verification
- 📈 **Analytics & Reporting** - Comprehensive data insights
- 🛡️ **Error Handling** - Structured error types with detailed information
- 🔄 **Auto-retry Logic** - Resilient API calls with exponential backoff
- ⚡ **Concurrent Safe** - Thread-safe operations for concurrent usage

## Installation

```bash
go get github.com/xrplsale/go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/xrplsale/go-sdk"
)

func main() {
    // Initialize the client
    client := xrplsale.NewClient("your-api-key")
    
    // Or with custom configuration
    client = xrplsale.NewClientWithConfig(&xrplsale.Config{
        APIKey:      "your-api-key",
        Environment: xrplsale.Production, // or xrplsale.Testnet
        Debug:       true,
    })
    
    ctx := context.Background()
    
    // Create a new project
    project, err := client.Projects.Create(ctx, &xrplsale.CreateProjectRequest{
        Name:        "My DeFi Protocol",
        Description: "Revolutionary DeFi protocol on XRPL",
        TokenSymbol: "MDP",
        TotalSupply: "100000000",
        Tiers: []xrplsale.Tier{
            {
                Tier:          1,
                PricePerToken: "0.001",
                TotalTokens:   "20000000",
            },
        },
        SaleStartDate: "2025-02-01T00:00:00Z",
        SaleEndDate:   "2025-03-01T00:00:00Z",
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Project created: %s\n", project.ID)
}
```

## Authentication

### XRPL Wallet Authentication

```go
ctx := context.Background()

// Generate authentication challenge
challenge, err := client.Auth.GenerateChallenge(ctx, "rYourWalletAddress...")
if err != nil {
    log.Fatal(err)
}

// Sign the challenge with your wallet
// (implementation depends on your wallet library)
signature := signMessage(challenge.Challenge)

// Authenticate
authResponse, err := client.Auth.Authenticate(ctx, &xrplsale.AuthRequest{
    WalletAddress: "rYourWalletAddress...",
    Signature:     signature,
    Timestamp:     challenge.Timestamp,
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Authentication successful: %s\n", authResponse.Token)

// The token is automatically set in the client for subsequent requests
```

## Core Services

### Projects Service

```go
ctx := context.Background()

// List active projects
projects, err := client.Projects.GetActive(ctx, 1, 10)
if err != nil {
    log.Fatal(err)
}

// Get project details
project, err := client.Projects.Get(ctx, "proj_abc123")
if err != nil {
    log.Fatal(err)
}

// Launch a project
launchedProject, err := client.Projects.Launch(ctx, "proj_abc123")
if err != nil {
    log.Fatal(err)
}

// Get project statistics
stats, err := client.Projects.GetStats(ctx, "proj_abc123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total raised: %s XRP\n", stats.TotalRaisedXRP)
```

### Investments Service

```go
ctx := context.Background()

// Create an investment
investment, err := client.Investments.Create(ctx, &xrplsale.CreateInvestmentRequest{
    ProjectID:       "proj_abc123",
    AmountXRP:       "100",
    InvestorAccount: "rInvestorAddress...",
})

// List investments for a project
investments, err := client.Investments.GetByProject(ctx, "proj_abc123", 1, 10)

// Get investor summary
summary, err := client.Investments.GetInvestorSummary(ctx, "rInvestorAddress...")

// Simulate an investment
simulation, err := client.Investments.Simulate(ctx, &xrplsale.SimulateInvestmentRequest{
    ProjectID: "proj_abc123",
    AmountXRP: "100",
})
fmt.Printf("Expected tokens: %s\n", simulation.TokenAmount)
```

### Analytics Service

```go
ctx := context.Background()

// Get platform analytics
analytics, err := client.Analytics.GetPlatformAnalytics(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total raised: %s XRP\n", analytics.TotalRaisedXRP)

// Get project-specific analytics
startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)
projectAnalytics, err := client.Analytics.GetProjectAnalytics(ctx, "proj_abc123", startDate, endDate)

// Get market trends
trends, err := client.Analytics.GetMarketTrends(ctx, "30d")

// Export data
export, err := client.Analytics.ExportData(ctx, &xrplsale.ExportDataRequest{
    Type:      "projects",
    Format:    "csv",
    StartDate: "2025-01-01",
    EndDate:   "2025-01-31",
})
fmt.Printf("Download URL: %s\n", export.DownloadURL)
```

## Webhook Integration

### HTTP Handler

```go
package main

import (
    "encoding/json"
    "io"
    "net/http"
    
    "github.com/xrplsale/go-sdk"
)

func webhookHandler(client *xrplsale.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Read the request body
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Failed to read body", http.StatusBadRequest)
            return
        }
        
        // Verify signature
        signature := r.Header.Get("X-XRPL-Sale-Signature")
        if !client.VerifyWebhookSignature(body, signature) {
            http.Error(w, "Invalid signature", http.StatusUnauthorized)
            return
        }
        
        // Parse webhook event
        event, err := client.ParseWebhookEvent(body)
        if err != nil {
            http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
            return
        }
        
        // Handle different event types
        switch event.Type {
        case "investment.created":
            handleNewInvestment(event.Data)
        case "project.launched":
            handleProjectLaunched(event.Data)
        case "tier.completed":
            handleTierCompleted(event.Data)
        }
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }
}

func handleNewInvestment(data map[string]interface{}) {
    // Process new investment
    fmt.Printf("New investment: %v XRP\n", data["amount_xrp"])
}

func main() {
    client := xrplsale.NewClient("your-api-key")
    
    http.HandleFunc("/webhooks", webhookHandler(client))
    http.ListenAndServe(":8080", nil)
}
```

### Gin Framework Integration

```go
func setupWebhooks(router *gin.Engine, client *xrplsale.Client) {
    router.POST("/webhooks", func(c *gin.Context) {
        body, _ := c.GetRawData()
        signature := c.GetHeader("X-XRPL-Sale-Signature")
        
        if !client.VerifyWebhookSignature(body, signature) {
            c.JSON(401, gin.H{"error": "Invalid signature"})
            return
        }
        
        event, err := client.ParseWebhookEvent(body)
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid payload"})
            return
        }
        
        // Process event...
        c.JSON(200, gin.H{"status": "ok"})
    })
}
```

## Error Handling

```go
import "github.com/xrplsale/go-sdk"

project, err := client.Projects.Get(ctx, "invalid-id")
if err != nil {
    switch e := err.(type) {
    case *xrplsale.NotFoundError:
        fmt.Println("Project not found")
    case *xrplsale.AuthenticationError:
        fmt.Println("Authentication failed")
    case *xrplsale.ValidationError:
        fmt.Printf("Validation error: %s\n", e.Message)
        fmt.Printf("Details: %v\n", e.Details)
    case *xrplsale.RateLimitError:
        fmt.Printf("Rate limit exceeded. Retry after: %s\n", e.RetryAfter)
    case *xrplsale.APIError:
        fmt.Printf("API error: %s (Code: %s)\n", e.Message, e.Code)
    default:
        fmt.Printf("Unexpected error: %v\n", err)
    }
}
```

## Context and Timeouts

```go
// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Use context for cancellation
projects, err := client.Projects.List(ctx, nil)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("Request timeout")
    }
}

// Cancel long-running operations
go func() {
    time.Sleep(2 * time.Second)
    cancel() // This will cancel any ongoing requests
}()
```

## Configuration Options

```go
client := xrplsale.NewClientWithConfig(&xrplsale.Config{
    APIKey:        "your-api-key",              // Required
    Environment:   xrplsale.Production,         // or xrplsale.Testnet
    BaseURL:       "",                          // Custom API URL (optional)
    Timeout:       30 * time.Second,            // Request timeout
    MaxRetries:    3,                           // Maximum retry attempts
    RetryWaitTime: 1 * time.Second,             // Base wait time between retries
    WebhookSecret: "your-webhook-secret",       // For webhook verification
    Debug:         false,                       // Enable debug logging
})
```

## Pagination

```go
// List projects with pagination
response, err := client.Projects.List(ctx, &xrplsale.ListProjectsOptions{
    Status:    "active",
    Page:      1,
    Limit:     50,
    SortBy:    "created_at",
    SortOrder: "desc",
})

if err != nil {
    log.Fatal(err)
}

for _, project := range response.Data {
    fmt.Printf("Project: %s\n", project.Name)
}

fmt.Printf("Page %d of %d\n", response.Pagination.Page, response.Pagination.TotalPages)
fmt.Printf("Total projects: %d\n", response.Pagination.Total)
```

## Concurrent Operations

```go
// Fetch multiple projects concurrently
projectIDs := []string{"proj_123", "proj_456", "proj_789"}
projects := make([]*xrplsale.Project, len(projectIDs))
errors := make([]error, len(projectIDs))

var wg sync.WaitGroup
for i, id := range projectIDs {
    wg.Add(1)
    go func(index int, projectID string) {
        defer wg.Done()
        project, err := client.Projects.Get(ctx, projectID)
        projects[index] = project
        errors[index] = err
    }(i, id)
}
wg.Wait()

// Process results
for i, project := range projects {
    if errors[i] != nil {
        fmt.Printf("Error fetching project %s: %v\n", projectIDs[i], errors[i])
    } else {
        fmt.Printf("Fetched project: %s\n", project.Name)
    }
}
```

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

## Development

```bash
# Clone the repository
git clone https://github.com/xrplsale/go-sdk.git
cd go-sdk

# Get dependencies
go mod download

# Run tests
go test ./...

# Build
go build ./...

# Format code
go fmt ./...

# Run linter
golangci-lint run
```

## Support

- 📖 [Documentation](https://docs.xrpl.sale)
- 💬 [Discord Community](https://discord.gg/xrpl-sale)
- 🐛 [Issue Tracker](https://github.com/xrplsale/go-sdk/issues)
- 📧 [Email Support](mailto:developers@xrpl.sale)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Links

- [XRPL.Sale Platform](https://xrpl.sale)
- [API Documentation](https://docs.xrpl.sale/api)
- [Other SDKs](https://docs.xrpl.sale/developers/sdk-downloads)
- [GitHub Organization](https://github.com/xrplsale)

---

Made with ❤️ by the XRPL.Sale team