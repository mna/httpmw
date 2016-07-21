// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logrequest implements a middleware that logs requests.
package logrequest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/httpmw"
)

var allFields = []string{
	"body_bytes_received",
	"body_bytes_sent",
	"duration",
	"end",
	"host",
	"id",
	"method",
	"origin",
	"path",
	"proto",
	"remote_addr",
	"start",
	"status",
	"uri",
	"user_agent",
}

// LogRequest holds the configuration for the LogRequest middleware.
type LogRequest struct {
	// Logger is the logger to use to log the requests.
	Logger httpmw.Logger

	// RequestIDHeader is the name of the header that contains the request
	// ID. Defaults to X-Request-Id.
	RequestIDHeader string

	// TimeFormat is the format to use to format timestamps, as supported by
	// the time package.
	TimeFormat string

	// DurationFormat is the format to use to log the duration of the request.
	// The value to format is the number of seconds in float64. Defaults to
	// %.3f (milliseconds).
	DurationFormat string

	// Fields is the list of field names to log. Defaults to all supported
	// fields. The supported fields are:
	//
	//     body_bytes_received: bytes in the request body
	//     body_bytes_sent: bytes in the response body
	//     duration: duration of the request
	//     end: date and time of the end of the request (UTC)
	//     host: host (and possibly port) of the request
	//     id: request ID
	//     method: method of the request (e.g. GET)
	//     origin: value of the Origin request header
	//     path: path section of the request URL
	//     proto: protocol and version (e.g. HTTP/1.1)
	//     remote_addr: address of the client
	//     start: date and time of the start of the request (UTC)
	//     status: status code of the response
	//     uri: raw request URI
	//     user_agent: value of the User-Agent request header
	//
	Fields []string
}

// Wrap returns a handler that records the start time, calls the handler h,
// records the end time and duration, and logs the request's fields as
// configured by the LogRequest.
func (lr *LogRequest) Wrap(h http.Handler) http.Handler {
	log := lr.Logger
	if log == nil {
		return h
	}
	dfmt := lr.DurationFormat
	if dfmt == "" {
		dfmt = "%.3f"
	}
	tf := lr.TimeFormat
	if tf == "" {
		tf = time.ANSIC
	}
	hd := lr.RequestIDHeader
	if hd == "" {
		hd = "X-Request-Id"
	}
	fields := lr.Fields
	if len(fields) == 0 {
		fields = allFields
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		h.ServeHTTP(w, r)
		end := time.Now().UTC()

		vals := map[string]string{
			"start":               start.Format(tf),
			"end":                 end.Format(tf),
			"duration":            fmt.Sprintf(dfmt, end.Sub(start).Seconds()),
			"proto":               r.Proto,
			"host":                r.Host,
			"method":              r.Method,
			"uri":                 r.RequestURI,
			"id":                  r.Header.Get(hd),
			"path":                r.URL.Path,
			"origin":              r.Header.Get("Origin"),
			"body_bytes_received": strconv.FormatInt(r.ContentLength, 10),
			"user_agent":          r.UserAgent(),
			"remote_addr":         r.RemoteAddr,
		}
		if ww, ok := w.(interface {
			Status() int
		}); ok {
			vals["status"] = strconv.Itoa(ww.Status())
		}
		if ww, ok := w.(interface {
			Size() int
		}); ok {
			vals["body_bytes_sent"] = strconv.Itoa(ww.Size())
		}

		args := make([]interface{}, 0, len(fields)*2)
		for _, f := range fields {
			args = append(args, f, vals[f])
		}
		lr.Logger.Log(args...)
	})
}
