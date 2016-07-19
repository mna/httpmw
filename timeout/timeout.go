// Package timeout implements a middleware that returns a 503 error
// if the request takes too long to execute.
package timeout

import (
	"net/http"
	"time"
)

// Timeout holds the configuration for the timeout middleware.
type Timeout struct {
	// Duration is the time allowed for the request to run.
	Duration time.Duration

	// Message is the message returned with the 503 status code if
	// the request timed out.
	Message string
}

// Wrap returns a handler that must run in the configured Duration
// otherwise it returns a status code 503. It calls http.TimeoutHandler
// to create the returned handler.
func (t *Timeout) Wrap(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, t.Duration, t.Message)
}
