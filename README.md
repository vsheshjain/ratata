# Ratata - Rate Limiter

**Ratata** is a lightweight rate-limiting library implemented in Go, designed to manage the number of requests or actions allowed for users over a given time. It uses the token bucket algorithm to ensure that actions or requests are handled efficiently while respecting a defined limit.

## Features

- Token bucket-based rate limiting.
- User-specific rate limits.
- Thread-safe implementation.
- Customizable bucket capacity and refill rates.

## Installation

To use Ratata in your Go project, add it to your dependencies:

```bash
go get github.com/vsheshjain/ratata
```

## Usage

### Basic Usage

Create a new rate-limiting bucket by defining the capacity and refill rate:

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/ratata"
)

func main() {
    // Create a new rate limiter with a capacity of 5 tokens, refilling every 1 second
    bucket := ratata.NewRatataBucket(5, time.Second)

    for i := 0; i < 10; i++ {
        allowed := bucket.Allow()
        if allowed {
            fmt.Println("Request allowed")
        } else {
            fmt.Println("Rate limit exceeded")
        }
        time.Sleep(200 * time.Millisecond)  // Simulate some delay between requests
    }
}
```

### Gin Web Framework Example
You can easily integrate Ratata with the Gin web framework to limit requests per user by incorporating it in your auth middleware:

```go
// UserRateLimiter returns a Gin middleware that applies rate limiting
// based on the user ID passed as a query parameter. 
func UserRateLimiter(tb *ratata.RatataBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is allowed to proceed based on their user ID
		if tb.AllowUser(c.Query("user_id")) {
			c.Next() // Proceed to the next handler
		} else {
			c.AbortWithStatus(http.StatusTooManyRequests) // Abort and send 429 response
		}
	}
}

// NewRouter initializes a new Gin router and sets up the routes.
// It creates a new RatataBucket with a capacity of 5 tokens that refills
// every 10 seconds. The /health endpoint is protected by the UserRateLimiter
// middleware, returning a health status if allowed.
func NewRouter(ctx context.Context) *gin.Engine {
	router := gin.New()
	tb := ratata.NewRatataBucket(5, 10*time.Second) // Create a new rate limiter

	// Define the /health endpoint with rate limiting
	router.GET("/health", UserRateLimiter(tb), func(c *gin.Context) {
		health := serviceDetails{
			Message: "We are up.",
			Time:    time.Now().Unix(), // Current Unix timestamp
		}
		c.JSON(http.StatusOK, health) // Return the health status
	})
	return router
}


```
