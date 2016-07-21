// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package turtles

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrintfLogger(t *testing.T) {
	var buf bytes.Buffer
	stdl := log.New(&buf, "", 0)
	l := PrintfLogger(stdl.Printf)

	cases := []struct {
		in  []interface{}
		out string
	}{
		{nil, ""},
		{[]interface{}{"a"}, ""},
		{[]interface{}{"a", 1}, "a=1\n"},
		{[]interface{}{"a", 1, "b"}, "a=1\n"},
		{[]interface{}{"a", 1, "b", 2}, "a=1 b=2\n"},
		{[]interface{}{"a", 1, "b", 2, "c"}, "a=1 b=2\n"},
		{[]interface{}{"a", 1, "b", 2, "c", true}, "a=1 b=2 c=true\n"},
		{[]interface{}{"a", 1, "b", 2, "c", true, "d"}, "a=1 b=2 c=true\n"},
		{[]interface{}{"a", 1, "b", 2, "c", true, "d", "some value"}, "a=1 b=2 c=true d=\"some value\"\n"},
		{[]interface{}{"a", 1, "b", 2, "c", true, "d", time.Second}, "a=1 b=2 c=true d=\"1s\"\n"},
	}
	for i, c := range cases {
		err := l.Log(c.in...)
		if assert.NoError(t, err, "%d: Log", i) {
			assert.Equal(t, c.out, buf.String(), "%d: expected output", i)
		}
		buf.Reset()
	}
}

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
