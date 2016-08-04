// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package httpmw supports creating middleware chains for HTTP
// handlers.
package httpmw

import "net/http"

// StatusHandler is an integer that handles HTTP requests by writing itself
// as status code. No body is sent.
type StatusHandler int

// ServeHTTP implements http.Handler for the StatusHandler.
func (s StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(s))
	w.Write(nil)
}

// Wrapper defines the Wrap method required to build a middleware-style
// chain of calls.
type Wrapper interface {
	Wrap(http.Handler) http.Handler
}

// WrapperFunc is a function type that implements the Wrapper interface,
// useful to adapt middleware from other packages.
type WrapperFunc func(http.Handler) http.Handler

// Wrap implements the Wrapper interface for a WrapperFunc by calling the function
// with h as argument.
func (fn WrapperFunc) Wrap(h http.Handler) http.Handler {
	return fn(h)
}

// Wrap wraps the HTTP handler h with the provided middleware ws.
// It returns a handler that will call
// ws[0] -> ws[1] -> ... -> ws[n-1] -> h.
// Each Wrapper may stop the chain of calls by not
// calling the next handler in the chain.
func Wrap(h http.Handler, ws ...Wrapper) http.Handler {
	for i := len(ws) - 1; i >= 0; i-- {
		h = ws[i].Wrap(h)
	}
	return h
}
