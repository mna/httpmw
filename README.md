# httpmw [![GoDoc](https://godoc.org/github.com/PuerkitoBio/httpmw?status.png)][godoc] [![Build Status](https://semaphoreci.com/api/v1/mna/httpmw/branches/master/badge.svg)](https://semaphoreci.com/mna/httpmw)

Package httpmw is a collection of bite-sized middleware with chaining support. Uses the standard library's `http.Handler` interface. It's 🐢  🐢  🐢  all the way down. See the [godoc][] for full documentation.

## Installation

```
$ go get github.com/PuerkitoBio/httpmw/...
```

Use `-u` to update, `-t` to install test dependencies.

## Example

```
func myHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hello, middleware...")
}

// add a random request ID to all requests, with default configuration
var rid requestid.RequestID

// extract the real client remote address, with default configuration
var ra remoteip.RemoteIP

// limit the request body size to 1024 bytes
bl := bodylimit.BodyLimit{N: 1024}

// add CSRF support - this uses the github.com/gorilla/csrf external
// package, but it is very easy to adapt any http.Handler-compliant
// middleware out there. csrf.Protect returns a middleware-friendly
// `func(http.Handler) http.Handler`, which can be adapted to an
// httpmw.Wrapper using httpmw.WrapperFunc (much like http.Handler/
// http.HandlerFunc).
protect := httpmw.WrapperFunc(csrf.Protect([]byte(/* the secret */)))

h := httpmw.Wrap(http.HandlerFunc(myHandler), &rid, &ra, &bl, protect)

// serve using the middleware-augmented handler
log.Fatal(http.ListenAndServe(":9000", h))
```

## License

The [BSD 3-clause][bsd] license, see LICENSE file.

[bsd]: http://opensource.org/licenses/BSD-3-Clause
[godoc]: http://godoc.org/github.com/PuerkitoBio/httpmw

