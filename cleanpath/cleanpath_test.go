// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cleanpath

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/httpmw"
	"github.com/stretchr/testify/assert"
)

func TestCleanPath(t *testing.T) {
	cases := []struct {
		path    string
		mode    TrailingSlashMode
		newPath string
		want    int
	}{
		{"", Leave, "/", 301},
		{"", Add, "/", 301},
		{"", Remove, "/", 301},

		{"/", Leave, "", 200},
		{"/", Add, "", 200},
		{"/", Remove, "", 200},

		{"a/b", Leave, "/a/b", 301},
		{"a/b", Add, "/a/b/", 301},
		{"a/b", Remove, "/a/b", 301},

		{"/a/b", Leave, "", 200},
		{"/a/b", Add, "/a/b/", 301},
		{"/a/b", Remove, "", 200},

		{"/a/b/", Leave, "", 200},
		{"/a/b/", Add, "", 200},
		{"/a/b/", Remove, "/a/b", 301},

		{"/a/b/./..", Leave, "/a", 301},
		{"/a/b/./..", Add, "/a/", 301},
		{"/a/b/./..", Remove, "/a", 301},

		{"/a/b/./..//", Leave, "/a/", 301},
		{"/a/b/./..//", Add, "/a/", 301},
		{"/a/b/./..//", Remove, "/a", 301},
	}
	for i, c := range cases {
		cp := &CleanPath{TrailingSlash: c.mode}
		h := httpmw.Wrap(httpmw.StatusHandler(200), cp)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("", c.path, nil)

		h.ServeHTTP(w, r)
		assert.Equal(t, c.want, w.Code, "%d: status", i)
		assert.Equal(t, c.newPath, w.Header().Get("Location"), "%d: location", i)
	}
}
