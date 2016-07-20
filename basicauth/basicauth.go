// Package basicauth implements a basic authentication middleware.
package basicauth

import (
	"fmt"
	"net/http"
)

// BasicAuth holds the configuration for the basic authentication
// middleware.
type BasicAuth struct {
	// User is the valid username.
	User string

	// Password is the valid password.
	Password string

	// AuthFunc is the function called to authenticated the provided
	// user and password. It returns true if the credentials are valid,
	// false if they are not, or an error if the check failed, in which
	// case the middleware returns a status code 500.
	//
	// If an AuthFunc is specified, User and Password are ignored.
	AuthFunc func(string, string) (bool, error)

	// Realm is the realm of the basic authentication, specified in the
	// WWW-Authenticate header when the authentication fails.
	Realm string
}

// Wrap returns a handler that validates the authentication credentials
// before calling the handler h.
func (ba *BasicAuth) Wrap(h http.Handler) http.Handler {
	fn := ba.AuthFunc
	if fn == nil {
		fn = func(u, p string) (bool, error) {
			return u == ba.User && p == ba.Password, nil
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if ok {
			success, err := fn(u, p)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			ok = success
		}
		if !ok {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", ba.Realm))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
