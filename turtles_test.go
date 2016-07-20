package turtles

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/turtles/augmentedrw"
	"github.com/stretchr/testify/assert"
)

func TestAugmentedRW(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := w.(interface {
			Size() int
			Status() int
		})
		assert.True(t, ok, "implements Size and Status")
		w.Write(nil)
	})

	hh := Wrap(h, WrapperFunc(augmentedrw.Wrap))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "/", nil)
	hh.ServeHTTP(w, r)
}
