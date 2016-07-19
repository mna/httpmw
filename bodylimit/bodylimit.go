// Package bodylimit implements an HTTP middleware that limits the number
// of bytes that can be read from the request body.
package bodylimit

import "net/http"

// BodyLimit holds the configuration for the middleware handler.
type BodyLimit struct {
	// N is the maximum number of bytes that can be read from the request
	// body before an error is returned.
	N int64
}

// Wrap returns a handler that limits the number of bytes that can be
// read from the request body before calling the handler h. It calls
// http.MaxBytesReader and set the request's body to the returned
// io.ReadCloser.
func (bl *BodyLimit) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, bl.N)
		h.ServeHTTP(w, r)
	})
}
