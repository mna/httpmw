// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpmw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	nop := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	w := func(char byte) WrapperFunc {
		return func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte{char})
				h.ServeHTTP(w, r)
			})
		}
	}

	cases := []struct {
		h    http.Handler
		code int
		body string
	}{
		{Wrap(StatusHandler(123)), 123, ""},
		{Wrap(nop, w('a')), 200, "a"},
		{Wrap(nop, w('a'), w('b')), 200, "ab"},
		{Wrap(nop, w('a'), w('b'), w('c')), 200, "abc"},
	}
	for i, c := range cases {
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest("", "/", nil)
		c.h.ServeHTTP(rw, req)

		assert.Equal(t, c.code, rw.Code, "%d: status code", i)
		assert.Equal(t, c.body, rw.Body.String(), "%d: body", i)
	}
}

func TestStatusHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "/", nil)

	h := StatusHandler(204)
	h.ServeHTTP(w, r)
	assert.Equal(t, 204, w.Code, "status")
}
