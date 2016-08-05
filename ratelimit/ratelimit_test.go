// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PuerkitoBio/httpmw"
	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	rl := &RateLimit{RPS: 2, Capacity: 2, MaxWait: 10 * time.Millisecond}
	h := httpmw.Wrap(httpmw.StatusHandler(200), rl)

	for _, want := range []int{200, 200, 429} {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("", "/", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, want, w.Code, "status")
	}

	// wait the interval, should be good for 2 more
	time.Sleep(time.Second)
	for _, want := range []int{200, 200, 429} {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("", "/", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, want, w.Code, "status")
	}
}
