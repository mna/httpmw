// Package augmentedrw implements a middleware that replaces the standard
// http.ResponseWriter with one that records the Size and Status of the
// response. This is primarily useful for logging requests.
package augmentedrw

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// Wrap returns a handler that calls h with an augmented http.ResponseWriter,
// that is, one that records the Size and Status code of the response.
func Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do not create the augmented response writer if it already implements
		// Size and Status.
		if _, ok := w.(interface {
			Size() int
			Status() int
		}); !ok {
			w = &responseWriter{w, 0, 0}
		}
		h.ServeHTTP(w, r)
	})
}

// responseWriter is an augmented response writer that keeps track
// of the response's status and body size.
type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack is not supported")
	}
	return hj.Hijack()
}

func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w *responseWriter) Flush() {
	f, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		f.Flush()
	}
}
