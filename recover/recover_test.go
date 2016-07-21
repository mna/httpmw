// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recover

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/turtles"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(io.EOF)
	})

	var rec Recover
	h := turtles.Wrap(fn, &rec)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "/", nil)

	h.ServeHTTP(w, r)
	assert.Equal(t, 500, w.Code, "status")
}
