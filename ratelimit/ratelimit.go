// Package ratelimit implements a rate limiter middleware handler.
// It uses the github.com/juju/ratelimit package to control the rate.
// This package uses a token bucket implementation.
package ratelimit

import (
	"net/http"
	"time"

	"github.com/juju/ratelimit"
)

// RateLimit holds the configuration for the RateLimit middleware handler.
// For example, to allow 100 requests per second and a wait time of 50ms
// to wait for an available token for the request to be allowed, set
// Requests=100, Interval=1s, MaxWait=50ms.
type RateLimit struct {
	// Requests is the maximum number of requests to allow per Interval.
	Requests int64

	// Interval is the time interval during which the specified number of
	// Requests are allowed.
	Interval time.Duration

	// MaxWait is the maximum time to wait for an available token for a
	// request to be allowed. If no token is available, the request is
	// denied and a status code 429 is returned.
	MaxWait time.Duration
}

// Wrap returns a handler that allows only the configured number of requests
// per specified time interval. The wrapped handler h is called only if the
// request is allowed by the rate limiter, otherwise a status code 429 is returned.
//
// Each call to Wrap creates a new, distinct rate limiter bucket that controls
// access to h.
func (rl *RateLimit) Wrap(h http.Handler) http.Handler {
	bucket := ratelimit.NewBucketWithQuantum(rl.Interval, rl.Requests, rl.Requests)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !bucket.WaitMaxDuration(1, rl.MaxWait) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		h.ServeHTTP(w, r)
	})
}
