// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package recover implements a middleware that recovers from
// panics.
package recover

import (
	"net/http"
	"runtime"

	"github.com/PuerkitoBio/httpmw"
)

// Recover holds the configuration for the middleware to recover
// from panics.
type Recover struct {
	// Logger is used to log the panic's details, if non-nil.
	Logger httpmw.Logger

	// StackTrace indicates if the stack trace should be logged
	// in addition to the panic.
	StackTrace bool
}

// Wrap returns a handler that recovers from panics by returning a
// 500 status code and optionally logging the panic and stack trace.
func (rv *Recover) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				if rv.Logger != nil {
					args := []interface{}{"panic", e}
					if rv.StackTrace {
						b := make([]byte, 4096)
						if n := runtime.Stack(b, false); n > 0 {
							args = append(args, "stack", string(b[:n]))
						}
					}
					rv.Logger.Log(args...)
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
