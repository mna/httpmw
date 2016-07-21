To put into README eventually:

## CSRF

Using gorilla/csrf:

```
wrapper := csrf.Protect([]byte("the secret key..."), ...options)
h := httpmw.Wrap(h, httpmw.WrapperFunc(wrapper))
// h is CSRF-protected, use csrf.Token(r) to retrieve the token on user login.
```

## GZIP

Using http://godoc.org/github.com/NYTimes/gziphandler.
