// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ratelimit

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PuerkitoBio/turtles"
	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	rl := &RateLimit{Requests: 2, Interval: 100 * time.Millisecond, MaxWait: 10 * time.Millisecond}
	h := turtles.Wrap(turtles.StatusHandler(200), rl)

	for _, want := range []int{200, 200, 429} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "/", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, want, w.Code, "status")
	}

	// wait the interval, should be good for 2 more
	time.Sleep(100 * time.Millisecond)
	for _, want := range []int{200, 200, 429} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "/", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, want, w.Code, "status")
	}
}