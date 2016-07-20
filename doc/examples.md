To put into README eventually:

## CSRF

Using gorilla/csrf:

```
wrapper := csrf.Protect([]byte("the secret key..."), ...options)
h := turtles.Wrap(h, turtles.WrapperFunc(wrapper))
// h is CSRF-protected, use csrf.Token(r) to retrieve the token on user login.
```

