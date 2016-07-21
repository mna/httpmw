// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cleanpath implements a middleware that cleans the requested
// path to a canonical form. It redirects to the clean path if the
// requested path is not the canonical one. Some request multiplexers
// already do this automatically to some extent (e.g. http.ServeMux cleans
// the . and .. parts of the path, others handle the trailing slash),
// but this middleware handles this uniformly regardless of the mux used.
package cleanpath

import (
	"net/http"
	"path"
)

// TrailingSlashMode specifies how the trailing slash should be
// handled on the request URL's Path.
type TrailingSlashMode int

const (
	// Leave keeps the trailing slash the way it was received.
	Leave TrailingSlashMode = iota
	// Add enforces the presence of a trailing slash.
	Add
	// Remove enforces the removal of a trailing slash.
	Remove
)

// CleanPath holds the configuration for the middleware.
type CleanPath struct {
	// TrailingSlash specifies how the trailing slash of the path should
	// be handled. Defaults to Leave, which keeps it as it was received.
	// The mode does not apply to the root slash, which is always present
	// if the path is otherwise empty.
	TrailingSlash TrailingSlashMode
}

// Wrap returns a handler that redirects with a status code 301 to the
// canonical path if the requested path is not as expected. It cleans
// the . and .. parts and handles the trailing slash according to the
// middleware configuration.
//
// If the path is already in a canonical form, it calls the handler h.
func (cp *CleanPath) Wrap(h http.Handler) http.Handler {
	mode := cp.TrailingSlash
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "CONNECT" {
			if p := cleanPath(r.URL.Path, mode); p != r.URL.Path {
				url := *r.URL
				url.Path = p
				http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// cleanPath comes from https://golang.org/src/net/http/server.go,
// adapted to support the trailing slash mode.
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func cleanPath(p string, mode TrailingSlashMode) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	switch mode {
	case Add:
		if np[len(np)-1] != '/' {
			np += "/"
		}
	case Remove:
		if len(np) > 1 && np[len(np)-1] == '/' {
			np = np[:len(np)-1]
		}
	}
	return np
}
