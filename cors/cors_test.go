// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCORSWrap(t *testing.T) {
	allowedOrigins := []string{
		"http://a",
		"https://b",
	}
	cases := []struct {
		origin string
		code   int
	}{
		{"", 204},                // not a CORS request
		{"blah", 403},            // invalid origin
		{allowedOrigins[0], 204}, // ok
		{allowedOrigins[1], 204}, // ok
	}

	for i, c := range cases {
		opts := &CORS{
			ExposeHeaders:    []string{"X-Csrf-Token"},
			AllowCredentials: true,
			AllowedOrigins:   allowedOrigins,
		}

		fn := func(w http.ResponseWriter, r *http.Request) {
			ori := w.Header().Get("Access-Control-Allow-Origin")
			assert.Equal(t, c.origin, ori, "%d: allowed origin", i)

			if c.origin != "" {
				assert.Equal(t, "X-Csrf-Token", w.Header().Get("Access-Control-Expose-Headers"), "%d: expose headers", i)
				assert.Empty(t, w.Header().Get("Vary"), "%d: vary", i)
			}
			w.WriteHeader(204)
		}

		h := opts.Wrap(http.HandlerFunc(fn))
		rr := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://host", nil)
		require.NoError(t, err)
		if c.origin != "" {
			r.Header.Set("Origin", c.origin)
		}

		h.ServeHTTP(rr, r)
		assert.Equal(t, c.code, rr.Code, "status code")
	}
}

func TestCORSHandler(t *testing.T) {
	allowedOrigins := []string{
		"http://a",
		"https://b",
	}
	cases := []struct {
		origin    string
		code      int
		wantMeths string
	}{
		{"", 200, ""},                         // not a CORS request
		{"blah", 403, ""},                     // invalid origin
		{allowedOrigins[0], 200, "GET, POST"}, // ok
		{allowedOrigins[1], 200, "GET, POST"}, // ok
	}

	for i, c := range cases {
		opts := &CORS{
			AllowCredentials: true,
			ExposeHeaders:    []string{"X-Csrf-Token"},
			MaxAge:           time.Hour,
			AllowedOrigins:   allowedOrigins,
			AllowedMethods:   []string{"GET", "POST"},
		}

		rr := httptest.NewRecorder()
		r, err := http.NewRequest("OPTIONS", "http://host", nil)
		require.NoError(t, err)
		if c.origin != "" {
			r.Header.Set("Origin", c.origin)
		}

		opts.ServeHTTP(rr, r)

		if assert.Equal(t, c.code, rr.Code, "%d: code", i) && c.code < 400 {
			got := rr.Header().Get("Access-Control-Allow-Methods")
			assert.Equal(t, c.wantMeths, got, "%d: allowed methods", i)

			ori := rr.Header().Get("Access-Control-Allow-Origin")
			assert.Equal(t, c.origin, ori, "%d: allowed origin", i)

			if c.origin != "" {
				assert.Equal(t, "X-Csrf-Token", rr.Header().Get("Access-Control-Expose-Headers"), "%d: expose headers", i)
				assert.Equal(t, "Origin", rr.Header().Get("Vary"), "%d: vary", i)
			}
		}
	}
}
