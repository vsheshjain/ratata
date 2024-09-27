// Package ratata provides a token bucket rate-limiting implementation for controlling
// the rate of actions performed by users based on their unique user IDs. The library
// helps manage requests and ensure users don't exceed their allowed rate limits.
package ratata

import (
	"sync"
	"time"
)

// RatataBucket represents a token bucket with a defined capacity and refill rate.
// It controls the rate of actions for a user based on the number of available tokens.
type RatataBucket struct {
	capacity   int           // Maximum number of tokens the bucket can hold.
	tokens     int           // Current number of tokens in the bucket.
	refillRate time.Duration // Duration to wait before adding a new token.
	lastRefill time.Time     // Time of the last token refill.
	mu         sync.Mutex    // Mutex to protect concurrent access to the bucket's fields.
}

var (
	userBuckets = make(map[string]*RatataBucket) // Map to store token buckets for each user.
	bucketMu    sync.Mutex                       // Mutex to protect concurrent access to the userBuckets map.
)

// NewRatataBucket creates and returns a new token bucket with a specified capacity and refill rate.
func NewRatataBucket(capacity int, refillRate time.Duration) *RatataBucket {
	return &RatataBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// refillRatata refills the bucket with tokens based on the elapsed time since the last refill.
// It ensures that tokens do not exceed the bucket's capacity.
func (rb *RatataBucket) refillRatata() {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rb.lastRefill)

	// Calculate how many tokens to add based on the time elapsed and refill rate.
	newTokens := int(elapsed / rb.refillRate)

	if newTokens > 0 {
		rb.tokens += newTokens
		if rb.tokens > rb.capacity {
			rb.tokens = rb.capacity // Ensure tokens do not exceed capacity.
		}
		rb.lastRefill = now // Update the last refill time.
	}
}

// Allow checks if a token is available and consumes one if so.
// Returns true if an action is allowed (token available), false otherwise.
func (rb *RatataBucket) Allow() bool {
	rb.refillRatata() // Refill tokens before allowing the action.

	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.tokens > 0 {
		rb.tokens-- // Consume one token.
		return true
	}
	return false
}

// AllowUser checks or creates a token bucket for a specific user and then checks if an action is allowed.
// It returns true if the user is allowed to perform the action (token available), false otherwise.
func (rb *RatataBucket) AllowUser(userID string) bool {
	bucketMu.Lock()
	defer bucketMu.Unlock()
	
	// Initialize a new bucket for the user if it doesn't exist.
	if userBuckets[userID] == nil {
		userBuckets[userID] = NewRatataBucket(rb.capacity, rb.refillRate)
	}
	userBucket := userBuckets[userID]

	return userBucket.Allow()
}
