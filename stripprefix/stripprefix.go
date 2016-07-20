// Package stripprefix implements a middleware handler that strips a prefix
// from the request URL's Path.
package stripprefix

import "net/http"

// StripPrefix holds the configuration for the StripPrefix middleware handler.
type StripPrefix struct {
	// Prefix is the prefix to remove from the request URL's Path.
	Prefix string
}

// Wrap returns a handler that strips the prefix from the request URL's
// Path. It calls http.StripPrefix to create the returned handler.
func (sp *StripPrefix) Wrap(h http.Handler) http.Handler {
	return http.StripPrefix(sp.Prefix, h)
}