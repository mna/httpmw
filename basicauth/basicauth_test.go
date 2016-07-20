package basicauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/turtles"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	cases := []struct {
		conf      *BasicAuth
		user, pwd string
		want      int
	}{
		{conf: &BasicAuth{User: "a", Password: "b"}, user: "a", pwd: "b", want: 200},
		{conf: &BasicAuth{User: "a", Password: "b"}, user: "x", pwd: "b", want: 401},
		{conf: &BasicAuth{User: "a", Password: "b"}, user: "a", pwd: "x", want: 401},
		{conf: &BasicAuth{User: "a", Password: "b"}, user: "x", pwd: "x", want: 401},
		{conf: &BasicAuth{User: "a", Password: "b", Realm: "r"}, user: "x", pwd: "x", want: 401},
		{conf: &BasicAuth{AuthFunc: func(u, p string) (bool, error) {
			return u == "a" && p == "b", nil
		}}, user: "a", pwd: "b", want: 200},
		{conf: &BasicAuth{AuthFunc: func(u, p string) (bool, error) {
			return u == "a" && p == "b", nil
		}}, user: "a", pwd: "x", want: 401},
		{conf: &BasicAuth{AuthFunc: func(u, p string) (bool, error) {
			return false, errors.New("error")
		}}, user: "a", pwd: "b", want: 500},
		{conf: &BasicAuth{AuthFunc: func(u, p string) (bool, error) {
			return true, errors.New("error")
		}}, user: "a", pwd: "b", want: 500},
		{conf: &BasicAuth{User: "a", Password: "b", AuthFunc: func(u, p string) (bool, error) {
			return false, nil
		}}, user: "a", pwd: "b", want: 401},
	}
	for i, c := range cases {
		h := turtles.Wrap(turtles.StatusHandler(200), c.conf)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("", "/", nil)
		r.SetBasicAuth(c.user, c.pwd)
		h.ServeHTTP(w, r)

		if assert.Equal(t, c.want, w.Code, "%d: status code", i) && c.want >= 400 && c.want < 500 {
			header := w.Header().Get("WWW-Authenticate")
			want := c.conf.Realm
			if want == "" {
				want = DefaultRealm
			}
			want = fmt.Sprintf("Basic realm=%q", want)
			assert.Equal(t, want, header, "%d: WWW-Authenticate")
		}
	}
}
