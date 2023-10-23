module github.com/go-pkgz/auth/_example

go 1.15

replace github.com/go-pkgz/auth => ../

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-oauth2/oauth2/v4 v4.5.1
	github.com/go-pkgz/auth v1.18.0
	github.com/go-pkgz/lgr v0.10.4
	github.com/go-pkgz/rest v1.16.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/tidwall/buntdb v1.2.8 // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
)
