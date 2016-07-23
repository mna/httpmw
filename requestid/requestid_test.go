// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package requestid

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/httpmw"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	cases := []struct {
		rid    *RequestID
		preset string
	}{
		{&RequestID{}, ""},
		{&RequestID{Header: "XYZ"}, ""},
		{&RequestID{Header: "XYZ"}, "abc"},
		{&RequestID{Header: "XYZ", Len: 12}, ""},
		{&RequestID{Header: "XYZ", Len: 12}, "abc"},
		{&RequestID{Header: "XYZ", Len: 12, ForceSet: true}, ""},
		{&RequestID{Header: "XYZ", Len: 12, ForceSet: true}, "abc"},
	}
	for i, c := range cases {
		key := "X-Request-Id"
		if c.rid.Header != "" {
			key = c.rid.Header
		}

		h := httpmw.Wrap(httpmw.StatusHandler(200), c.rid)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("", "/", nil)
		if c.preset != "" {
			r.Header.Set(key, c.preset)
		}
		h.ServeHTTP(w, r)

		assert.Equal(t, 200, w.Code, "%d: status", i)
		got := r.Header.Get(key)
		t.Logf("%d: got request ID %q", i, got)

		if c.preset != "" && !c.rid.ForceSet {
			assert.Equal(t, c.preset, got, "%d: id", i)
			continue
		}
		wantLen := c.rid.Len
		if wantLen == 0 {
			wantLen = 8
		}
		assert.Equal(t, wantLen, len(got), "%d: length", i)
		assert.NotEqual(t, c.preset, got, "%d: not preset value", i)
	}
}
