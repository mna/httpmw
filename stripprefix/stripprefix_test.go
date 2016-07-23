// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stripprefix

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/httpmw"
	"github.com/stretchr/testify/assert"
)

func TestStripPrefix(t *testing.T) {
	cases := []struct {
		prefix string
		path   string
		want   string
		code   int
	}{
		{"", "/", "/", 200},
		{"/api", "/api/x", "/x", 200},
		{"/api", "/blah/x", "/blah/x", 404},
	}

	for i, c := range cases {
		sp := &StripPrefix{c.prefix}
		h := httpmw.Wrap(httpmw.StatusHandler(200), sp)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("", c.path, nil)

		h.ServeHTTP(w, r)
		assert.Equal(t, c.code, w.Code, "%d: status", i)
		assert.Equal(t, c.want, r.URL.Path, "%d: path", i)
	}
}
