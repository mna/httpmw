// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package augmentedrw

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/httpmw"
	"github.com/stretchr/testify/assert"
)

func TestAugmentedRW(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww, ok := w.(interface {
			Size() int
			Status() int
		})
		assert.True(t, ok, "implements Size and Status")
		fmt.Fprint(w, "ok")
		assert.Equal(t, 2, ww.Size(), "size")
	})

	hh := httpmw.Wrap(h, httpmw.WrapperFunc(Wrap))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "/", nil)
	hh.ServeHTTP(w, r)

	assert.Equal(t, w.Code, 200, "status")
}
