// Package remoteip implements a middleware that extracts the effective
// remote client IP address and sets it on the request's RemoteAddr field.
package remoteip

import (
	"net"
	"net/http"
	"strings"
)

// DefaultHeaders is the list of headers inspected and trusted to get the
// effective remote client IP address.
var DefaultHeaders = []string{"Cf-Connecting-Ip", "X-Forwarded-For", "X-Real-Ip"}

func extractClientIP(h http.Header, keys []string) string {
	for _, key := range keys {
		if v := h.Get(key); v != "" {
			if ips := strings.Fields(strings.TrimSpace(v)); len(ips) > 0 {
				if pip := net.ParseIP(ips[0]); pip != nil {
					return pip.String()
				}
			}
		}
	}
	return ""
}

// RemoteIP holds the configuration for the remote IP middleware.
type RemoteIP struct {
	// Headers is the list of headers to use to get the effective remote
	// client IP address. If it is empty, DefaultHeaders is used.
	Headers []string
}

// Wrap returns a handler that assigns the request's RemoteAddr field
// if it finds a valid IP address in the configured header keys, before
// calling the handler h.
func (rip *RemoteIP) Wrap(h http.Handler) http.Handler {
	keys := rip.Headers
	if len(keys) == 0 {
		keys = DefaultHeaders
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rip := extractClientIP(r.Header, keys); rip != "" {
			r.RemoteAddr = rip
		}
		h.ServeHTTP(w, r)
	})
}
