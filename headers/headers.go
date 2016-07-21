// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package headers defines a middleware that adds static headers to
// the requests.
package headers

import "net/http"

// Headers is an http.Header map that implements the httpmw.Wrapper
// interface so that the headers are added to each request using the
// middleware.
type Headers http.Header

// Add adds the value v to the header k.
func (hd Headers) Add(k, v string) {
	http.Header(hd).Add(k, v)
}

// Set sets the value v to the header k, replacing any existing value.
func (hd Headers) Set(k, v string) {
	http.Header(hd).Set(k, v)
}

// Get returns the first value of the header k.
func (hd Headers) Get(k string) string {
	return http.Header(hd).Get(k)
}

// Del removes the header k.
func (hd Headers) Del(k string) {
	http.Header(hd).Del(k)
}

// Wrap returns a handler that adds the headers to the response's Header.
// Values are added, so that if a header key already exists, values are
// appended.
func (hd Headers) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range hd {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		h.ServeHTTP(w, r)
	})
}
