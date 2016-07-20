// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cors implements a CORS middleware and handler for OPTIONS requests.
package cors

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CORS holds the configuration for the cors middleware.
type CORS struct {
	// AllowCredentials indicates if the request is allowed to send
	// credentials information along with the request.
	AllowCredentials bool

	// MaxAge sets the time to cache the OPTIONS response.
	MaxAge time.Duration

	// AllowHeaders is list of headers allowed by the CORS endpoint.
	AllowHeaders []string

	// ExposeHeaders is a list of headers exposed to the client by the CORS endpoint.
	ExposeHeaders []string

	// AllowedOrigins is a list of whitelisted, allowed origins. It should
	// include the scheme and port, as required. The value is lowercased
	// and compared to the Origin header received from the client. The
	// special value "*" can be used to allow any origin.
	//
	// Ignored if AllowedOriginsFunc is set.
	AllowedOrigins []string

	// AllowedOriginsFunc is a function that is called with the Origin header
	// received from the client. The origin is allowed if it returns true.
	// If this field is non-nil, AllowedOrigins is ignored and this function
	// is called for each request.
	AllowedOriginsFunc func(string) bool

	// AllowedMethods indicates the list of allowed HTTP methods. The
	// OPTIONS method should not be included.
	AllowedMethods []string
}

// ServeHTTP is the handler for CORS OPTIONS preflight requests.
// It sets the CORS headers and returns either 200 if the request
// is allowed, or 403 if it isn't.
func (c *CORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := setCORSHeaders(w, r, c); err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Add vary by origin for cache
	vary := w.Header().Get("Vary")
	if vary != "" {
		vary += ", "
	}
	vary += "Origin"
	w.Header().Set("Vary", vary)

	w.WriteHeader(200)
}

// Wrap returns a handler that sets the CORS headers to allow a
// cross-origin request if the request origin is a whitelisted origin.
// The handler h is called only if the request is allowed by the CORS
// policy.
func (c *CORS) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := setCORSHeaders(w, r, c); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func isIn(list []string, v string) bool {
	for _, vv := range list {
		if vv == "*" {
			return true
		}
		if strings.ToLower(vv) == strings.ToLower(v) {
			return true
		}
	}
	return false
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request, c *CORS) error {
	ori := r.Header.Get("Origin")
	if ori == "" {
		// not a CORS request, passthrough
		return nil
	}

	if c.AllowedOriginsFunc != nil {
		if !c.AllowedOriginsFunc(ori) {
			return errors.New("invalid CORS origin: " + ori)
		}
	} else if !isIn(c.AllowedOrigins, ori) {
		return errors.New("invalid CORS origin: " + ori)
	}

	// simple CORS request (not preflight OPTIONS) requires only
	// allow-origin, credentials and expose headers.
	w.Header().Set("Access-Control-Allow-Origin", ori)
	if c.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.ExposeHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(c.ExposeHeaders, ", "))
	}

	// other CORS headers are only for preflight requests (OPTIONS)
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.AllowedMethods, ", "))
		if c.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(int(c.MaxAge.Seconds())))
		}
		if len(c.AllowHeaders) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.AllowHeaders, ", "))
		}
	}
	return nil
}
