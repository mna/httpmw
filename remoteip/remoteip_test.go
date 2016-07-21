// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package remoteip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/turtles"
	"github.com/stretchr/testify/assert"
)

func TestRemoteIP(t *testing.T) {
	var rip RemoteIP
	h := turtles.Wrap(turtles.StatusHandler(200), &rip)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "/", nil)

	r.Header.Set("Cf-Connecting-Ip", "blah")
	r.Header.Set("X-Forwarded-For", " 12.34.56.78 11.11.111.111")
	h.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code, "status")
	assert.Equal(t, "12.34.56.78", r.RemoteAddr, "remote address")
}
