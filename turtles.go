// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package turtles supports creating middleware chains for HTTP
// handlers.
package turtles

import (
	"bytes"
	"fmt"
	"net/http"
)

// Logger defines the Log method that is used to log structured
// data, in tuples of alternating keys/values. The go-kit logger
// satisfies this interface (github.com/go-kit/kit/log).
type Logger interface {
	Log(...interface{}) error
}

// PrintfLogger is an adapter to use Printf-style functions as a
// Logger in the middlewares that accept one. For example,
// the stdlib's log.Printf function can be used via this adapter.
type PrintfLogger func(string, ...interface{})

// Log implements Logger for the PrintfLogger function adapter.
func (fn PrintfLogger) Log(args ...interface{}) error {
	var buf bytes.Buffer
	for i := 0; i < len(args)-1; i += 2 {
		if i > 0 {
			buf.WriteByte(' ')
		}
		verb := "%v" //catch-all formatter
		val := args[i+1]

		// use quoted string formatter if possible
		switch val.(type) {
		case string:
			verb = "%q"
		case fmt.Stringer:
			verb = "%q"
		}

		fmt.Fprintf(&buf, "%s="+verb, args[i], args[i+1])
	}
	if buf.Len() > 0 {
		fn(buf.String())
	}
	return nil
}

// StatusHandler is an integer that handles HTTP requests by writing itself
// as status code. No body is sent.
type StatusHandler int

// ServeHTTP implements http.Handler for the StatusHandler.
func (s StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(s))
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
