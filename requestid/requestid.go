// Package requestid implements a middleware that generates a random request
// ID.
package requestid

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

// RequestID holds the configuration for the request ID middleware.
type RequestID struct {
	// ForceSet replaces the existing value for the request ID Header if true.
	// If false, it sets a request ID only if there was none. Defaults to false.
	ForceSet bool

	// Len is the length of the generated request ID. The ID is hex-encoded,
	// this is the length of the final hex-encoded string, not the number of
	// random bytes used. Defaults to 8.
	Len int

	// Header is the name of the header to use to store the request ID. Defaults
	// to X-Request-Id.
	Header string
}

// for tests
var testForceRandErr bool

// Wrap returns a handler that sets a random request ID header before calling
// the handler h.
func (rid *RequestID) Wrap(h http.Handler) http.Handler {
	header := rid.Header
	if header == "" {
		header = "X-Request-Id"
	}
	force := rid.ForceSet
	n := rid.Len
	if n <= 0 {
		n = 8
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// generate an ID if there is none or ForceSet is true
		if id := r.Header.Get(header); id == "" || force {
			// the number of random bytes is Len / 2 (since we then hex-encode the bytes)
			b := make([]byte, hex.DecodedLen(n))
			if _, err := rand.Read(b); err == nil && !testForceRandErr {
				r.Header.Set(header, hex.EncodeToString(b))
			} else {
				// fallback on timestamp
				ts := time.Now().UnixNano()
				v := strconv.FormatInt(ts, 10)
				if len(v) > n {
					// take the last n bytes, more randomness
					v = v[len(v)-n:]
				}
				r.Header.Set(header, v)
			}
		}
		h.ServeHTTP(w, r)
	})
}
