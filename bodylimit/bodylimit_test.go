// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bodylimit

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/turtles"
	"github.com/stretchr/testify/assert"
)

func TestBodyLimit(t *testing.T) {
	cases := []struct {
		N       int64
		wantErr bool
	}{
		{0, false},
		{1, true},
		{9, true},
		{10, false},
		{11, false},
	}
	for i, c := range cases {
		bl := &BodyLimit{N: c.N}
		h := turtles.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.Copy(ioutil.Discard, r.Body)
			if c.wantErr {
				if assert.Error(t, err, "%d: error", i) {
					assert.Contains(t, err.Error(), "http: request body too large", "%d: message", i)
				}
			} else {
				assert.NoError(t, err, "%d: error", i)
			}
			w.WriteHeader(200)
		}), bl)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("", "/", io.LimitReader(rand.Reader, 10))
		h.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code, "%d: status", i)
	}
}
