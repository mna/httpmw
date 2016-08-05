// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
// RPS=100, MaxWait=50ms.
type RateLimit struct {
	// RPS is the number of requests per seconds. Tokens will fill at an
	// interval that closely respects that RPS value.
	RPS int64

	// Capacity is the maximum number of tokens that can be available in
	// the bucket. The bucket starts at full capacity. If the capacity is
	// <= 0, it is set to the RPS.
	Capacity int64

	// MaxWait is the maximum time to wait for an available token for a
	// request to be allowed. If no token is available, the request is
	// denied without waiting and a status code 429 is returned.
	MaxWait time.Duration
}

// Wrap returns a handler that allows only the configured number of requests.
// The wrapped handler h is called only if the request is allowed by the rate
// limiter, otherwise a status code 429 is returned.
//
// Each call to Wrap creates a new, distinct rate limiter bucket that controls
// access to h.
func (rl *RateLimit) Wrap(h http.Handler) http.Handler {
	cap := rl.Capacity
	if rl.Capacity <= 0 {
		cap = rl.RPS
	}
	bucket := ratelimit.NewBucketWithRate(float64(rl.RPS), cap)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !bucket.WaitMaxDuration(1, rl.MaxWait) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		h.ServeHTTP(w, r)
	})
}
