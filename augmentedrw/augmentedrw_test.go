package augmentedrw

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/turtles"
	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww, ok := w.(interface {
			Size() int
			Status() int
		})
		assert.True(t, ok, "implements Size and Status")
		fmt.Fprint(w, "ok")
		assert.Equal(t, 2, ww.Size(), "size")
	})

	hh := turtles.Wrap(h, turtles.WrapperFunc(Wrap))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "/", nil)
	hh.ServeHTTP(w, r)

	assert.Equal(t, w.Code, 200, "status")
}
