// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrequest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/httpmw"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestLogRequest(t *testing.T) {
	var buf bytes.Buffer
	l := log.NewLogfmtLogger(&buf)
	lr := &LogRequest{Logger: l, DurationFormat: "%.5f", Fields: []string{"duration", "method", "path"}}
	h := httpmw.Wrap(httpmw.StatusHandler(200), lr)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "/", nil)
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code, "status")
	assert.Contains(t, buf.String(), " method=GET path=/", "expected output")
}
