// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package headers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/httpmw"
	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {
	head := make(Headers)
	head.Add("A", "a")
	head.Add("B", "b")
	head.Add("C", "c")

	h := httpmw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("D", "d")
		w.Header().Add("A", "z")
		w.Header().Set("B", "x")
		w.WriteHeader(200)
	}), head)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "/", nil)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code, "status")
	assert.Equal(t, 3, len(head), "length of middleware")
	assert.Equal(t, 4, len(w.HeaderMap), "length of response")

	assert.Equal(t, map[string][]string{
		"A": {"a"},
		"B": {"b"},
		"C": {"c"},
	}, map[string][]string(head), "middleware content")

	assert.Equal(t, map[string][]string{
		"A": {"a", "z"},
		"B": {"x"},
		"C": {"c"},
		"D": {"d"},
	}, map[string][]string(w.HeaderMap), "response content")
}
