// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timeout

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PuerkitoBio/httpmw"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	cases := []struct {
		sleep   string
		timeout string
		want    int
	}{
		{"10ms", "20ms", 200},
		{"10ms", "5ms", 503},
		{"10ms", "0", 503},
	}
	for i, c := range cases {
		dur, _ := time.ParseDuration(c.timeout)
		to := &Timeout{Duration: dur}
		h := httpmw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dur, _ := time.ParseDuration(c.sleep)
			time.Sleep(dur)
			w.WriteHeader(200)
		}), to)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "/", nil)
		h.ServeHTTP(w, r)

		assert.Equal(t, c.want, w.Code, "%d: status", i)
	}
}
