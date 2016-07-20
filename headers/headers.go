// Package headers defines a middleware that adds static headers to
// the requests.
package headers

import "net/http"

// Headers is an http.Header map that implements the turtles.Wrapper
// interface so that the headers are added to each request using the
// middleware.
type Headers http.Header

// Wrap returns a handler that adds the headers to the request's Header.
// Values are added, so that if a header key already exists, values are
// appended.
func (hd Headers) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range hd {
			for _, vv := range v {
				r.Header.Add(k, vv)
			}
		}
		h.ServeHTTP(w, r)
	})
}
